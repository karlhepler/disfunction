package disfunction

import (
	"context"
	"sync"
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
	Owner string
	Since time.Time
	Until time.Time
}

type RandomRes interface {
	Send(RandomMsg)
}

type RandomMsg struct {
	Patch github.Patch
}

func (hdl *RandomHandler) Handle(req RandomReq, res RandomRes) error {
	ctx := req.Context
	owner := github.Owner(req.Owner)
	var wg sync.WaitGroup

	commits, errs := hdl.GitHub.ListCommits(ctx,
		github.FilterCommitsByOwner(owner),
		github.ListCommitsSince(req.Since),
		github.ListCommitsUntil(req.Until),
	)
	wg.Add(1)
	go func(errs <-chan error) {
		defer wg.Done()
		for err := range errs {
			hdl.HandleErr(err)
		}
	}(errs)

	patches, errs := hdl.GitHub.ListPatchesByCommits(ctx, commits)
	wg.Add(1)
	go func(errs <-chan error) {
		defer wg.Done()
		for err := range errs {
			hdl.HandleErr(err)
		}
	}(errs)

	for patch := range patches {
		res.Send(RandomMsg{Patch: patch})
	}

	wg.Wait()

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
