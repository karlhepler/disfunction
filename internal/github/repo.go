package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

func (c *Client) ListReposByOrg(ctx context.Context, org Org) (<-chan *github.Repository, <-chan error) {
	outchan, errchan := make(chan *github.Repository), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			org := org.String()
			c.Debugf("GitHub.Repositories.ListByOrg(org=%s, page=%d)", org, opt.Page)
			repos, res, err := c.GitHub.Repositories.ListByOrg(ctx, org, opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing repos by org, org=%s opt%+v", org, opt)
			}
			for _, repo := range repos {
				c.Debugf("\trepo=%s", *repo.FullName)
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
