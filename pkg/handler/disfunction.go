package handler

import (
	"context"
	"sync"
	"time"

	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/github"
	"github.com/karlhepler/disfunction/internal/log"
	"github.com/karlhepler/disfunction/internal/parse"
)

// Disfunction one possible implementation enabling this use case:
// TODO(karlhepler): generate this comment on each handler.
type Disfunction struct {
	gh  *github.Client
	log log.Logger
}

func NewDisfunction(ghtoken string, log log.Logger) (*Disfunction, error) {
	gh, err := github.NewClient(ghtoken, log)
	return &Disfunction{gh, log}, err
}

type DisfunctionReq struct {
	Owner string
	Since time.Time
	Until time.Time
	Ctx   context.Context
}

func (req DisfunctionReq) Context() context.Context {
	return req.Ctx
}

type DisfunctionRes struct {
	Patch github.Patch
	Ctx   context.Context
}

func (res DisfunctionRes) Context() context.Context {
	return res.Ctx
}

// Disfunction.Handle is a usecase.Handler
//
// TODO(karlhepler): automatically run something to generate these comments
// every time an implementation of it is found. Preferrably on save. I wonder
// if there is a way to tap into go fmt to do this. That would be automatic and
// would also work for most developers.
func (hdl *Disfunction) Handle(req DisfunctionReq, res Sender[DisfunctionRes]) {
	ctx := req.Context()
	owner := github.Owner(req.Owner)
	var wg sync.WaitGroup

	// list all commits
	commits, errs := hdl.gh.ListCommits(ctx,
		github.FilterCommitsByOwner(owner),
		github.ListCommitsSince(req.Since),
		github.ListCommitsUntil(req.Until),
	)
	channel.GoForEach(ctx, &wg, errs, hdl.log.Error)

	// list all patches from commits
	patches, errs := hdl.gh.ListPatchesByCommits(ctx, commits)
	channel.GoForEach(ctx, &wg, errs, hdl.log.Error)

	// list all new function declarations for all patches
	gofuncs, errs := parse.ListAddedGoFuncsByPatches(ctx, patches)
	channel.GoForEach(ctx, &wg, errs, hdl.log.Error)
	channel.ForEach(ctx, gofuncs, func(gofunc parse.GoFunc) {
		//
	})

	wg.Wait()
}
