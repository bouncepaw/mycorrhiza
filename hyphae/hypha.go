package hyphae

import (
	"regexp"
	"strings"
	"sync"
)

// hyphaNamePattern is a pattern which all hyphae names must match.
var hyphaNamePattern = regexp.MustCompile(`[^?!:#@><*|"'&%{}]+`)

// IsValidName checks for invalid characters and path traversals.
func IsValidName(hyphaName string) bool {
	if !hyphaNamePattern.MatchString(hyphaName) {
		return false
	}
	for _, segment := range strings.Split(hyphaName, "/") {
		if segment == ".git" || segment == ".." {
			return false
		}
	}
	return true
}

// Hypha is the hypha you know and love.
type Hypha interface {
	sync.Locker

	// CanonicalName returns the canonical name of the hypha.
	//
	//     util.CanonicalName(h.CanonicalName()) == h.CanonicalName()
	CanonicalName() string
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
