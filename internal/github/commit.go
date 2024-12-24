package github

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/channel"
)

type Owner string

func (o Owner) String() string {
	return string(o)
}

type Repo string

func (r Repo) String() string {
	return string(r)
}

type OwnerRepo struct {
	Owner
	Repo
}

func (or OwnerRepo) String() string {
	return fmt.Sprintf("%s/%s", or.Owner, or.Repo)
}
func (or OwnerRepo) OwnerStr() string {
	return or.Owner.String()
}
func (or OwnerRepo) RepoStr() string {
	return or.Repo.String()
}

type OwnerRepoCommit struct {
	OwnerRepo
	*github.RepositoryCommit
}
type DateRange struct {
	Since time.Time
	Until time.Time
}
type RepositoryCommit github.RepositoryCommit

func (c *Client) ListOwnerCommitsByDateRange(ctx context.Context, owner Owner, date DateRange) (<-chan OwnerRepoCommit, <-chan error) {
	outchan, errchan := make(chan OwnerRepoCommit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		repos, errs := c.ListOwnerRepos(ctx, owner)
		go channel.Forward(errs, errchan) // TODO(karlhepler): Make a Forward wrapping function that lets me set a Sprintf string to wrap over the error
		for repo := range repos {
			ownrepo := OwnerRepo{Owner: owner, Repo: Repo(*repo.Name)}
			commits, errs := c.ListOwnerRepoCommitsByDateRange(ctx, ownrepo, date)
			go channel.Forward(errs, errchan)
			channel.Forward(commits, outchan)
		}

		for err := range errchan {
			log.Println(err)
		}
	}()
	return outchan, errchan
}

func (c *Client) ListOwnerRepoCommitsByDateRange(ctx context.Context, ownrepo OwnerRepo, date DateRange) (<-chan OwnerRepoCommit, <-chan error) {
	outchan, errchan := make(chan OwnerRepoCommit), make(chan error)
	go func() {
		defer close(outchan)
		defer close(errchan)

		opt := &github.CommitsListOptions{
			Since:       date.Since,
			Until:       date.Until,
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			commits, res, err := c.GitHub.Repositories.ListCommits(ctx, ownrepo.OwnerStr(), ownrepo.RepoStr(), opt)
			if err != nil {
				errchan <- fmt.Errorf("error listing owner/repository commits; ownrepo=%s opt=%+v: %w", ownrepo, opt, err)
			}
			for _, commit := range commits {
				outchan <- OwnerRepoCommit{ownrepo, commit}
			}
			if res == nil || res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}()
	return outchan, errchan
}
