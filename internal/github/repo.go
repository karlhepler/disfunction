package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/funk"
)

type listReposConfig struct {
	owner Owner
}

func FilterReposByOwner(owner Owner) funk.Option[listReposConfig] {
	return func(config *listReposConfig) {
		config.owner = owner
	}
}

func (c *Client) ListRepos(ctx context.Context, opts ...funk.Option[listReposConfig]) (<-chan *github.Repository, <-chan error) {
	var config = funk.ConfigWithOptions[listReposConfig](opts...)
	return channel.Async(func(outchan chan *github.Repository, errchan chan error) {
		owner := config.owner.String()
		opt := &github.RepositoryListByAuthenticatedUserOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			c.log.Debugf("*github.Client.Repositories.ListByAuthenticatedUser(page=%d)", opt.Page)
			repos, res, err := c.gh.Repositories.ListByAuthenticatedUser(ctx, opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing repos by authenticated user; opt=%+v: %w", opt, err)
			}

			for _, repo := range repos {
				if owner == "" || owner == *repo.Owner.Login {
					c.log.Debugf("\trepo=%s", *repo.FullName)
					outchan <- repo
				}
			}

			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	})
}
