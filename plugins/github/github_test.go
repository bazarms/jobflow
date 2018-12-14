package github

import (
	//"fmt"
	"context"
	//"os"
	//"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"

	//"github.com/uthng/gojobs"
	//"github.com/uthng/gojobs/plugins/github"
	log "github.com/uthng/golog"
)

type commit struct {
	sha     string
	message string
}

type release struct {
	tag       string
	commitish string
}

type fakeClient struct {
	commits  []*commit
	releases []*release
}

func (c *fakeClient) GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error) {
	log.Infoln("GetLatestRelease")

	return nil, nil, nil
}

func (c *fakeClient) ListCommits(ctx context.Context, owner, repo string, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	commits := []*github.RepositoryCommit{}

	if len(c.commits) > 0 {
		for _, c := range c.commits {
			// In case of no sha in opt is specified (add all)
			// or sha specified
			if opt.SHA == "" || opt.SHA == c.sha {
				commits = append(commits, &github.RepositoryCommit{
					SHA: &c.sha,
					Commit: &github.Commit{
						Message: &c.message,
					},
				})
			}
		}
	}

	return commits, nil, nil
}

func (c *fakeClient) GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	release := &github.RepositoryRelease{}

	if len(c.releases) > 0 {
		for _, r := range c.releases {
			if r.tag == tag {
				release.TagName = &r.tag
				release.TargetCommitish = &r.commitish
				return release, nil, nil
			}
		}
	}
	return nil, nil, nil
}

func (c *fakeClient) DeleteRelease(ctx context.Context, owner, repo string, id int64) (*github.Response, error) {
	return nil, nil
}

func (c *fakeClient) DeleteReleaseAsset(ctx context.Context, owner, repo string, id int64) (*github.Response, error) {
	return nil, nil
}

func (c *fakeClient) CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error) {
	return release, nil, nil
}

func newFakeClient() *client {
	c := &client{
		ctx: context.Background(),
	}

	fc := &fakeClient{}
	c.repositories = fc

	return c
}

func getPtrString(str string) *string {
	return &str
}

func getPtrBool(b bool) *bool {
	return &b
}

func TestGenerateChangeLog(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name   string
		client *client
		output string
	}{
		{
			"Disabled",
			&client{
				ctx:           ctx,
				changelog:     false,
				changelogType: COMMIT,
				repositories: &fakeClient{
					commits: []*commit{
						{
							sha:     "197b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "feat(core): subject 1\n\nBody1\n\nClosed #1, resolved #500",
						},
					},
				},
			},
			"",
		},
		{
			"TypeError",
			&client{
				ctx:           ctx,
				changelog:     true,
				changelogType: 4,
				repositories: &fakeClient{
					commits: []*commit{
						{
							sha:     "197b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "feat(core): subject 1\n\nBody1\n\nClosed #1, resolved #500",
						},
					},
				},
			},
			"",
		},
		{
			"TypeCommitOK",
			&client{
				ctx:           ctx,
				changelog:     true,
				changelogType: COMMIT,
				repositories: &fakeClient{
					commits: []*commit{
						{
							sha:     "364b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "feat(core): subject 1\n\nBody1\n\nClosed #1, resolved #500",
						},
						{
							sha:     "989b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "fix: subject 2 (#1234)\n\nBody2\n\nFixed #2, fix #120",
						},
						{
							sha:     "197b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "Subject 3\n\nBody3\n\nClosed #3, fixed ex_repo/ex_user#234, fixes #200",
						},
					},
				},
			},
			"[364b53e] core: subject 1, (#1, #500)\n[989b53e] subject 2 (#1234), (#2, #120)\n[197b53e] Subject 3, (#3, ex_repo/ex_user#234, #200)\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changelog := tc.client.generateChangelog("")

			assert.Equal(t, tc.output, changelog)
		})
	}
}

func TestCreateRelease(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name    string
		client  *client
		release *github.RepositoryRelease
	}{
		{
			"NoRelease_TypeCommitOK",
			&client{
				ctx:           ctx,
				tag:           "0.2.1",
				commitish:     "master",
				name:          "0.2.1",
				changelog:     true,
				changelogType: COMMIT,
				repositories: &fakeClient{
					commits: []*commit{
						{
							sha:     "364b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "feat(core): subject 1\n\nBody1\n\nClosed #1, resolved #500",
						},
						{
							sha:     "989b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "fix: subject 2 (#1234)\n\nBody2\n\nFixed #2, fix #120",
						},
						{
							sha:     "197b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "Subject 3\n\nBody3\n\nClosed #3, fixed ex_repo/ex_user#234, fixes #200",
						},
					},
					releases: []*release{
						{
							tag:       "0.2.0",
							commitish: "197b53e7abf3b56b8e984c55ce9bebef8ee016eb",
						},
					},
				},
			},
			&github.RepositoryRelease{
				TagName:         getPtrString("0.2.1"),
				TargetCommitish: getPtrString("master"),
				Name:            getPtrString("0.2.1"),
				Body:            getPtrString("[364b53e] core: subject 1, (#1, #500)\n[989b53e] subject 2 (#1234), (#2, #120)\n[197b53e] Subject 3, (#3, ex_repo/ex_user#234, #200)\n"),
				Draft:           getPtrBool(false),
				Prerelease:      getPtrBool(false),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			release, err := tc.client.createRelease()
			log.Infoln(release)
			assert.Nil(t, err)
			assert.Equal(t, tc.release, release)
		})
	}
}
