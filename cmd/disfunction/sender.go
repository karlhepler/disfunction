package main

import (
	"fmt"

	"github.com/karlhepler/disfunction/pkg/api"
	"github.com/lithammer/dedent"
)

type ConsoleSender struct{}

func NewConsoleSender() ConsoleSender {
	return ConsoleSender{}
}

func (console ConsoleSender) Send(res api.DisfunctionRes) {
	gofunc := res.GoFunc

	fmt.Printf(dedent.Dedent(`
			Commit Author: %s
			Commit URL: %s
			Function: %s
		`),
		gofunc.RepositoryCommit.GetCommit().GetAuthor().GetName(),
		gofunc.RepositoryCommit.GetHTMLURL(),
		gofunc.Line,
	)
}
