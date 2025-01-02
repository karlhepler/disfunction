package parse

import (
	"strings"

	"github.com/karlhepler/disfunction/internal/github"
)

type GoFunc struct {
	Line string
	github.Commit
}

func MatchGoFunc(line string) bool {
	return strings.Contains(line, "func ")
}
