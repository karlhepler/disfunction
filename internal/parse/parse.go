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

func ForEachLineMatch(s string, onMatch func(string), matchers ...LineMatcher) error {
	return ForEachLine(s, func(line string) {
		// If there are no matchers, then match everything.
		isMatch := len(matchers) == 0

		// matchers are integrated with OR
		// if any single matcher matches, then it's a match
		for _, match := range matchers {
			isMatch = isMatch || match(line)
			if isMatch {
				break
			}
		}

		if isMatch {
			onMatch(line[1:])
		}
	})
}

func MatchGitAdd(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "+")
}

func MatchGoFunc(line string) bool {
	line = strings.TrimSpace(line)
	return strings.Contains(line, "func ")
}
