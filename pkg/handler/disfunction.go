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

// TODO(karlhepler):devex Also have a way to call the most fequently called functions.
// That shouldn't be too hard. I'll just have to save state somewhere.
// I'll need to consume all data in a timeline and present the most current
// state, along with the ability to see more. That data could be used to power
// a lsp to give that data to you in your editor. Get a list of functions
// sorted by popularity. Also put a special effect on them in the editor view.
// There should also be a hover window over ctrl+k that gives the custom info
// (wrapping the normal info).
//
// ChatGPT:
// The golang.org/x/tools module provides helpful libraries for building tools
// related to Go, including packages like golang.org/x/tools/lsp that you can
// reference or extend. Use Go’s standard libraries for handling JSON-RPC and
// HTTP/2, as they are required for LSP communication.

// Disfunction one possible implementation enabling this use case:
// TODO(karlhepler):devex generate this comment on each handler.
type Disfunction struct {
	gh  *github.Client
	log log.Logger
}

func NewDisfunction(ghtoken string, log log.Logger) (*Disfunction, error) {
	gh, err := github.NewClient(ghtoken, log)
	return &Disfunction{gh, log}, err
}

type DisfunctionReq struct {
	Ctx          context.Context
	Since        time.Time
	Until        time.Time
	AllowedRepos github.RepoAllowList
	AllowedFiles github.FileAllowList
}

func (req DisfunctionReq) Context() context.Context {
	return req.Ctx
}

type DisfunctionRes struct {
	Ctx    context.Context
	GoFunc parse.GoFunc
}

func (res DisfunctionRes) Context() context.Context {
	return res.Ctx
}

func (r DisfunctionRes) Send(res DisfunctionRes) {
	//
}

// Disfunction.Handle is a usecase.Handler
//
// TODO(karlhepler):devex automatically run something to generate these comments
// every time an implementation of it is found. Preferrably on save. I wonder
// if there is a way to tap into go fmt to do this. That would be automatic and
// would also work for most developers.
func (hdl *Disfunction) Handle(req DisfunctionReq, res Sender[DisfunctionRes]) {
	var wg sync.WaitGroup

	var ctx = req.Context()

	commits, errs := hdl.gh.ListCommits(ctx,
		github.ListCommitsSince(req.Since),
		github.ListCommitsUntil(req.Until),
		github.ListCommitsExclusiveTo(req.AllowedRepos),
		github.ListCommitsWithDetail(true),
		github.ListCommitsToFiles(req.AllowedFiles),
	)
	channel.GoForEach(ctx, &wg, errs, hdl.log.Error)

	gofuncs, errs := channel.Map(ctx, commits, func(commit github.Commit, outs chan<- parse.GoFunc, errs chan<- error) {
		for _, file := range commit.Files {
			if file.Patch == nil {
				continue
			}

			err := parse.ForEachLineMatch(
				*file.Patch,
				parse.MatchAll(parse.MatchGitAdd, parse.MatchGoFunc),
				func(line string) {
					outs <- parse.GoFunc{
						Line:   line,
						Commit: commit,
					}
				},
			)

			if err != nil {
				errs <- err
			}
		}
	})
	channel.GoForEach(ctx, &wg, errs, hdl.log.Error)

	channel.ForEach(ctx, gofuncs, func(gofunc parse.GoFunc) {
		res.Send(DisfunctionRes{ctx, gofunc})
	})

	wg.Wait()
}
