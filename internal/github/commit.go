package github

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
)

func (c *Client) ListOwnerCommitsByDateRange(ctx context.Context, owner Owner, date DateRange) (<-chan OwnerRepoCommit, <-chan error) {
	outchan, errchan := make(chan OwnerRepoCommit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)
		var wg sync.WaitGroup

		repos, errs := c.ListReposByOwner(ctx, owner)
		channel.GoForward(&wg, errs, errchan) // TODO(karlhepler): Make a Forward wrapping function that lets me set a Sprintf string to wrap over the error
		for repo := range repos {
			ownrepo := OwnerRepo{Owner: owner, Repo: Repo(*repo.Name)}
			commits, errs := c.ListOwnerRepoCommitsByDateRange(ctx, ownrepo, date)
			channel.GoForward(&wg, errs, errchan)
			channel.GoForward(&wg, commits, outchan)
		}

		wg.Wait()
	}()
	return outchan, errchan
}

func (c *Client) ListOwnerRepoCommitsByDateRange(ctx context.Context, ownrepo OwnerRepo, date DateRange) (<-chan OwnerRepoCommit, <-chan error) {
	outchan, errchan := make(chan OwnerRepoCommit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.CommitsListOptions{
			Since:       date.Since,
			Until:       date.Until,
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			owner, repo := ownrepo.OwnerStr(), ownrepo.RepoStr()
			c.Debugf("GitHub.Repositories.ListCommits(owner=%s, repo=%s, page=%d)", owner, repo, opt.Page)
			commits, res, err := c.GitHub.Repositories.ListCommits(ctx, owner, repo, opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing owner/repository commits; owner/repo=%s opt=%+v: %w", ownrepo, opt, err)
			}
			for _, commit := range commits {
				c.Debugf("\tsha=%s", *commit.SHA)
				outchan <- OwnerRepoCommit{ownrepo, commit}
			}
			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
