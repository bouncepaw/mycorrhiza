package hyphae

import "sync"

type HyphaKind int

const (
	HyphaEmpty HyphaKind = iota
	HyphaText
	HyphaMedia
)

// Hypher is a temporary name for this interface. The name will become Hypha, once the struct with the said name is deprecated for good.
type Hypher interface {
	sync.Locker

	CanonicalName() string
	Kind() HyphaKind

	HasTextPart() bool
	TextPartPath() string
}

// DeleteHypha deletes the hypha from the storage.
func DeleteHypha(h Hypher) {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.CanonicalName())
	decrementCount()
	byNamesMutex.Unlock()
	h.Unlock()
}
