package hyphae

import (
	"sync"
)

type TextualHypha struct {
	sync.RWMutex

	canonicalName string
	mycoFilePath  string
}

func (t *TextualHypha) DoesExist() {
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

func ExtendTextualToMedia(t *TextualHypha, mediaFilePath string) *MediaHypha {
	return &MediaHypha{
		canonicalName: t.CanonicalName(),
		mycoFilePath:  t.TextFilePath(),
		mediaFilePath: mediaFilePath,
	}
}
