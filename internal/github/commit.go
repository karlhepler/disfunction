package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
)

type RepositoryCommit struct {
	*github.RepositoryCommit
	Owner string
	Repo  string
}

func (c *Client) ListOwnerCommitsByDateRange(ctx context.Context, owner string, since, until time.Time, outchan chan<- MetaRepositoryCommit, errchan chan<- error) {
	defer close(outchan)
	defer close(errchan)

	repoChan, repoErrChan := make(chan *github.Repository), make(chan error)
	go c.ListOwnerRepos(ctx, owner, repoChan, repoErrChan)
	go channel.Forward(repoErrChan, errchan)

	for repo := range repoChan {
		commitChan, commitErrChan := make(chan *github.RepositoryCommit), make(chan error)
		go c.ListRepoCommitsByDateRange(ctx, owner, *repo.Name, since, until, commitChan, commitErrChan)
		go channel.Forward(commitErrChan, errchan)
		for commit := range commitChan {
			outchan <- RepositoryCommit{commit, owner, *repo.Name}
		}
	}
}

func (c *Client) ListRepoCommitsByDateRange(ctx context.Context, owner, repo string, since, until time.Time, outchan chan<- *github.RepositoryCommit, errchan chan<- error) {
	defer close(outchan)
	defer close(errchan)

	opt := &github.CommitsListOptions{
		Since:       since,
		Until:       until,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		commits, res, err := c.GitHub.Repositories.ListCommits(ctx, owner, repo, opt)
		if err != nil {
			errchan <- fmt.Errorf("error listing repository commits, owner=%s repo=%s opt%+v: %w", owner, repo, opt, err)
		}
		for _, commit := range commits {
			outchan <- commit
		}
		if res == nil || res.NextPage == 0 {
			break
		}
		opt.Page = res.NextPage
	}
}
