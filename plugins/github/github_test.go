package github

import (
	"context"
	"fmt"
	"os"
	//"strings"
	"testing"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/stretchr/testify/assert"
	//"github.com/uthng/jobflow/job"
	//"github.com/uthng/jobflow/plugins/github"
	//log "github.com/uthng/golog"
)

type commit struct {
	sha     string
	message string
	date    time.Time
}

type release struct {
	tag       string
	commitish string
	createAt  time.Time
}

type tag struct {
	ref string
	sha string
}

type fakeClient struct {
	commits  []*commit
	releases []*release
	tags     []*tag
	//repositories repositoriesService
	//git          gitService
}

func (c *fakeClient) GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error) {
	var date time.Time
	var index int

	// Just simulate error
	if c.releases == nil {
		return nil, nil, fmt.Errorf("list of releases is nil")
	}

	// Loop release list to find out the release
	// with createAt most recent
	for i, r := range c.releases {
		if date.IsZero() || date.Unix() < r.createAt.Unix() {
			date = r.createAt
			index = i
		}
	}

	release := &github.RepositoryRelease{
		TagName:         &c.releases[index].tag,
		TargetCommitish: &c.releases[index].commitish,
		CreatedAt: &github.Timestamp{
			Time: c.releases[index].createAt,
		},
	}

	return release, nil, nil
}

func (c *fakeClient) ListCommits(ctx context.Context, owner, repo string, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	commits := []*github.RepositoryCommit{}

	// Just simulate error
	if c.commits == nil {
		return nil, nil, fmt.Errorf("list of commits is nil")
	}

	if len(c.commits) > 0 {
		for _, c := range c.commits {
			ok := false
			// Check if date since or until are set
			if opt.Since.IsZero() && opt.Until.IsZero() {
				ok = true
			} else if !opt.Since.IsZero() && opt.Until.IsZero() {
				if c.date.Unix() >= opt.Since.Unix() {
					ok = true
				}
			} else if opt.Since.IsZero() && !opt.Until.IsZero() {
				if c.date.Unix() <= opt.Until.Unix() {
					ok = true
				}
			} else {
				if c.date.Unix() >= opt.Since.Unix() && c.date.Unix() <= opt.Until.Unix() {
					ok = true
				}
			}

			if ok {
				commits = append(commits, &github.RepositoryCommit{
					SHA: &c.sha,
					Commit: &github.Commit{
						Message: &c.message,
						Committer: &github.CommitAuthor{
							Date: &c.date,
						},
					},
				})
			}
		}
	}

	return commits, nil, nil
}

func (c *fakeClient) ListReleases(ctx context.Context, owner, repo string, opt *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
	releases := []*github.RepositoryRelease{}

	// Just simulate error
	if c.releases == nil {
		return nil, nil, fmt.Errorf("error")
	}

	for _, r := range c.releases {
		release := &github.RepositoryRelease{}
		release.TagName = &r.tag
		release.TargetCommitish = &r.commitish
		release.CreatedAt = &github.Timestamp{
			Time: r.createAt,
		}
		releases = append(releases, release)
	}

	return releases, nil, nil
}

func (c *fakeClient) GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	release := &github.RepositoryRelease{}

	// Just simulate error
	if tag == "" || c.releases == nil {
		return nil, nil, fmt.Errorf("error")
	}

	if len(c.releases) > 0 {
		for _, r := range c.releases {
			if r.tag == tag {
				release.TagName = &r.tag
				release.TargetCommitish = &r.commitish
				release.CreatedAt = &github.Timestamp{
					Time: r.createAt,
				}
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

func (c *fakeClient) ListReleaseAssets(ctx context.Context, owner, repo string, id int64, opt *github.ListOptions) ([]*github.ReleaseAsset, *github.Response, error) {
	return nil, nil, nil
}

func (c *fakeClient) UploadReleaseAsset(ctx context.Context, owner, repo string, id int64, opt *github.UploadOptions, file *os.File) (*github.ReleaseAsset, *github.Response, error) {
	return nil, nil, nil
}

func (c *fakeClient) CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error) {
	// Just simulate error
	if release == nil {
		return nil, nil, fmt.Errorf("error")
	}

	return release, nil, nil
}

func (c *fakeClient) GetCommit(ctx context.Context, owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error) {

	// Just simulate error
	if sha == "" || c.commits == nil {
		return nil, nil, fmt.Errorf("error")
	}

	for _, c := range c.commits {
		if c.sha == sha {
			ghCommit := &github.RepositoryCommit{
				SHA: &c.sha,
				Commit: &github.Commit{
					Message: &c.message,
				},
			}
			return ghCommit, nil, nil
		}
	}

	return nil, nil, nil
}

func (c *fakeClient) CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error) {
	// Just simulate error
	if ref == nil {
		return nil, nil, fmt.Errorf("error")
	}

	return ref, nil, nil
}

func (c *fakeClient) GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error) {

	// Just simulate error
	if ref == "" || c.tags == nil {
		return nil, nil, fmt.Errorf("error")
	}

	for _, t := range c.tags {
		if t.ref == ref {
			refTag := &github.Reference{
				Ref: &t.ref,
				Object: &github.GitObject{
					SHA: &t.sha,
				},
			}
			return refTag, nil, nil
		}
	}

	// No match => error as real func does
	return nil, nil, fmt.Errorf("error")
}

