package hyphae

import "sync"

type HyphaKind int

const (
	HyphaEmpty HyphaKind = iota
	HyphaText
	HyphaMedia
)

// Hypher is a temporary name for this interface. The name will become MediaHypha, once the struct with the said name is deprecated for good.
type Hypher interface {
	sync.Locker

	CanonicalName() string
	Kind() HyphaKind
	DoesExist() bool

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

// RenameHyphaTo renames a hypha and performs respective changes in the storage.
func RenameHyphaTo(h Hypher, newName string) {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.CanonicalName())
	h.(*MediaHypha).SetName(newName)
	byNames[h.CanonicalName()] = h.(*MediaHypha)
	byNamesMutex.Unlock()
	h.Unlock()
}
