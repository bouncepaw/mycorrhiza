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

// Hypha is a hypha you know and love.
type Hypha interface {
	sync.Locker

	CanonicalName() string
}

// DeleteHypha deletes the hypha from the storage.
func DeleteHypha(h ExistingHypha) {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.CanonicalName())
	decrementCount()
	byNamesMutex.Unlock()
	h.Unlock()
}

// Insert inserts the hypha into the storage, possibly overwriting the previous hypha with the same name. Count incrementation is done if needed. You cannot insert an empty hypha.
func Insert(h ExistingHypha) (madeNewRecord bool) {
	_, recorded := byNames[h.CanonicalName()]

	byNamesMutex.Lock()
	byNames[h.CanonicalName()] = h
	byNamesMutex.Unlock()

	if !recorded {
		incrementCount()
	}

	return !recorded
}

// ByName returns a hypha by name. It returns an *EmptyHypha if there is no such hypha. This function is the only source of empty hyphae.
func ByName(hyphaName string) (h Hypha) {
	byNamesMutex.Lock()
	defer byNamesMutex.Unlock()
	h, recorded := byNames[hyphaName]
	if recorded {
		return h
	}
	return &EmptyHypha{
		canonicalName: hyphaName,
	}
}
