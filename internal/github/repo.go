package github

import (
	"context"
	"fmt"
	"slices"

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

func FilterReposByRepos(repos []Repo) funk.Option[listReposConfig] {
	return func(config *listReposConfig) {
		config.repos = repos
	}
}

func (c *Client) ListRepos(ctx context.Context, opts ...funk.Option[listReposConfig]) (<-chan *github.Repository, <-chan error) {
	var config = funk.ConfigWithOptions[listReposConfig](opts)
	return channel.Async(func(outchan chan *github.Repository, errchan chan error) {
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
				if !isOwnerMatch(config.owner, repo.Owner) {
					continue
				}
				if !isRepoMatch(config.repos, repo) {
					continue
				}
				c.log.Debugf("\trepo=%s", *repo.FullName)
				outchan <- repo
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

func isRepoMatch(allowlist []Repo, candidate *github.Repository) bool {
	if len(allowlist) == 0 {
		return true
	}
	if candidate == nil {
		return false
	}
	if candidate.Name == nil {
		return false
	}
	return slices.ContainsFunc(allowlist, func(repo Repo) bool {
		return repo.Name == *candidate.Name
	})
}
