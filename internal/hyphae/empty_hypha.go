package hyphae

import "sync"

// EmptyHypha is a hypha that does not exist and is not stored anywhere. You get one when querying for a hypha that was not created before.
type EmptyHypha struct {
	sync.RWMutex

	canonicalName string
}

func (e *EmptyHypha) CanonicalName() string {
	return e.canonicalName
}

// ExtendEmptyToTextual returns a new textual hypha with the same name as the given empty hypha. The created hypha is not stored yet.
func ExtendEmptyToTextual(e *EmptyHypha, mycoFilePath string) *TextualHypha {
	return &TextualHypha{
		canonicalName: e.CanonicalName(),
		mycoFilePath:  mycoFilePath,
	}
}

// ExtendEmptyToMedia returns a new media hypha with the same name as the given empty hypha. The created hypha is not stored yet.
func ExtendEmptyToMedia(e *EmptyHypha, mediaFilePath string) *MediaHypha {
	return &MediaHypha{
		canonicalName: e.CanonicalName(),
		mycoFilePath:  "",
		mediaFilePath: mediaFilePath,
	}
}
