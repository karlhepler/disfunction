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

type Repo struct {
	Owner
	Name string
}

func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

type Commit struct {
	Repo
	*github.RepositoryCommit
}

type DateRange struct {
	Since time.Time
	Until time.Time
}

type Patch struct {
	Commit
	Patch string
}

func (patch Patch) String() string {
	// TODO(karlhepler): Work on this to collect the infomration I need
	return fmt.Sprintf(dedent.Dedent(`
		---
		Repo: %s
		Author: %s
		Commit SHA: %s
		---
		%s
	`), patch.Repo, *patch.Author.Login, *patch.SHA, patch.Patch)
}
