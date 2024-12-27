package github

import (
	"context"
	"fmt"

	"github.com/karlhepler/disfunction/internal/channel"
)

func (c *Client) ListPatchesByCommits(ctx context.Context, commits <-chan Commit) (<-chan Patch, <-chan error) {
	outchan, errchan := make(chan Patch), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		channel.ForEach(ctx, commits, func(commit Commit) {
			c.log.Debugf("*github.Client.Repositories.GetCommit(owner=%s, repo=%s, sha=%s)", commit.Repo.Owner, commit.Repo.Name, *commit.SHA)
			meta, res, err := c.gh.Repositories.GetCommit(ctx, commit.Repo.Owner.String(), commit.Repo.Name, *commit.SHA, nil)
			if err != nil {
				errchan <- fmt.Errorf("error getting commit; repo=%s sha=%s", commit.Repo, *commit.SHA)
			}

			for _, file := range meta.Files {
				if file.Patch != nil {
					c.log.Debugf("\tfile=%s", *file.Filename)
					outchan <- Patch{commit, *file.Patch}
				}
			}

			if res == nil {
				return
			}
		})
	}()
	return outchan, errchan
}