func (c *fakeClient) DeleteRef(ctx context.Context, owner string, repo string, ref string) (*github.Response, error) {
	// Just simulate error
	if ref == "" {
		return nil, fmt.Errorf("error")
	}

	for i, t := range c.tags {
		if t.ref == ref {
			c.tags = append(c.tags[:i], c.tags[i+1:]...)
		}
	}

	return nil, nil
}

//func newFakeClient() *client {
//c := &client{
//ctx: context.Background(),
//}

//fc := &fakeClient{}
//c.repositories = fc
//c.git = fc

//return c
//}

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
			changelog := tc.client.generateChangelog(time.Time{}, time.Time{})

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
				commitish:     "197b53e7abf3b56b8e984c55ce9bebef8ee016eb",
				name:          "0.2.1",
				changelog:     true,
				changelogType: COMMIT,
				repositories: &fakeClient{
					commits: []*commit{
						{
							sha:     "364b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "feat(core): subject 1\n\nBody1\n\nClosed #1, resolved #500",
							date:    time.Date(2018, time.November, 5, 10, 0, 0, 0, time.UTC),
						},
						{
							sha:     "989b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "fix: subject 2 (#1234)\n\nBody2\n\nFixed #2, fix #120",
							date:    time.Date(2018, time.November, 10, 10, 0, 0, 0, time.UTC),
						},
						{
							sha:     "111b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "fix: subject 21\n\nBody21\n\nFixed #21",
							date:    time.Date(2018, time.November, 11, 10, 0, 0, 0, time.UTC),
						},
						{
							sha:     "197b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							message: "Subject 3\n\nBody3\n\nClosed #3, fixed ex_repo/ex_user#234, fixes #200",
							date:    time.Date(2018, time.November, 15, 10, 0, 0, 0, time.UTC),
						},
					},
					releases: []*release{
						{
							tag:       "0.2.0",
							commitish: "989b53e7abf3b56b8e984c55ce9bebef8ee016eb",
							createAt:  time.Date(2018, time.November, 10, 10, 0, 0, 0, time.UTC),
						},
					},
				},
				git: &fakeClient{
					tags: []*tag{
						{
							ref: "refs/tags/0.1.9",
							sha: "555b53e7abf3b56b8e984c55ce9bebef8ee016eb",
						},
						{
							ref: "refs/tags/0.2.0",
							sha: "989b53e7abf3b56b8e984c55ce9bebef8ee016eb",
						},
					},
				},
			},
			&github.RepositoryRelease{
				TagName:         getPtrString("0.2.1"),
				TargetCommitish: getPtrString("197b53e7abf3b56b8e984c55ce9bebef8ee016eb"),
				Name:            getPtrString("0.2.1"),
				Body:            getPtrString("[111b53e] subject 21, (#21)\n[197b53e] Subject 3, (#3, ex_repo/ex_user#234, #200)\n"),
				Draft:           getPtrBool(false),
				Prerelease:      getPtrBool(false),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//log.SetVerbosity(log.DEBUG)
			release, err := tc.client.createRelease()

			assert.Nil(t, err)
			assert.Equal(t, tc.release, release)
		})
	}
}
