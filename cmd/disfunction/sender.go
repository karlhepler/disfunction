package main

import (
	"fmt"

	"github.com/karlhepler/disfunction/pkg/handler"
	"github.com/lithammer/dedent"
)

type ConsoleSender struct{}

func NewConsoleSender() ConsoleSender {
	return ConsoleSender{}
}

func (console ConsoleSender) Send(res handler.DisfunctionRes) {
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
