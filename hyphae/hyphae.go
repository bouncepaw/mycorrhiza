// Package hyphae is for the Hypha type, hypha storage and stuff like that. It shall not depend on mycorrhiza modules other than util.
package hyphae

import (
	"log"
	"regexp"
	"sync"
)

// HyphaPattern is a pattern which all hyphae must match.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"'&%{}]+`)

type Hypha struct {
	sync.RWMutex

	Name       string // Canonical name
	Exists     bool
	TextPath   string // == "" => no text part
	BinaryPath string // == "" => no attachment
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

// Insert inserts the hypha into the storage. A previous record is used if possible. Count incrementation is done if needed.
func (h *Hypha) Insert() (justRecorded bool) {
	hp, recorded := byNames[h.Name]
	if recorded {
		hp.MergeIn(h)
	} else {
		storeHypha(h)
		IncrementCount()
	}

	return !recorded
}

func (h *Hypha) InsertIfNew() (justRecorded bool) {
	if !h.Exists {
		return h.Insert()
	}
	return false
}

func (h *Hypha) Delete() {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.Name)
	DecrementCount()
	byNamesMutex.Unlock()
	h.Unlock()
}

func (h *Hypha) RenameTo(newName string) {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.Name)
	h.Name = newName
	byNames[h.Name] = h
	byNamesMutex.Unlock()
	h.Unlock()
}

// MergeIn merges in content file paths from a different hypha object. Prints warnings sometimes.
func (h *Hypha) MergeIn(oh *Hypha) {
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
