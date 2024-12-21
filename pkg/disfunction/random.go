package disfunction

import (
	"context"
	"time"

	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/function"
)

type RandomReq struct {
	context.Context
	Opts struct {
		Since time.Time
		Until time.Time
		Kinds []function.Kind
	}
	Deps struct {
		GitHub interface {
			ListOwnerCommitsByDateRange(ctx context.Context, owner string, since, until time.Time) ([]*github.Commit, error)
			ListPatchesByCommits(context.Context, []*github.Commit, chan<- string, chan<- error)
		}
	}
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

func Random(req RandomReq, res RandomRes) {
	var ctx = req.Context
	var gh = req.Deps.GitHub

	since, until := req.Opts.Since, req.Opts.Until
	commits, err := gh.ListCommitsByDateRange(ctx, since, until)
	if err != nil {
		res.Log(ErrorLog(err))
	}

	outchan, errchan := make(chan string), make(chan error)
	go gh.ListPatchesByCommits(ctx, commits, outchan, errchan)
	go func() {
		for err := range errchan {
			res.Log(ErrorLog(err))
			res.Send(RandomMsg{
				Status:  StatusError,
				Message: "internal error",
			})
		}
	}()

	for patch := range outchan {
		res.Send(RandomMsg{Message: patch})
	}
}
