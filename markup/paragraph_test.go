package markup

import (
	"fmt"
	"testing"
)

/*
func TestGetTextNode(t *testing.T) {
	tests := [][]string{
		// input   textNode  rest
		{"barab", "barab", ""},
		{"test, ", "test", ", "},
		{"/test/", "", "/test/"},
		{"\\/test/", "/test", "/"},
		{"test \\/ar", "test /ar", ""},
		{"test //italian// test", "test ", "//italian// test"},
	}
	for _, triplet := range tests {
		a, b := getTextNode([]byte(triplet[0]))
		if a != triplet[1] || string(b) != triplet[2] {
			t.Error(fmt.Sprintf("Wanted: %q\nGot: %q %q", triplet, a, b))
		}
	}
}
*/

func TestParagraphToHtml(t *testing.T) {
	tests := [][]string{
		{"a simple paragraph", "a simple paragraph"},
		{"//italic//", "<em>italic</em>"},
		{"Embedded //italic//", "Embedded <em>italic</em>"},
		{"double //italian// //text//", "double <em>italian</em> <em>text</em>"},
		{"it has `mono`", "it has <code>mono</code>"},
		{"this is a left **bold", "this is a left <strong>bold</strong>"},
		{"this line has a ,comma, two of them", "this line has a ,comma, two of them"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
	}
	for _, test := range tests {
		if ParagraphToHtml(test[0]) != test[1] {
			t.Error(fmt.Sprintf("%q: Wanted %q, got %q", test[0], test[1], ParagraphToHtml(test[0])))
		}
	}
}
