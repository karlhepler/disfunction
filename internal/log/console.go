package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Console bool

func NewConsole(debugModeEnabled bool) Console {
	return Console(debugModeEnabled)
}

func (debugModeEnabled Console) Debug(as ...any) {
	if console := &debugModeEnabled; debugModeEnabled {
		console.debug(as...)
	}
}

func (console Console) debug(a ...any) {
	var format string
	for range a {
		format += "%+v\n"
	}

	console.debugf(format, a...)
}

func (debugModeEnabled Console) Debugf(format string, a ...any) {
	if console := &debugModeEnabled; debugModeEnabled {
		console.debugf(format, a...)
	}
}

func (console Console) debugf(format string, a ...any) {
	format = "[ DEBUG ] " + format
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	log.Printf(format, a...)
}

func (console Console) Error(err error) {
	log.Printf("[ ERROR ] %s\n", err)
	err = errors.Unwrap(err)
	for i := 1; err != nil; i++ {
		fmt.Fprintf(os.Stderr, "\t%d: %s\n", i, err)
		err = errors.Unwrap(err)
	}
}
