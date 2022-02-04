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
	HyphaText HyphaKind = iota
	HyphaMedia
)

// Hypher is a temporary name for this interface. The name will become NonEmptyHypha, once the struct with the said name is deprecated for good.
type Hypher interface {
	sync.Locker

	CanonicalName() string

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
	h.(*NonEmptyHypha).SetName(newName)
	byNames[h.CanonicalName()] = h.(*NonEmptyHypha)
	byNamesMutex.Unlock()
	h.Unlock()
}

// insert inserts the hypha into the storage, possibly overwriting the previous hypha with the same name. Count incrementation is done if needed.
func insert(h Hypher) (madeNewRecord bool) {
	_, recorded := byNames[h.CanonicalName()]

	byNamesMutex.Lock()
	byNames[h.CanonicalName()] = h
	byNamesMutex.Unlock()

	if !recorded {
		incrementCount()
	}

	return !recorded
}

// InsertIfNew checks whether the hypha exists and returns `true` if it didn't and has been created.
func InsertIfNew(h Hypher) (madeNewRecord bool) {
	switch ByName(h.CanonicalName()).(type) {
	case *EmptyHypha:
		return insert(h)
	default:
		return false
	}
}

// ByName returns a hypha by name. It returns an *EmptyHypha if there is no such hypha. This function is the only source of empty hyphae.
func ByName(hyphaName string) (h Hypher) {
	byNamesMutex.Lock()
	defer byNamesMutex.Unlock()
	h, recorded := byNames[hyphaName]
	if recorded {
		return h
	}
	return NewEmptyHypha(hyphaName)
}
