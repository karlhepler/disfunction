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

type listCommitsConfig struct {
	owner Owner
	since time.Time
	until time.Time
}

func FilterCommitsByOwner(owner Owner) funk.Option[listCommitsConfig] {
	return func(config *listCommitsConfig) {
		config.owner = owner
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

func (c *Client) ListCommits(ctx context.Context, opts ...funk.Option[listCommitsConfig]) (<-chan Commit, <-chan error) {
	var config = funk.ConfigWithOptions[listCommitsConfig](opts...)
	return channel.Async(func(outchan chan Commit, errchan chan error) {
		var wg sync.WaitGroup

		repos, errs := c.ListRepos(ctx, FilterReposByOwner(config.owner))
		channel.GoForward(ctx, &wg, errs, errchan)
		channel.ForEach(ctx, repos, func(repo *github.Repository) {
			commits, errs := c.ListCommitsByRepo(ctx,
				Repo{
					Owner: Owner(*repo.Owner.Login),
					Name:  *repo.Name,
				},
				ListCommitsByRepoSince(config.since),
				ListCommitsByRepoUntil(config.until),
			)
			channel.GoForward(ctx, &wg, errs, errchan)
			channel.GoForward(ctx, &wg, commits, outchan)
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
	var config = funk.ConfigWithOptions[listCommitsByRepoConfig](opts...)
	return channel.Async(func(outchan chan Commit, errchan chan error) {
		opt := &github.CommitsListOptions{
			Since:       config.since,
			Until:       config.until,
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
