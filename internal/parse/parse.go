package parse

import (
	"bufio"
	"strings"
)

type LineMatcherFunc func(line string) bool

func ForEachLine(s string, onLine func(string)) error {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		onLine(scanner.Text())
	}
	return scanner.Err()
}

func ForEachLineMatch(s string, match MatcherFunc[string], onMatch func(string)) error {
	return ForEachLine(s, func(line string) {
		if match(line) {
			onMatch(line)
		}
	})
}

type MatcherFunc[T any] func(T) bool

// All matchers must match.
// If any one doesn't match, then it doesn't match.
// If there are no matchers defined, then it matches everything.
func MatchAll[T any](matchers ...MatcherFunc[T]) MatcherFunc[T] {
	return func(t T) bool {
		var isMatch = true

		for _, match := range matchers {
			isMatch = isMatch && match(t)
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
func MatchOne[T any](matchers ...MatcherFunc[T]) MatcherFunc[T] {
	return func(t T) bool {
		var isMatch = false

		for _, match := range matchers {
			isMatch = isMatch || match(t)
			if isMatch {
				break
			}
		}

		return len(matchers) == 0 || isMatch
	}
}
