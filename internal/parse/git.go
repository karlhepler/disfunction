package parse

import "strings"

func MatchGitAdd(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "+")
}
