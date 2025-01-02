package github

import (
	"context"
	"fmt"

	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/lithammer/dedent"
)

type Patch struct {
	Commit
	Patch string
}

func (patch Patch) String() string {
	// TODO(karlhepler): Work on this to collect the infomration I need
	return fmt.Sprintf(dedent.Dedent(`
		---
		Repo: %s
		Author: %s
		Commit SHA: %s
		---
		%s
	`), patch.Repo, *patch.Author.Login, *patch.SHA, patch.Patch)
}

func (c *Client) ListPatchesByCommits(ctx context.Context, commits <-chan Commit) (<-chan Patch, <-chan error) {
	return channel.Async(func(outchan chan Patch, errchan chan error) {
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
	})
}
