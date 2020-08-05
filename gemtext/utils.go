package gemtext

import (
	"strings"
)

// Function that returns a function that can strip `prefix` and trim whitespace when called.
func remover(prefix string) func(string) string {
	return func(l string) string {
		return strings.TrimSpace(strings.TrimPrefix(l, prefix))
	}
}

// Remove #, ## or ### from beginning of `line`.
func removeHeadingOctothorps(line string) string {
	f := remover("#")
	return f(f(f(line)))
}

// Return a canonical representation of a hypha `name`.
func canonicalName(name string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))
}
