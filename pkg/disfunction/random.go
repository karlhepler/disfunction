package disfunction

import (
	"context"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/function"
)

type Random struct {
	GitHub interface {
		ListOwnerCommitsByDateRange(ctx context.Context, owner string, since, until time.Time, outchan chan<- *github.RepositoryCommit, errchan chan<- error)
		ListPatchesByCommits(ctx context.Context, commits []*github.RepositoryCommit, outchan chan<- string, errchan chan<- error)
	}
}

type RandomReq struct {
	context.Context
	Owner string
	Since time.Time
	Until time.Time
	Kinds []function.Kind
}

type RandomRes interface {
	Log(string)
	Send(RandomMsg)
}

type RandomMsg struct {
	Status
	Message string
	Patches []string
}

func (r *Random) Handle(req RandomReq, res RandomRes) {
	var ctx = req.Context

	commitChan, errchan := make(chan *github.RepositoryCommit), make(chan error)
	go r.GitHub.ListOwnerCommitsByDateRange(ctx, req.Owner, req.Since, req.Until, commitChan, errchan)
	go func() {
		for err := range errchan {
			res.Log(ErrorLog(err))
			res.Send(RandomMsg{
				Status:  StatusError,
				Message: "internal error",
			})
		}
	}()

	var commits []*github.RepositoryCommit
	for commit := range commitChan {
		commits = append(commits, commit)
	}

	patchChan, errchan := make(chan string), make(chan error)
	go r.GitHub.ListPatchesByCommits(ctx, commits, patchChan, errchan)
	go func() {
		for err := range errchan {
			res.Log(ErrorLog(err))
			res.Send(RandomMsg{
				Status:  StatusError,
				Message: "internal error",
			})
		}
	}()

	for patch := range patchChan {
		res.Send(RandomMsg{Message: patch})
	}
}
