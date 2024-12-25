package parse

import (
	"bufio"
	"strings"
)

type LineMatcher func(line string) bool

func ForEachLine(s string, onLine func(string)) error {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		onLine(scanner.Text())
	}
	return scanner.Err()
}

// MatchAll by default. Wrap your list of matchers in MatchOne for "OR" behavior.
func ForEachLineMatch(s string, onMatch func(string), matchers ...LineMatcher) error {
	return ForEachLine(s, func(line string) {
		if MatchAll(matchers...)(line) {
			onMatch(line[1:])
		}
	})
}

func MatchGitAdd(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "+")
}

func MatchGoFunc(line string) bool {
	return strings.Contains(line, "func ")
}

// All matchers must match.
// If any one doesn't match, then it doesn't match.
// If there are no matchers defined, then it matches everything.
func MatchAll(matchers ...LineMatcher) LineMatcher {
	return func(line string) bool {
		var isMatch = true

		for _, match := range matchers {
			isMatch = isMatch && match(line)
			if !isMatch {
				break
			}
		}

		return isMatch
	}
}

// Only one matcher must match.
// If ALL of them don't match, then it doesn't match.
// If there are no matchers defined, then it matches everything.
func MatchOne(matchers ...LineMatcher) LineMatcher {
	return func(line string) bool {
		var isMatch = false

		for _, match := range matchers {
			isMatch = isMatch || match(line)
			if isMatch {
				break
			}
		}

		return len(matchers) == 0 || isMatch
	}
}
