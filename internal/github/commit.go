package github

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
)

type listCommitsConfig struct {
	owner Owner
	since time.Time
	until time.Time
}

type listCommitsOption func(*listCommitsConfig)

func FilterCommitsByOwner(owner Owner) listCommitsOption {
	return func(config *listCommitsConfig) {
		config.owner = owner
	}
}

func ListCommitsSince(since time.Time) listCommitsOption {
	return func(config *listCommitsConfig) {
		config.since = since
	}
}

func ListCommitsUntil(until time.Time) listCommitsOption {
	return func(config *listCommitsConfig) {
		config.until = until
	}
}

func (c *Client) ListCommits(ctx context.Context, opts ...listCommitsOption) (<-chan Commit, <-chan error) {
	var config listCommitsConfig
	for _, opt := range opts {
		opt(&config)
	}

	outchan, errchan := make(chan Commit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)
		var wg sync.WaitGroup

		repos, errs := c.ListRepos(ctx, FilterReposByOwner(config.owner))
		channel.GoForward(&wg, errs, errchan)
		channel.ForEach(repos, func(repo *github.Repository) {
			commits, errs := c.ListCommitsByRepo(ctx,
				Repo{
					Owner: Owner(*repo.Owner.Login),
					Name:  *repo.Name,
				},
				ListCommitsByRepoSince(config.since),
				ListCommitsByRepoUntil(config.until),
			)
			channel.GoForward(&wg, errs, errchan)
			channel.Forward(commits, outchan)
		})

		wg.Wait()
	}()
	return outchan, errchan
}

type listCommitsByRepoConfig struct {
	since time.Time
	until time.Time
}

type listCommitsByRepoOption func(*listCommitsByRepoConfig)

func ListCommitsByRepoSince(since time.Time) listCommitsByRepoOption {
	return func(config *listCommitsByRepoConfig) {
		config.since = since
	}
}

func ListCommitsByRepoUntil(until time.Time) listCommitsByRepoOption {
	return func(config *listCommitsByRepoConfig) {
		config.until = until
	}
}

func (c *Client) ListCommitsByRepo(ctx context.Context, repo Repo, opts ...listCommitsByRepoOption) (<-chan Commit, <-chan error) {
	var config listCommitsByRepoConfig
	for _, opt := range opts {
		opt(&config)
	}

	outchan, errchan := make(chan Commit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.CommitsListOptions{
			Since:       config.since,
			Until:       config.until,
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			c.Log.Debugf("GitHub.Repositories.ListCommits(owner=%s, repo=%s, page=%d)", repo.Owner, repo, opt.Page)
			commits, res, err := c.GitHub.Repositories.ListCommits(ctx, repo.Owner.String(), repo.Name, opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing commits by repo; repo=%s opt=%+v: %w", repo, opt, err)
			}

			for _, commit := range commits {
				c.Log.Debugf("\tsha=%s", *commit.SHA)
				outchan <- Commit{repo, commit}
			}

			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
