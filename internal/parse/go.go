package parse

import (
	"context"
	"strings"

	"github.com/karlhepler/disfunction/internal/channel"
	"github.com/karlhepler/disfunction/internal/github"
)

type GoFunc struct {
	Line string
	github.Patch
}

func MatchGoFunc(line string) bool {
	return strings.Contains(line, "func ")
}

// This is a channel.ProcessorFunc. Remember that now.
// --
// NOTE(karlhepler): Looking at this a little later now, I'm not sure I agree
// with it. This is a huge jump. We're going from patches to --  not just strings --
// very specifically filtered strings.
func ListAddedGoFuncsByPatches(ctx context.Context, patches <-chan github.Patch) (<-chan GoFunc, <-chan error) {
	return channel.Async(func(outchan chan GoFunc, errchan chan error) {
		channel.ForEach(ctx, patches, func(patch github.Patch) {
			if err := ForEachLineMatch(
				patch.Patch, MatchAll[string](MatchGitAdd, MatchGoFunc),
				func(line string) {
					outchan <- GoFunc{line, patch}
				},
			); err != nil {
				errchan <- err
			}
		})
	})
}

// func MatchGoFile(commit github.Commit) bool {
// 	//
// }
