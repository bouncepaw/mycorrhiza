package markup

import (
	"unicode"
)

// MatchesHorizontalLine checks if the string can be interpreted as suitable for rendering as <hr/>.
//
// The rule is: if there are more than 4 characters "-" in the string, then make it a horizontal line.
// Otherwise it is a paragraph (<p>).
func MatchesHorizontalLine(line string) bool {
	counter := 0

	// Check initially that the symbol is "-". If it is not a "-", it is most likely a space or another character.
	// With unicode.IsLetter() we can separate spaces and characters.
	for _, ch := range line {
		if ch == '-' {
			counter++
			continue
		}
		// If we bump into any other character (letter) in the line, it is immediately an incorrect horizontal line.
		// There is no point in counting further, we end the loop.
		if unicode.IsLetter(ch) {
			counter = 0
			break
		}
	}

	if counter >= 4 {
		return true
	}

	return false
}
