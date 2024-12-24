package disfunction

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/karlhepler/disfunction/internal/github"
)

type RandomHandler struct {
	GitHub *github.Client
}

func NewRandomHandler(ghtoken string) (*RandomHandler, error) {
	gh, err := github.NewClient(ghtoken)

	// TODO(karlhepler): this should somehow be defined in an outer layer
	gh.Debugf = func(format string, a ...any) {
		format = "[DEBUG] " + format
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		log.Printf(format, a...)
	}

	return &RandomHandler{GitHub: gh}, err
}

type RandomReq struct {
	context.Context
	Org   string
	Since time.Time
	Until time.Time
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
		res.Send(RandomMsg{
			Message: patch.String(),
		})
	}

	return nil
}

func (hdl *RandomHandler) HandleErrs(errs <-chan error, res RandomRes) {
	for err := range errs {
		hdl.HandleErr(err, res)
	}
}

func (hdl *RandomHandler) HandleErr(err error, res RandomRes) {
	res.Log(ErrorLog(err))
	res.Send(RandomMsg{
		Status:  StatusError,
		Message: "internal error",
	})
}
