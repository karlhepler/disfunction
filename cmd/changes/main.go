package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/github"
	"github.com/karlhepler/disfunction/internal/log"
)

func main() {
	var ctx = context.Background()
	console := log.NewConsole(true)

	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		console.Error(errors.New("missing required env var GITHUB_TOKEN"))
		os.Exit(1)
	}

	gh, err := github.NewClient(token, console)
	if err != nil {
		console.Error(err)
		os.Exit(2)
	}

	mustVal := newMustVal[time.Time](console)
	repo := github.NewRepository("paxosglobal/pax")
	commits, errs := gh.ListCommitsByRepo(ctx, repo,
		github.ListCommitsByRepoSince(mustVal(time.Parse(time.DateOnly, "2024-01-01"))),
		github.ListCommitsByRepoUntil(mustVal(time.Parse(time.DateOnly, "2025-01-01"))),
	)

	var wg sync.WaitGroup
	channel.GoForEach(ctx, &wg, errs, func(err error) {
		console.Error(err)
	})

	numCommits := 0
	channel.ForEach(ctx, commits, func(commit github.Commit) {
		numCommits++
	})
	fmt.Println("Changes per Year = ", numCommits)
	fmt.Println("Changes per Day = ", numCommits/366.0)

	wg.Wait()
}

func newMustVal[T any](log log.Logger) func(T, error) T {
	return func(val T, err error) T {
		if err != nil {
			log.Error(err)
			os.Exit(3)
		}
		return val
	}
}
