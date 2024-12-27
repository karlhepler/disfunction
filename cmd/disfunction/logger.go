package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type ConsoleLogger bool

func NewConsoleLogger(debugModeEnabled bool) ConsoleLogger {
	return ConsoleLogger(debugModeEnabled)
}

func (debugModeEnabled ConsoleLogger) Debug(as ...any) {
	if console := &debugModeEnabled; debugModeEnabled {
		console.debug(as...)
	}
}

func (console ConsoleLogger) debug(a ...any) {
	var format string
	for range a {
		format += "%+v\n"
	}

	console.debugf(format, a...)
}

func (debugModeEnabled ConsoleLogger) Debugf(format string, a ...any) {
	if console := &debugModeEnabled; debugModeEnabled {
		console.debugf(format, a...)
	}
}

func (console ConsoleLogger) debugf(format string, a ...any) {
	format = "[ DEBUG ] " + format
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	log.Printf(format, a...)
}

func (console ConsoleLogger) Error(err error) {
	log.Printf("[ ERROR ] %s\n", err)
	err = errors.Unwrap(err)
	for i := 1; err != nil; i++ {
		fmt.Fprintf(os.Stderr, "\t%d: %s\n", i, err)
		err = errors.Unwrap(err)
	}
}
