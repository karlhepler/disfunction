package github

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/funk"
)

type Commit struct {
	Repo
	*github.RepositoryCommit
}

type listCommitsConfig struct {
	since time.Time
	until time.Time
	// commits []Commit // I don't think I need this
	repoAllowList []Repo
	withDetail    bool
}

func ListCommitsInReposExclusiveTo(repos []Repo) funk.Option[listCommitsConfig] {
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

func ListCommitsWithDetail() funk.Option[listCommitsConfig] {
	return func(config *listCommitsConfig) {
		config.withDetail = true
	}
}

func (c *Client) ListCommits(ctx context.Context, opts ...funk.Option[listCommitsConfig]) (<-chan Commit, <-chan error) {
	var config = funk.ConfigWithOptions[listCommitsConfig](opts)
	return channel.Async(func(outchan chan Commit, errchan chan error) {
		var wg sync.WaitGroup

		// If I have an allow list, and I can fill it with the right information, which I cane,
		// then there is no reason to do this at all. If I give a list of repos,
		// then just iterate over that list instead of doing this.
		var repos <-chan Repo
		var errs <-chan error
		if len(config.repoAllowList) > 0 {
			repos = channel.SendEachOnChannel(config.repoAllowList)
		} else {
			repos, errs = c.ListRepos(ctx)
			channel.GoFwd(ctx, &wg, errs, errchan)
		}

		channel.ForEach(ctx, repos, func(repo Repo) {
			commits, errs := c.ListCommitsByRepo(ctx, repo,
				ListCommitsByRepoSince(config.since),
				ListCommitsByRepoUntil(config.until),
			)
			channel.GoFwd(ctx, &wg, errs, errchan)
			// TODO(karlhepler): I need to change this
			channel.GoFwd(ctx, &wg, commits, outchan)
		})

		wg.Wait()
	})
}

func (c *Client) ListDetailedCommits(ctx context.Context, opts ...funk.Option[listCommitsConfig]) (<-chan Commit, <-chan error) {
	commits, listCommitsErrs := c.ListCommits(ctx, opts...)
	return channel.Async(func(outchan chan Commit, errchan chan error) {
		var wg sync.WaitGroup
		channel.GoFwd(ctx, &wg, listCommitsErrs, errchan)

		channel.ForEach(ctx, commits, func(commit Commit) {
			c.log.Debugf("*github.Client.Repositories.GetCommit(owner=%s, repo=%s, sha=%s)", commit.Repo.Owner, commit.Repo.Name, *commit.SHA)
			detailedCommit, res, err := c.gh.Repositories.GetCommit(ctx, commit.Repo.Owner.String(), commit.Repo.Name, *commit.SHA, nil)
			if err != nil {
				errchan <- fmt.Errorf("error getting commit; repo=%s sha=%s", commit.Repo, *commit.SHA)
			}

			outchan <- Commit{
				Repo:             commit.Repo,
				RepositoryCommit: detailedCommit,
			}

			if res == nil {
				return
			}
		})

		wg.Wait()
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

func (c *Client) ListCommitsByRepo(ctx context.Context, repo Repo, opts ...funk.Option[listCommitsByRepoConfig]) (<-chan Commit, <-chan error) {
	var config = funk.ConfigWithOptions[listCommitsByRepoConfig](opts)
	return channel.Async(func(outchan chan Commit, errchan chan error) {
		opt := &github.CommitsListOptions{
			Since: config.since,
			Until: config.until,
			/* Path: "go.mod",
				^ I can't use this here, but thought I'd write a note about it.
				This tells GitHub to only return commits that include this path.
			  In this example, only commits that affected `go.mod` will be returned.
			*/
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			c.log.Debugf("*github.Client.Repositories.ListCommits(owner=%s, repo=%s, page=%d)", repo.Owner, repo, opt.Page)
			commits, res, err := c.gh.Repositories.ListCommits(ctx, repo.Owner.String(), repo.Name, opt)
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
