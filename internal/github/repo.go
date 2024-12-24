package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListReposByOwner(ctx context.Context, owner Owner) (<-chan *github.Repository, <-chan error) {
	outchan, errchan := make(chan *github.Repository), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			repos, res, err := c.GitHub.Repositories.ListByOrg(ctx, owner.String(), opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing repos by owner, owner=%s opt%+v", owner, opt)
			}
			for _, repo := range repos {
				outchan <- repo
			}
			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
