package disfunction

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/github"
	"github.com/karlhepler/disfunction/internal/log"
	"github.com/karlhepler/disfunction/internal/parse"
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

func (hdl *RandomHandler) Handle(req RandomReq, res RandomRes) {
	ctx := req.Context
	owner := github.Owner(req.Owner)
	var wg sync.WaitGroup

	// list all commits
	commits, errs := hdl.GitHub.ListCommits(ctx,
		github.FilterCommitsByOwner(owner),
		github.ListCommitsSince(req.Since),
		github.ListCommitsUntil(req.Until),
	)
	channel.GoForEach(&wg, errs, hdl.HandleErr)

	// list all patches from commits
	patches, errs := hdl.GitHub.ListPatchesByCommits(ctx, commits)
	channel.GoForEach(&wg, errs, hdl.HandleErr)

	// list all new function declarations for all patches
	channel.ForEach(patches, func(patch github.Patch) {
		var onMatch = func(line string) {
			fmt.Println(line)
		}

		if err := parse.ForEachLineMatch(
			patch.Patch, onMatch,
			parse.MatchGitAdd, parse.MatchGoFunc,
		); err != nil {
			hdl.HandleErr(err)
		}
	})

	wg.Wait()
}

func (hdl *RandomHandler) HandleErr(err error) {
	hdl.Log.Error(err)
}
