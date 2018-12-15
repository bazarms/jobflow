package github

import "C"
import (
	"context"
	"fmt"
	//"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/github"
	//"github.com/spf13/cast"
	"golang.org/x/oauth2"

	"github.com/uthng/gojobs"
	log "github.com/uthng/golog"
)

type repositoriesService interface {
	GetCommit(ctx context.Context, owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error)
	ListCommits(ctx context.Context, owner, repo string, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	DeleteRelease(ctx context.Context, owner, repo string, id int64) (*github.Response, error)
	DeleteReleaseAsset(ctx context.Context, owner, repo string, id int64) (*github.Response, error)
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
}

type client struct {
	ctx context.Context

	// Repository
	user       string
	repository string
	tag        string
	commitish  string

	// Release
	name        string
	description string

	// Changelog
	changelog     bool
	changelogType int

	// Options
	draft      bool
	prerelease bool
	delete     bool
	replace    bool
	soft       bool

	repositories repositoriesService
}

const (
	// COMMIT = 0
	COMMIT = iota
	// ISSUE = 1
	ISSUE = 1
)

var module = gojobs.Module{
	Name:        "git",
	Version:     "0.1",
	Description: "Commands related to git",
}

// List of available commands for this module
var commands = []gojobs.Cmd{
	{
		Name:   "release",
		Func:   CmdRelease,
		Module: module,
	},
}

// Init initializes module by registering all its commands
// to command registry
func init() {
	for _, cmd := range commands {
		gojobs.CmdRegister(cmd)
	}
}

// CmdRelease creates a release on github: changelog, release and maybe upload assets
//
// Params:
// - user: github user
// - token: github api token
// - repository: repository name
// - version: tag name
// - commitish: branch name or commit SHA
// - name: release name
// - description: release description. Possible values: "changelog" or your content. If param is not specified, "changelog" is used.
// - changelog: true/false. Generate changelog or not. Default value: true event if it is not specified
// - changelog_type: 0 = "commit", 1 = "issue". By default 0 = "commit" is used
// - draft: true/false. Just a draft and no publish
// - prerelease: true/false
// - delete: true/false. Delete release and its git tag in advance if it exists
// - replace: replace artifacts if it is already uploaded
// - soft: true/false. Stop uploading if the same tag already exists
func CmdRelease(params map[string]interface{}) *gojobs.CmdResult {
	//var value interface{}
	var result = gojobs.NewCmdResult()
	//var token string
	//var user string
	//var repository string
	//var version string
	//var commitish string
	//var name string
	//var description = "changelog"
	//var changelog = true
	//var changelogType = "commit"

	//value, ok := params["token"]
	//if ok == false {
	//result.Error = fmt.Errorf("param token missing")
	//return result
	//}
	//token = cast.ToString(value)

	//value, ok = params["user"]
	//if ok == false {
	//result.Error = fmt.Errorf("param user missing")
	//return result
	//}
	//user = cast.ToString(value)

	//value, ok = params["repository"]
	//if ok == false {
	//result.Error = fmt.Errorf("param repository missing")
	//return result
	//}
	//repository = cast.ToString(value)

	//value, ok = params["version"]
	//if ok == false {
	//result.Error = fmt.Errorf("param version missing")
	//return result
	//}
	//version = cast.ToString(value)

	//value, ok = params["commitish"]
	//if ok == false {
	//result.Error = fmt.Errorf("param commit missing")
	//return result
	//}
	//commitish = cast.ToString(value)

	//value, ok = params["name"]
	//if ok == false {
	//result.Error = fmt.Errorf("param name missing")
	//return result
	//}
	//name = cast.ToString(value)

	//value, ok = params["description"]
	//if ok {
	//description = cast.ToString(value)
	//}

	//value, ok = params["changelog"]
	//if ok {
	//changelog = cast.ToBool(value)
	//}

	//value, ok = params["changelog_type"]
	//if ok {
	//changelogType = cast.ToInt(value)
	//}

	//client := newClientByToken(token)

	// list all repositories for the authenticated user
	//latestRelease, _, _ := client.Repositories.GetLatestRelease(client.ctx, user, repository)
	//log.Infoln("latest release", latestRelease)

	return result
}

/////////////////// INTERNAL FUNCTION //////////////////

// newClientByToken returns a new github client with context
// using github token
func newClientByToken(token string) *client {
	c := &client{
		ctx: context.Background(),
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(c.ctx, ts)

	githubClient := github.NewClient(tc)

	c.repositories = githubClient.Repositories

	return c
}

// generateChangelog generates changelog using issues or commits
//
// If changelog is enabled and its type is issue,
// it tries to get the latest release first
// to find the sha from which it will generate changelog.
// If no release hasnt created yet, it will do from the beginning.
func (c *client) generateChangelog(from, to time.Time) string {
	var commits []*github.RepositoryCommit
	var msgs string
	var err error

	log.Debugw("Generate changelog", "from", from.Local(), "to", to.Local())
	if !c.changelog {
		log.Warnln("Option changelog is disabled. Do nothing")
		return msgs
	}

	if c.changelogType < COMMIT || c.changelogType > ISSUE {
		log.Warnw("Option changelogType is not supported. Do nothing", "changelog_type", c.changelogType)
		return msgs
	}

	// A starting commit given from the latest release
	// and changelog type is issue
	if c.changelogType == COMMIT {
		opt := &github.CommitsListOptions{}

		if !from.IsZero() {
			opt.Since = from
		}

		if !to.IsZero() {
			opt.Until = to
		}

		commits, _, err = c.repositories.ListCommits(c.ctx, c.user, c.repository, opt)
		log.Debugw("List commits", "opt", opt, "commits", commits)
		if err != nil {
			log.Errorw("Error while getting list of commits", "from", from)
			return msgs
		}

		for _, commit := range commits {
			msgs += formatCommitChangelog(commit.GetSHA(), commit.Commit.GetMessage()) + "\n"
		}
	} else {

	}

	return msgs
}

// createRelease creates a new github release with
// the given version (tag name), name and commitish.
//
// If version (tag) is already created and delete option is set,
// it will delete existing release, its tags to create new ones completely.
// If replace is set, it will only regenerate and replace the existing
// changelog and artifacts.
//
// Before, it creates a new release, it will do some calls to perform:
// - Generate changelog
func (c *client) createRelease() (*github.RepositoryRelease, error) {
	var from time.Time
	var to time.Time

	// Check if release that we want to create already exists or not
	wantedRelease, _, err := c.repositories.GetReleaseByTag(c.ctx, c.user, c.repository, c.tag)
	if err != nil {
		log.Errorw("Cannot get latest release", "tag", c.tag, "err", err)
		return nil, err
	}
	log.Debugw("Get releases by tag", "tag", c.tag, "release", wantedRelease)
	// Exist a release with same tag
	// Check different options
	if wantedRelease != nil {
		// Delete all: release, assets to redo
		if c.replace {
			_, err := c.repositories.DeleteRelease(c.ctx, c.user, c.repository, *wantedRelease.ID)
			if err != nil {
				log.Errorw("Cannot delete existing release", "tag", c.tag, "err", err)
				return nil, err
			}

			_, err = c.repositories.DeleteReleaseAsset(c.ctx, c.user, c.repository, *wantedRelease.ID)
			if err != nil {
				log.Errorw("Cannot delete existing release assets", "tag", c.tag, "err", err)
				return nil, err
			}
		}
	}

	// Get latest release => date
	// if no release until now => no date from
	latestRelease, _, err := c.repositories.GetLatestRelease(c.ctx, c.user, c.repository)
	log.Debugw("Get latest release", "release", latestRelease)
	if err != nil {
		log.Errorw("Error while getting latest release", "err", err)
		return nil, err
	}

	if latestRelease != nil {
		from = latestRelease.GetCreatedAt().Time
	}

	// Get date of commitish => date to
	commit, _, err := c.repositories.GetCommit(c.ctx, c.user, c.repository, c.commitish)
	log.Debugw("Get commit", "commitish", c.commitish, "commit", commit)
	if err != nil {
		log.Errorw("Error while getting commit", "sha", c.commitish, "err", err)
		return nil, err
	}

	if commit == nil {
		log.Errorw("No commit found", "sha", c.commitish)
		return nil, fmt.Errorf("No commit found for sha: %s", c.commitish)
	}

	to = commit.Commit.GetCommitter().GetDate()

	// Get list of commits from the date of latest release to
	// the wanted release created if exists or to the date of
	// the commitish in param.
	msgs := c.generateChangelog(from, to)
	if msgs == "" {
		return nil, fmt.Errorf("cannot generate changelog")
	}

	// Create new release
	newRelease := &github.RepositoryRelease{
		TagName:         &c.tag,
		TargetCommitish: &c.commitish,
		Name:            &c.name,
		Body:            &msgs,
		Draft:           &c.draft,
		Prerelease:      &c.prerelease,
	}

	r, _, err := c.repositories.CreateRelease(c.ctx, c.user, c.repository, newRelease)
	if err != nil {
		log.Errorw("Cannot create new release", "release", r)
		return nil, err
	}

	return r, nil
}

// formatCommitChangelog parses commit message to build a changelog message
// with the following format: [<scope>:] <subject> <sha>, [issues...]
//
// It is recommanded to write commit msg with the format below: <type>(<scope>): <subject>\n\n<body>\n\n<footer>
func formatCommitChangelog(sha, msg string) string {
	var format string
	var issues string
	var args []interface{}

	lines := strings.Split(msg, "\n\n")
	summary := lines[0]
	footer := lines[len(lines)-1]

	// Add short SHA
	args = append(args, "["+sha[0:7]+"]")
	format += "%s "

	// Check to find type & scope & subject in summary line
	// Regexp: 2 patterns: <type(scope): > & <type: ><subject>
	// 0. feat(core): subject 1
	// 1. feat(core):
	// 2. (core)
	// 3. core
	// 4. subject 1
	re, err := regexp.Compile(`(\w+(\((\w+)\)): |\w+: )?(.*)`)
	if err != nil {
	}

	resSummary := re.FindStringSubmatch(summary)

	// Scope
	if resSummary[3] != "" {
		args = append(args, resSummary[3])
		format += "%s: "
	}

	// Subject
	if resSummary[4] != "" {
		args = append(args, resSummary[4])
		format += "%s,"
	}

	// Check to capture #issue using gitlub closing keywords
	// Ex: Closed #3, fixed ex_repo/ex_user#234, fixes #200
	// 0. ["Closed #3" "#3"]
	// 1.["fixed ex_repo/ex_user#234" "ex_repo/ex_user#234"]
	// 2. ["fixed #234" "#234"]
	re, err = regexp.Compile(`(?:[C|c]lose[s|d]{0,1}|[F|f]ixe[s|d]{0,1}|[F|f]ix|[R|r]esolve[s|d]{0,1}) ([^ ]*#\d{1,5})`)
	if err != nil {
	}
	resFooter := re.FindAllStringSubmatch(footer, -1)

	issues += "("
	for k, v := range resFooter {
		issues += v[1]
		if k < len(resFooter)-1 {
			issues += ", "
		}
	}
	issues += ")"
	args = append(args, issues)

	if len(resFooter) > 0 {
		format += " %s"
	}

	changelog := fmt.Sprintf(format, args...)

	return changelog
}
