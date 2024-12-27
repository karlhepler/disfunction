package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
)

type listReposConfig struct {
	owner Owner
}

type listReposOption func(*listReposConfig)

func FilterReposByOwner(owner Owner) listReposOption {
	return func(config *listReposConfig) {
		config.owner = owner
	}
}

func (c *Client) ListRepos(ctx context.Context, opts ...listReposOption) (<-chan *github.Repository, <-chan error) {
	var config listReposConfig
	for _, opt := range opts {
		opt(&config)
	}

	outchan, errchan := make(chan *github.Repository), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

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
	}()
	return outchan, errchan
}
