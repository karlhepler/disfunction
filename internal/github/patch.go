package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListPatchesByCommits(ctx context.Context, commits <-chan Commit) (<-chan Patch, <-chan error) {
	outchan, errchan := make(chan Patch), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.ListOptions{PerPage: 100}

		for commit := range commits {
			c.Log.Debugf("GitHub.Repositories.GetCommit(owner=%s, repo=%s, sha=%s, page=%d)", commit.Repo.Owner, commit.Repo.Name, *commit.SHA, opt.Page)
			meta, res, err := c.GitHub.Repositories.GetCommit(ctx, commit.Repo.Owner.String(), commit.Repo.Name, *commit.SHA, opt)
			if err != nil {
				errchan <- fmt.Errorf("error getting commit; repo=%s sha=%s", commit.Repo, *commit.SHA)
			}

			for _, file := range meta.Files {
				c.Log.Debugf("\tfile=%s", *file.Filename)
				outchan <- Patch{commit, *file.Patch}
			}

			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
