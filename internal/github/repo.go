package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/funk"
)

type Repository = github.Repository

func NewRepository(ownerRepo string) *Repository {
	repo := new(Repository)

	// default to repo name
	if !strings.Contains(ownerRepo, "/") {
		ownerRepo = "/" + ownerRepo
	}

	parts := strings.Split(ownerRepo, "/")
	if ownerLogin := parts[0]; ownerLogin != "" {
		repo.Owner = &github.User{
			Login: &ownerLogin,
		}
	}
	if len(parts) > 1 {
		if repoName := parts[1]; repoName != "" {
			repo.Name = &repoName
		}
	}

	return repo
}

type RepoAllowList []*github.Repository

func (allowList RepoAllowList) Allows(ghRepo *Repository) bool {
	for _, repo := range allowList {
		// Effectively: */*
		if repo.Owner == nil && repo.Name == nil {
			return true
		}

		// Effectively: owner/*
		if repo.Owner != nil && repo.Name == nil {
			match := *repo.Owner.Login == *ghRepo.Owner.Login
			if match == true {
				return true
			}

			continue // next!
		}

		// Effectively: */repo
		if repo.Owner == nil && repo.Name != nil {
			match := *repo.Name == *ghRepo.Name
			if match == true {
				return true
			}

			continue // next!
		}

		// Effectively: owner/repo
		if repo.Owner != nil && repo.Name != nil {
			match := *repo.Owner.Login == *ghRepo.Owner.Login && *repo.Name == *ghRepo.Name
			if match == true {
				return true
			}

			continue // next! (this is the end. for now...)
		}
	}

	return false // I don't know if this should be true or false, so I just picked one at random.
}

type listReposConfig struct {
	allowlist RepoAllowList
}

func ListReposOwnedBy(user *github.User) funk.Option[listReposConfig] {
	return func(config *listReposConfig) {
		config.allowlist = append(config.allowlist, &Repository{Owner: user})
	}
}

func ListReposAllowedBy(allowlist []*Repository) funk.Option[listReposConfig] {
	return func(config *listReposConfig) {
		config.allowlist = allowlist
	}
}

func (c *Client) ListRepos(ctx context.Context, opts ...funk.Option[listReposConfig]) (<-chan *Repository, <-chan error) {
	var config = funk.ConfigWithOptions[listReposConfig](opts)
	return channel.Async(func(outchan chan *Repository, errchan chan error) {
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
				if config.allowlist.Allows(repo) {
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
