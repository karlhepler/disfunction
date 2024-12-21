package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListPatchesByCommits(ctx context.Context, commits []*github.RepositoryCommit, outchan chan<- string, errchan chan<- error) {
	defer close(outchan)
	defer close(errchan)

	opt := &github.ListOptions{PerPage: 100}
	for _, commit := range commits {
		// TODO(karlhepler): replace hardcoded owner and repo with variables
		commit, res, err := c.GitHub.Repositories.GetCommit(ctx, "karlhepler", "disfunction", *commit.SHA, opt)
		if err != nil {
			errchan <- fmt.Errorf("error getting repository commit for owner=karlhepler repo=disfunction sha=" + *commit.SHA)
		}
		if res == nil || res.NextPage == 0 {
			break
		}
		opt.Page = res.NextPage

		// add the commit patches to the patches
		for _, file := range commit.Files {
			outchan <- *file.Patch
		}
	}
}
