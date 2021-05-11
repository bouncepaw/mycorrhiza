package markup

import (
	"strings"
)

// Function that returns a function that can strip `prefix` and trim whitespace when called.
func remover(prefix string) func(string) string {
	return func(l string) string {
		return strings.TrimSpace(strings.TrimPrefix(l, prefix))
	}
}
