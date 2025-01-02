package github

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/funk"
)

type Commit struct {
	*github.Repository
	*github.RepositoryCommit
}

type listCommitsConfig struct {
	since         time.Time
	until         time.Time
	repoAllowList RepoAllowList
	withDetail    bool
	fileAllowList FileAllowList
}

type FileAllowList []string

func (allowlist FileAllowList) Allows(file *github.CommitFile) bool {
	for _, pattern := range allowlist {
		match, err := filepath.Match(pattern, *file.Filename)
		if err != nil {
			continue
		}
		if match == true {
			return true
		}
	}
	return false
}

func ListCommitsExclusiveTo(repos []*github.Repository) funk.Option[listCommitsConfig] {
	return func(config *listCommitsConfig) {
		config.repoAllowList = repos
	}
}

func ListCommitsSince(since time.Time) funk.Option[listCommitsConfig] {
	return func(config *listCommitsConfig) {
		config.since = since
	}
}

func ListCommitsUntil(until time.Time) funk.Option[listCommitsConfig] {
	return func(config *listCommitsConfig) {
		config.until = until
	}
}

func ListCommitsToFiles(patterns []string) funk.Option[listCommitsConfig] {
	return func(config *listCommitsConfig) {
		config.fileAllowList = patterns
	}
}

func ListCommitsWithDetail(withDetail bool) funk.Option[listCommitsConfig] {
	return func(config *listCommitsConfig) {
		config.withDetail = withDetail
	}
}

func (c *Client) ListCommits(ctx context.Context, opts ...funk.Option[listCommitsConfig]) (<-chan Commit, <-chan error) {
	var config = funk.ConfigWithOptions[listCommitsConfig](opts)
	return channel.Async(func(outchan chan Commit, errchan chan error) {
		var wg sync.WaitGroup

		allowedRepos, errs := c.ListRepos(ctx, ListReposAllowedBy(config.repoAllowList))
		channel.GoFwd(ctx, &wg, errs, errchan)

		channel.ForEach(ctx, allowedRepos, func(repo *github.Repository) {
			commits, errs := c.ListCommitsByRepo(ctx, repo,
				ListCommitsByRepoSince(config.since),
				ListCommitsByRepoUntil(config.until),
			)
			channel.GoFwd(ctx, &wg, errs, errchan)

			// must get detail if there is a file allow list
			if len(config.fileAllowList) > 0 {
				config.withDetail = true
			}

			if config.withDetail == true {
				commits, errs = channel.Map(ctx, commits, func(commit Commit, outs chan<- Commit, errs chan<- error) {
					ownerLogin, repoName := *commit.Repository.Owner.Login, *commit.Repository.Name
					c.log.Debugf("*github.Client.Repositories.GetCommit(owner=%s, repo=%s, sha=%s)", ownerLogin, repoName, *commit.SHA)
					detail, _, err := c.gh.Repositories.GetCommit(ctx, ownerLogin, repoName, *commit.SHA, nil)
					if err != nil {
						errs <- fmt.Errorf("error getting commit; repo=%s sha=%s", *commit.Repository.FullName, *commit.SHA)
					}
					outs <- Commit{Repository: commit.Repository, RepositoryCommit: detail}
				})
				channel.GoFwd(ctx, &wg, errs, errchan)
			}

			if len(config.fileAllowList) > 0 {
				commits = channel.Filter(ctx, commits, func(commit Commit) bool {
					for _, file := range commit.Files {
						if config.fileAllowList.Allows(file) {
							return true
						}
					}
					return false
				})
			}

			channel.Fwd(ctx, commits, outchan)
		})

		wg.Wait() // I'm not confident that I actually need this
	})
}

type listCommitsByRepoConfig struct {
	since time.Time
	until time.Time
}

func ListCommitsByRepoSince(since time.Time) funk.Option[listCommitsByRepoConfig] {
	return func(config *listCommitsByRepoConfig) {
		config.since = since
	}
}

func ListCommitsByRepoUntil(until time.Time) funk.Option[listCommitsByRepoConfig] {
	return func(config *listCommitsByRepoConfig) {
		config.until = until
	}
}

func (c *Client) ListCommitsByRepo(ctx context.Context, repo *github.Repository, opts ...funk.Option[listCommitsByRepoConfig]) (<-chan Commit, <-chan error) {
	var config = funk.ConfigWithOptions[listCommitsByRepoConfig](opts)
	return channel.Async(func(outchan chan Commit, errchan chan error) {
		opt := &github.CommitsListOptions{
			Since: config.since,
			Until: config.until,
			/* Path: "go.mod",
				^ I can't use this here, but thought I'd write a note about it.
				This tells GitHub to only return commits that include this path.
			  In this example, only commits that affected `go.mod` will be returned.
				From what I can tell, this can only be a single path and it must be
				relative to git root. It can't be a glob or a list of paths.
			*/
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			c.log.Debugf("*github.Client.Repositories.ListCommits(owner=%s, repo=%s, page=%d)", repo.Owner, repo, opt.Page)
			commits, res, err := c.gh.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing commits by repo; repo=%s opt=%+v: %w", repo, opt, err)
			}

			for _, commit := range commits {
				c.log.Debugf("\tsha=%s", *commit.SHA)
				outchan <- Commit{repo, commit}
			}

			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	})
}
