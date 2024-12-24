package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListPatchesByOwnerRepoCommits(ctx context.Context, commits <-chan OwnerRepoCommit) (<-chan OwnerRepoCommitPatch, <-chan error) {
	outchan, errchan := make(chan OwnerRepoCommitPatch), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.ListOptions{PerPage: 100}

		for commit := range commits {
			owner, repo := commit.Owner.String(), commit.Repo.String()
			c.Debugf("GitHub.Repositories.GetCommit(owner=%s, repo=%s, sha=%s, page=%d)", owner, repo, *commit.SHA, opt.Page)
			meta, res, err := c.GitHub.Repositories.GetCommit(ctx, owner, repo, *commit.SHA, opt)
			if err != nil {
				errchan <- fmt.Errorf("error getting repository commit; owner/repo=%s sha=%s", commit.OwnerRepo, *commit.SHA)
			}

			for _, file := range meta.Files {
				c.Debugf("\tfile=%s", *file.Filename)
				outchan <- OwnerRepoCommitPatch{commit, *file.Patch}
			}

			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
