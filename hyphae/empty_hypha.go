package hyphae

import "sync"

type EmptyHypha struct {
	sync.RWMutex

	canonicalName string
}

func (e *EmptyHypha) CanonicalName() string {
	return e.canonicalName
}

func (e *EmptyHypha) DoesExist() bool {
	return false
}

func (e *EmptyHypha) HasTextPart() bool {
	return false
}

func (e *EmptyHypha) TextPartPath() string {
	return ""
}

// NewEmptyHypha returns an empty hypha struct with given name.
func NewEmptyHypha(hyphaName string) *EmptyHypha {
	return &EmptyHypha{
		canonicalName: hyphaName,
	}
}

func FillEmptyHyphaUpToMediaHypha(e *EmptyHypha) *MediaHypha { // sic!
	return &MediaHypha{
		name:       e.CanonicalName(),
		TextPath:   "",
		binaryPath: "",
	}
}
