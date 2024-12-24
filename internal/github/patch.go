package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

type OwnerRepoCommitPatch struct {
	OwnerRepoCommit
	Patch string
}

func (orcp OwnerRepoCommitPatch) String() string {
	return fmt.Sprintf(`
		OwnerRepo: %s
		Commit SHA: %s
		Patch: %s
	`, orcp.OwnerRepo, *orcp.SHA, orcp.Patch)
}

func (c *Client) ListPatchesByOwnerRepoCommits(ctx context.Context, commits <-chan OwnerRepoCommit) (<-chan OwnerRepoCommitPatch, <-chan error) {
	outchan, errchan := make(chan OwnerRepoCommitPatch), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.ListOptions{PerPage: 100}

		for commit := range commits {
			meta, res, err := c.GitHub.Repositories.GetCommit(ctx, commit.Owner.String(), commit.Repo.String(), *commit.SHA, opt)
			if err != nil {
				errchan <- fmt.Errorf("error getting repository commit; ownrepo=%s sha=%s", commit.OwnerRepo, *commit.SHA)
			}

			for _, file := range meta.Files {
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
