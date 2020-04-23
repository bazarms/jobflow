// +build integration

package github_integration

import (
	"context"
	"fmt"
	"os"
	//"strings"
	"testing"
	//"time"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/google/go-github/v28/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	//mygithub "github.com/bazarms/jobflow/plugins/github"
	"github.com/bazarms/jobflow/job"
	log "github.com/uthng/golog"
)

var repoTestName string
var repoUser string
var repoToken string
var githubClient *github.Client
var readmeContent = []byte("This is the content of my file\nand the 2nd line of it")

func init() {
	repoToken = os.Getenv("GITHUB_AUTH_TOKEN")
	if repoToken == "" {
		fmt.Println("No token found in var env")
		os.Exit(-1)
	}

	repoUser = os.Getenv("GITHUB_USER")
	if repoUser == "" {
		fmt.Println("No user specified in var env")
		os.Exit(-1)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: repoToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	githubClient = github.NewClient(tc)

	_, err := prepareRepoTest()

	if err != nil {
		log.Fatalln(err)
	}
}

func prepareRepoTest() (*github.Repository, error) {
	for {
		repoTestName = fmt.Sprintf("test-%d", rand.Int())
		_, resp, err := githubClient.Repositories.Get(context.Background(), repoUser, repoTestName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				// found a non-existent repo, perfect
				break
			}

			return nil, err
		}
	}

	// create the repository
	repo, _, err := githubClient.Repositories.Create(context.Background(), "", &github.Repository{Name: github.String(repoTestName), AutoInit: github.Bool(false)})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return repo, nil
}

func deleteRepoTest() error {
	_, err := githubClient.Repositories.Delete(context.Background(), repoUser, repoTestName)
	return err
}

func TestCmdRelease(t *testing.T) {
	flow := "data/flow/flow.yml"
	testCases := []struct {
		job     string
		commits []string
		params  map[string]interface{}
		output  *github.RepositoryRelease
	}{
		{
			"FirstRelease",
			[]string{
				"fix(feat1): commit 2\n\nCloses #1",
			},
			&github.RepositoryRelease{
				TagName: github.String("0.1.0"),
				Name:    github.String("0.1.0"),
				Body:    github.String("\\[\\w+\\] feat1: commit 2, \\(#1\\)\n\\[\\w+\\] feat1: commit 1\n"),
				Assets: []github.ReleaseAsset{
					{
						Name: github.String("asset1.tar.gz"),
					},
					{
						Name: github.String("asset2.tar.gz"),
					},
				},
			},
		},
		{
			"ReplaceRelease",
			[]string{
				"fix(feat1): commit 3",
			},
			&github.RepositoryRelease{
				TagName: github.String("0.1.0"),
				Name:    github.String("0.1.0"),
				Body:    github.String("\\[\\w+\\] feat1: commit 3\n\\[\\w+\\] feat1: commit 2, \\(#1\\)\n\\[\\w+\\] feat1: commit 1\n"),
				Assets: []github.ReleaseAsset{
					{
						Name: github.String("asset3.tar.gz"),
					},
					{
						Name: github.String("asset4.tar.gz"),
					},
				},
			},
		},
		{
			"SecondRelease",
			[]string{
				"feat(feat2): commit 4\nCloses #2",
				"feat(feat3): commit 5\nCloses #3",
			},
			&github.RepositoryRelease{
				TagName: github.String("0.2.0"),
				Name:    github.String("0.2.0"),
				Body:    github.String("\\[\\w+\\] feat3: commit 5, \\(#3\\)\n\\[\\w+\\] feat2: commit 4, \\(#2\\)\n"),
				Assets: []github.ReleaseAsset{
					{
						Name: github.String("asset1.tar.gz"),
					},
					{
						Name: github.String("asset2.tar.gz"),
					},
				},
			},
		},
	}

	//log.SetVerbosity(log.DEBUG)

	// Load modules if exists
	pluginDir := "../../bin/" + runtime.GOOS + "_" + runtime.GOARCH + "/plugins"
	err := job.LoadModules(pluginDir)
	if err != nil {
		log.Fatalln(err)
	}

	jf := config.ReadFlowFile(flow)
	jf.PluginDir = pluginDir

	// Add first file to get tree & commit
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("feat(feat1): commit 1"),
		Content:   readmeContent,
		Branch:    github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String(repoUser), Email: github.String("user@example.com")},
	}

	res, _, err := githubClient.Repositories.CreateFile(context.Background(), repoUser, repoTestName, "myNewFile.md", opts)
	if !assert.Nil(t, err) {
		return
	}

	readmeCommit := res.Commit

	for _, tc := range testCases {
		t.Run(tc.job, func(t *testing.T) {

			// Create commits
			for _, commit := range tc.commits {
				// Get ref
				ref, _, err := githubClient.Git.GetRef(context.Background(), repoUser, repoTestName, "refs/heads/master")
				if !assert.Nil(t, err) {
					return
				}

				// Get commit related to ref
				parent, _, err := githubClient.Repositories.GetCommit(context.Background(), repoUser, repoTestName, *ref.Object.SHA)
				if !assert.Nil(t, err) {
					return
				}
				// This is not always populated, but is needed.
				parent.Commit.SHA = parent.SHA

				// Create new commit
				newCommit, _, err := githubClient.Git.CreateCommit(context.Background(), repoUser, repoTestName, &github.Commit{Message: github.String(commit), Parents: []github.Commit{*parent.Commit}, Tree: readmeCommit.Tree})
				if !assert.Nil(t, err) {
					return
				}

				// Attach new commit to ref
				ref.Object.SHA = newCommit.SHA
				_, _, err = githubClient.Git.UpdateRef(context.Background(), repoUser, repoTestName, ref, false)
				if !assert.Nil(t, err) {
					return
				}
			}

			err = jf.RunJob(tc.job)
			if !assert.Nil(t, err) {
				return
			}

			// Check release
			result := jf.GetTaskResult(tc.job, "task-1")
			release := result.Result["release"].(*github.RepositoryRelease)
			assert.Equal(t, tc.output.GetTagName(), release.GetTagName())
			assert.Equal(t, tc.output.GetName(), release.GetName())
			assert.Regexp(t, tc.output.GetBody(), release.GetBody())

			// Check assets
			assets, _, err := githubClient.Repositories.ListReleaseAssets(context.Background(), repoUser, repoTestName, release.GetID(), nil)
			if !assert.Nil(t, err) {
				return
			}

			for index, asset := range assets {
				assert.Equal(t, asset.GetName(), tc.output.Assets[index].GetName())
			}
		})
	}

	err = deleteRepoTest()
	if err != nil {
		fmt.Println(err)
	}
}
