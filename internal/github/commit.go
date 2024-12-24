package github

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
)

func (c *Client) ListOrgCommitsByDateRange(ctx context.Context, org Org, date DateRange) (<-chan OrgRepoCommit, <-chan error) {
	outchan, errchan := make(chan OrgRepoCommit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)
		var wg sync.WaitGroup

		repos, errs := c.ListReposByOrg(ctx, org)
		channel.GoForward(&wg, errs, errchan) // TODO(karlhepler): Make a Forward wrapping function that lets me set a Sprintf string to wrap over the error
		for repo := range repos {
			ownrepo := OrgRepo{Org: org, Repo: Repo(*repo.Name)}
			commits, errs := c.ListOrgRepoCommitsByDateRange(ctx, ownrepo, date)
			channel.GoForward(&wg, errs, errchan)
			channel.GoForward(&wg, commits, outchan)
		}

		wg.Wait()
	}()
	return outchan, errchan
}

func (c *Client) ListOrgRepoCommitsByDateRange(ctx context.Context, ownrepo OrgRepo, date DateRange) (<-chan OrgRepoCommit, <-chan error) {
	outchan, errchan := make(chan OrgRepoCommit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.CommitsListOptions{
			Since:       date.Since,
			Until:       date.Until,
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			org, repo := ownrepo.Org.String(), ownrepo.Repo.String()
			c.Debugf("GitHub.Repositories.ListCommits(org=%s, repo=%s, page=%d)", org, repo, opt.Page)
			commits, res, err := c.GitHub.Repositories.ListCommits(ctx, org, repo, opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing org/repo commits; org/repo=%s opt=%+v: %w", ownrepo, opt, err)
			}
			for _, commit := range commits {
				c.Debugf("\tsha=%s", *commit.SHA)
				outchan <- OrgRepoCommit{ownrepo, commit}
			}
			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
