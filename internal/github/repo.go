package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListOwnerRepos(ctx context.Context, owner string, outchan chan<- *github.Repository, errchan chan<- error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, res, err := c.GitHub.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			errchan <- fmt.Errorf("error listing repos by org, org=%s opt%+v", owner, opt)
		}

		if repos != nil {
			for _, repo := range repos {
				outchan <- repo
			}
		}

		if res == nil || res.NextPage == 0 {
			break
		}
		opt.Page = res.NextPage
	}
}
