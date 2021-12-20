// Package hyphae is for the Hypha type, hypha storage and stuff like that. It shall not depend on mycorrhiza modules other than util.
package hyphae

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/bouncepaw/mycorrhiza/files"
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

// Hypha keeps vital information about a hypha
type Hypha struct {
	sync.RWMutex

	Name       string // Canonical name
	Exists     bool
	TextPath   string // == "" => no text part
	BinaryPath string // == "" => no attachment
}

// TextPartPath returns rooted path to the file where the text part should be.
func (h *Hypha) TextPartPath() string {
	if h.TextPath == "" {
		return filepath.Join(files.HyphaeDir(), h.Name+".myco")
	}
	return h.TextPath
}

// HasAttachment is true if the hypha has an attachment.
func (h *Hypha) HasAttachment() bool {
	return h.BinaryPath != ""
}

var byNames = make(map[string]*Hypha)
var byNamesMutex = sync.Mutex{}

// EmptyHypha returns an empty hypha struct with given name.
func EmptyHypha(hyphaName string) *Hypha {
	return &Hypha{
		Name:       hyphaName,
		Exists:     false,
		TextPath:   "",
		BinaryPath: "",
	}
}

// ByName returns a hypha by name. It may have been recorded to the storage.
func ByName(hyphaName string) (h *Hypha) {
	h, recorded := byNames[hyphaName]
	if recorded {
		return h
	}
	return EmptyHypha(hyphaName)
}

func storeHypha(h *Hypha) {
	byNamesMutex.Lock()
	byNames[h.Name] = h
	byNamesMutex.Unlock()

	h.Lock()
	h.Exists = true
	h.Unlock()
}

// insert inserts the hypha into the storage. A previous record is used if possible. Count incrementation is done if needed.
func (h *Hypha) insert() (justRecorded bool) {
	hp, recorded := byNames[h.Name]
	if recorded {
		hp.mergeIn(h)
	} else {
		storeHypha(h)
		incrementCount()
	}

	return !recorded
}

// InsertIfNew checks whether hypha exists and returns `true` if it didn't and has been created.
func (h *Hypha) InsertIfNew() (justRecorded bool) {
	if !h.Exists {
		return h.insert()
	}
	return false
}

// Delete removes a hypha from the storage.
func (h *Hypha) Delete() {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.Name)
	decrementCount()
	byNamesMutex.Unlock()
	h.Unlock()
}

// RenameTo renames a hypha and performs respective changes in the storage.
func (h *Hypha) RenameTo(newName string) {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.Name)
	h.Name = newName
	byNames[h.Name] = h
	byNamesMutex.Unlock()
	h.Unlock()
}

// mergeIn merges in content file paths from a different hypha object. Prints warnings sometimes.
func (h *Hypha) mergeIn(oh *Hypha) {
	if h == oh {
		return
	}
	h.Lock()
	if h.TextPath == "" && oh.TextPath != "" {
		h.TextPath = oh.TextPath
	}
	if oh.BinaryPath != "" {
		if h.BinaryPath != "" {
			log.Println("There is a file collision for attachment of a hypha:", h.BinaryPath, "and", oh.BinaryPath, "-- going on with the latter")
		}
		h.BinaryPath = oh.BinaryPath
	}
	h.Unlock()
}
