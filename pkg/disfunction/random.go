package disfunction

import (
	"context"
	"time"

	"github.com/karlhepler/disfunction/internal/function"
	"github.com/karlhepler/disfunction/internal/github"
)

type Random struct {
	GitHub *github.Client
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
	SendErr(error)
	SendErrChan(<-chan error)
}

type RandomMsg struct {
	Status
	Message string
	Patches []string
}

func (r *Random) Handle(req RandomReq, res RandomRes) {
	ctx := req.Context
	owner := github.Owner(req.Owner)
	date := github.DateRange{
		Since: req.Since,
		Until: req.Until,
	}

	commits, errs := r.GitHub.ListOwnerCommitsByDateRange(ctx, owner, date)
	go r.handleErrs(errs, res)

	patches, errs := r.GitHub.ListPatchesByOwnerRepoCommits(ctx, commits)
	go r.handleErrs(errs, res)

	for patch := range patches {
		res.Send(RandomMsg{
			Message: patch.String(),
		})
	}
}

func (r *Random) handleErrs(errs <-chan error, res RandomRes) {
	for err := range errs {
		res.Log(ErrorLog(err))
		res.Send(RandomMsg{
			Status:  StatusError,
			Message: "internal error",
		})
	}
}
