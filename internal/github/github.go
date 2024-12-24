package github

import (
	"fmt"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/lithammer/dedent"
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

type OwnerRepoCommitPatch struct {
	OwnerRepoCommit
	Patch string
}

func (orcp OwnerRepoCommitPatch) String() string {
	return fmt.Sprintf(dedent.Dedent(`
		---
		Owner/Repo: %s
		Author: %s
		Commit SHA: %s
		---
		%s
	`), orcp.OwnerRepo, *orcp.Author.Login, *orcp.SHA, orcp.Patch)
}
