package github

import (
	"context"
	"time"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListOwnerCommitsByDateRange(ctx context.Context, owner string, since, until time.Time, outchan chan<- *github.RepositoryCommit, errchan chan<- error) {
	defer close(outchan)
	defer close(errchan)

	rchan, rechan := make(chan *github.Repository), make(chan error)
	go c.ListOwnerRepositories(ctx, owner, rchan, rechan)
	go func() {
		for err := range rechan {
			errchan <- err
		}
	}()

	for repo := range rchan {
		commits, err := c.ListRepoCommitsWithinDateRange(ctx, *repo.Name, since, until)
		if err != nil {
			c.Log.ErrorWithMessagef(err, "repo=%s; since=%+v; until=%+v", *repo.Name, since, until)
		}
		for _, commit := range commits {
			list = append(list, &Commit{RepositoryCommit: commit, Owner: org, Repo: *repo.Name})
		}
	}
}
