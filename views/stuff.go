package views

import (
	"fmt"
)

// inputSize returns the size in chars for an html input element
// (as a string) given the placeholder string.
func inputSize(placeholder string) string {
	charCount := len(placeholder)
	// Because size="0" is invalid, clamp above 1
	min := 1
	if charCount <= min {
		charCount = min
	}
	return fmt.Sprint(charCount)
}
