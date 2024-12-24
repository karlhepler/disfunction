package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListPatchesByOrgRepoCommits(ctx context.Context, commits <-chan OrgRepoCommit) (<-chan OrgRepoCommitPatch, <-chan error) {
	outchan, errchan := make(chan OrgRepoCommitPatch), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.ListOptions{PerPage: 100}

		for commit := range commits {
			org, repo := commit.Org.String(), commit.Repo.String()
			c.Debugf("GitHub.Repositories.GetCommit(org=%s, repo=%s, sha=%s, page=%d)", org, repo, *commit.SHA, opt.Page)
			meta, res, err := c.GitHub.Repositories.GetCommit(ctx, org, repo, *commit.SHA, opt)
			if err != nil {
				errchan <- fmt.Errorf("error getting org/repo commit; org/repo=%s sha=%s", commit.OrgRepo, *commit.SHA)
			}

			for _, file := range meta.Files {
				c.Debugf("\tfile=%s", *file.Filename)
				outchan <- OrgRepoCommitPatch{commit, *file.Patch}
			}

			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
