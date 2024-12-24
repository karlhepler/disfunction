package github

import (
	"fmt"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/lithammer/dedent"
)

type Org string

func (o Org) String() string {
	return string(o)
}

type Repo string

func (r Repo) String() string {
	return string(r)
}

type OrgRepo struct {
	Org
	Repo
}

func (or OrgRepo) String() string {
	return fmt.Sprintf("%s/%s", or.Org, or.Repo)
}

type OrgRepoCommit struct {
	OrgRepo
	*github.RepositoryCommit
}

type DateRange struct {
	Since time.Time
	Until time.Time
}

type RepositoryCommit github.RepositoryCommit

type OrgRepoCommitPatch struct {
	OrgRepoCommit
	Patch string
}

func (orcp OrgRepoCommitPatch) String() string {
	return fmt.Sprintf(dedent.Dedent(`
		---
		Org/Repo: %s
		Author: %s
		Commit SHA: %s
		---
		%s
	`), orcp.OrgRepo, *orcp.Author.Login, *orcp.SHA, orcp.Patch)
}
