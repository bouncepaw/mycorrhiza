package hyphae

import (
	"sync"
)

type Hypha struct {
	sync.RWMutex

	Name       string
	Exists     bool
	TextPath   string
	BinaryPath string
	OutLinks   []*Hypha
	BackLinks  []*Hypha
}

/*
// Insert inserts the hypha into the mycelium. It overwrites the previous record, if there was any, and returns false. If the was no previous record, return true.
func (h *Hypha) Insert() (justCreated bool) {
	var hp *Hypha
	hp, justCreated = ByName(h.Name)

	mycm.Lock()
	defer mycm.Unlock()
	if justCreated {
		mycm.byNames[hp.Name] = h
	} else {
		hp = h
	}

	return justCreated
}*/

// PhaseOut marks the hypha as non-existent. This is an idempotent operation.
func (h *Hypha) PhaseOut() {
	h.Lock()
	h.Exists = false
	h.OutLinks = make([]*Hypha, 0)
	h.TextPath = ""
	h.BinaryPath = ""
	h.Unlock()
}
