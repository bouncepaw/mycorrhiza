package hyphae

import (
	"regexp"
	"strings"
	"sync"
)

// HyphaPattern is a pattern which all hyphae names must match.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"'&%{}]+`)

// IsValidName checks for invalid characters and path traversals.
func IsValidName(hyphaName string) bool {
	if !HyphaPattern.MatchString(hyphaName) {
		return false
	}
	for _, segment := range strings.Split(hyphaName, "/") {
		if segment == ".git" || segment == ".." {
			return false
		}
	}
	return true
}

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

// insert inserts the hypha into the storage. A previous record is used if possible. Count incrementation is done if needed.
func insert(h Hypher) (madeNewRecord bool) {
	hp, recorded := byNames[h.CanonicalName()]
	if recorded {
		hp.(*MediaHypha).mergeIn(h)
	} else {
		storeHypha(h)
		incrementCount()
	}

	return !recorded
}

// InsertIfNew checks whether hypha exists and returns `true` if it didn't and has been created.
func InsertIfNew(h Hypher) (madeNewRecord bool) {
	if h.DoesExist() {
		return false
	}
	return insert(h)
}

func storeHypha(h Hypher) {
	byNamesMutex.Lock()
	byNames[h.CanonicalName()] = h
	byNamesMutex.Unlock()

	h.Lock()
	h.(*MediaHypha).Exists = true
	h.Unlock()
}

// ByName returns a hypha by name. It may have been recorded to the storage.
func ByName(hyphaName string) (h Hypher) {
	h, recorded := byNames[hyphaName]
	if recorded {
		return h
	}
	return NewEmptyHypha(hyphaName)
}
