package hyphae

import "sync"

type EmptyHypha struct {
	sync.RWMutex

	canonicalName string
}

func (e *EmptyHypha) CanonicalName() string {
	return e.canonicalName
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

func FillEmptyHyphaUpToTextualHypha(e *EmptyHypha, textPath string) *NonEmptyHypha { // sic!
	return &NonEmptyHypha{
		name:     e.CanonicalName(),
		TextPath: textPath,
	}
}

func FillEmptyHyphaUpToMediaHypha(e *EmptyHypha, binaryPath string) *NonEmptyHypha { // sic!
	return &NonEmptyHypha{
		name:       e.CanonicalName(),
		binaryPath: binaryPath,
	}
}
