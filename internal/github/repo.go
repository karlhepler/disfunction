package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/funk"
)

type Repo struct {
	Owner
	Name string
}

func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

type listReposConfig struct {
	owner Owner
	repos []Repo
}

func FilterReposByOwner(owner Owner) funk.Option[listReposConfig] {
	return func(config *listReposConfig) {
		config.owner = owner
	}
}

func (c *Client) ListRepos(ctx context.Context, opts ...funk.Option[listReposConfig]) (<-chan Repo, <-chan error) {
	var config = funk.ConfigWithOptions[listReposConfig](opts)
	return channel.Async(func(outchan chan Repo, errchan chan error) {
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
				if isOwnerMatch(config.owner, repo.Owner) {
					c.log.Debugf("\trepo=%s", *repo.FullName)
					outchan <- Repo{
						Owner: Owner{Login: *repo.Owner.Login},
						Name:  *repo.Name,
					}
				}
			}

			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	})
}

func isOwnerMatch(reference Owner, candidate *github.User) bool {
	if reference.Login == "" {
		return true
	}
	if candidate == nil {
		return false
	}
	if candidate.Login == nil {
		return false
	}
	return reference.Login == *candidate.Login
}
