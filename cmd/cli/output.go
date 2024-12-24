package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/karlhepler/disfunction/pkg/disfunction"
)

type Output struct {
	EnableDebug bool
}

func (out Output) Debugf(format string, a ...any) {
	if !out.EnableDebug {
		return
	}

	format = "[DEBUG] " + format
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	log.Printf(format, a...)
}

func (out Output) Error(err error) {
	log.Printf("[ERROR] %w\n", err)
}

func (out Output) Send(res disfunction.RandomMsg) {
	fmt.Println(res.Patch)
}
