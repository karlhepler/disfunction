package parse

import (
	"context"

	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/github"
)

type GoFunc struct {
	Line string
	github.Patch
}

// This is a ChannelProcessor
func ListAddedGoFuncsByPatches(ctx context.Context, patches <-chan github.Patch) (<-chan GoFunc, <-chan error) {
	return channel.Async(func(outchan chan GoFunc, errchan chan error) {
		channel.ForEach(ctx, patches, func(patch github.Patch) {
			var onLineMatch = func(line string) {
				outchan <- GoFunc{line, patch}
			}

			if err := ForEachLineMatch(
				patch.Patch, onLineMatch,
				MatchAll(MatchGitAdd, MatchGoFunc),
			); err != nil {
				errchan <- err
			}
		})
	})
}
