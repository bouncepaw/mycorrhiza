package hyphae

import "sync"

type EmptyHypha struct {
	sync.RWMutex

	canonicalName string
}

func (e *EmptyHypha) CanonicalName() string {
	return e.canonicalName
}

func ExtendEmptyToTextual(e *EmptyHypha, mycoFilePath string) *TextualHypha {
	return &TextualHypha{
		canonicalName: e.CanonicalName(),
		mycoFilePath:  mycoFilePath,
	}
}

func ExtendEmptyToMedia(e *EmptyHypha, mediaFilePath string) *MediaHypha {
	return &MediaHypha{
		canonicalName: e.CanonicalName(),
		mycoFilePath:  "",
		mediaFilePath: mediaFilePath,
	}
}
