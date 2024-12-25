package disfunction

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/karlhepler/disfunction/internal/channel"
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
		err := forEachAddedFunc(patch.Patch, matchGoFunc, func(line string) {
			fmt.Println(line)
		})

		if err != nil {
			hdl.HandleErr(err)
		}
	})

	wg.Wait()
}

func (hdl *RandomHandler) HandleErr(err error) {
	hdl.Log.Error(err)
}

func forEachLine(s string, cb func(string)) error {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		cb(scanner.Text())
	}
	return scanner.Err()
}

func forEachAddedLine(s string, cb func(string)) error {
	return forEachLine(s, func(line string) {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "+") {
			cb(line[1:])
		}
	})
}

func forEachAddedFunc(s string, match func(string) bool, cb func(string)) error {
	return forEachAddedLine(s, func(line string) {
		line = strings.TrimSpace(line)
		if match(line) {
			cb(line)
		}
	})
}

func matchGoFunc(line string) bool {
	return strings.HasPrefix(line, "func ")
}
