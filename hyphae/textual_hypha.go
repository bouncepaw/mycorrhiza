package hyphae

import (
	"sync"
)

// TextualHypha is a hypha with text, and nothing else. An article, a note, a poem, whatnot.
type TextualHypha struct {
	sync.RWMutex

	canonicalName string
	mycoFilePath  string
}

func (t *TextualHypha) CanonicalName() string {
	return t.canonicalName
}

func (t *TextualHypha) HasTextFile() bool {
	return true
}

func (t *TextualHypha) TextFilePath() string {
	return t.mycoFilePath
}

// ExtendTextualToMedia returns a new media hypha with the same name and text file as the given textual hypha. The new hypha is not stored yet.
func ExtendTextualToMedia(t *TextualHypha, mediaFilePath string) *MediaHypha {
	return &MediaHypha{
		canonicalName: t.CanonicalName(),
		mycoFilePath:  t.TextFilePath(),
		mediaFilePath: mediaFilePath,
	}
}
