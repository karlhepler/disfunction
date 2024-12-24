package disfunction

import (
	"context"
	"time"

	"github.com/karlhepler/disfunction/internal/github"
	"github.com/karlhepler/disfunction/internal/log"
)

type RandomHandler struct {
	GitHub *github.Client
	Log    log.Logger
}

func NewRandomHandler(ghtoken string, log log.Logger) (*RandomHandler, error) {
	gh, err := github.NewClient(ghtoken, log)
	return &RandomHandler{GitHub: gh, Log: log}, err
}

type RandomReq struct {
	context.Context
	Org   string
	Since time.Time
	Until time.Time
}

type RandomRes interface {
	Send(RandomMsg)
}

type RandomMsg struct {
	Patch github.OrgRepoCommitPatch
}

func (hdl *RandomHandler) Handle(req RandomReq, res RandomRes) error {
	ctx := req.Context
	org := github.Org(req.Org)
	date := github.DateRange{
		Since: req.Since,
		Until: req.Until,
	}

	commits, errs := hdl.GitHub.ListOrgCommitsByDateRange(ctx, org, date)
	go hdl.HandleErrs(errs, res)

	patches, errs := hdl.GitHub.ListPatchesByOrgRepoCommits(ctx, commits)
	go hdl.HandleErrs(errs, res)

	for patch := range patches {
		res.Send(RandomMsg{Patch: patch})
	}

	return nil
}

func (hdl *RandomHandler) HandleErrs(errs <-chan error, res RandomRes) {
	for err := range errs {
		hdl.HandleErr(err)
	}
}

func (hdl *RandomHandler) HandleErr(err error) {
	hdl.Log.Error(err)
}
