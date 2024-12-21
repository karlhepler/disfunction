package disfunction

import (
	"context"
	"time"

	"github.com/karlhepler/disfunction/internal/function"
	"github.com/karlhepler/disfunction/internal/github"
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
			ListCommitsByDateRange(ctx context.Context, since, until time.Time) ([]*github.Commit, error)
			ListPatchesByCommits(context.Context, []*github.Commit) ([]*github.Patch, error)
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
	Patches []*github.Patch
}

func Random(req RandomReq, res RandomRes) {
	var ctx = req.Context
	var gh = req.Deps.GitHub

	since, until := req.Opts.Since, req.Opts.Until
	commits, err := gh.ListCommitsByDateRange(ctx, since, until)
	if err != nil {
		res.Log(ErrorLog(err))
	}

	patches, err := gh.ListPatchesByCommits(ctx, commits)
	if err != nil {
		res.Log(ErrorLog(err))
		res.Send(RandomMsg{
			Status:  StatusError,
			Message: "internal error",
		})
		return
	}

	res.Send(RandomMsg{Patches: patches})
}
