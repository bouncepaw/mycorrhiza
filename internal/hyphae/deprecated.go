package hyphae

import (
	"os"
)

// FetchMycomarkupFile tries to read text file of the given hypha. If there is no file, empty string is returned.
//
// TODO: Get rid of this function.
func FetchMycomarkupFile(h Hypha) (string, error) {
	switch h := h.(type) {
	case *EmptyHypha:
		return "", nil
	case *MediaHypha:
		if !h.HasTextFile() {
			return "", nil
		}
		text, err := os.ReadFile(h.TextFilePath())
		if os.IsNotExist(err) {
			return "", nil
		} else if err != nil {
			return "", err
		}
		return string(text), nil
	case *TextualHypha:
		text, err := os.ReadFile(h.TextFilePath())
		if os.IsNotExist(err) {
			return "", nil
		} else if err != nil {
			return "", err
		}
		return string(text), nil
	}
	panic("unreachable")
}
