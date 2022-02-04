// Package hyphae is for the MediaHypha type, hypha storage and stuff like that. It shall not depend on mycorrhiza modules other than util.
package hyphae

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/bouncepaw/mycorrhiza/files"
)

// MediaHypha keeps vital information about a media hypha
type MediaHypha struct {
	sync.RWMutex

	name       string // Canonical name
	Exists     bool
	TextPath   string // == "" => no text part
	binaryPath string // == "" => no attachment
}

func (h *MediaHypha) SetName(s string) { h.name = s }

func (h *MediaHypha) BinaryPath() string     { return h.binaryPath }
func (h *MediaHypha) SetBinaryPath(s string) { h.binaryPath = s }

func (h *MediaHypha) CanonicalName() string {
	return h.name
}

func (h *MediaHypha) Kind() HyphaKind {
	if !h.DoesExist() {
		return HyphaEmpty
	}
	if h.HasAttachment() {
		return HyphaMedia
	}
	return HyphaText
}

func (h *MediaHypha) DoesExist() bool { // TODO: rename
	return h.Exists
}

func (h *MediaHypha) HasTextPart() bool {
	return h.TextPath != ""
}

// TextPartPath returns rooted path to the file where the text part should be.
func (h *MediaHypha) TextPartPath() string {
	if h.TextPath == "" {
		return filepath.Join(files.HyphaeDir(), h.name+".myco")
	}
	return h.TextPath
}

// HasAttachment is true if the hypha has an attachment.
func (h *MediaHypha) HasAttachment() bool {
	return h.binaryPath != ""
}

// mergeIn merges in content file paths from a different hypha object. Prints warnings sometimes.
func (h *MediaHypha) mergeIn(oh Hypher) {
	if h == oh {
		return
	}
	h.Lock()
	if h.TextPath == "" && oh.HasTextPart() {
		h.TextPath = oh.TextPartPath()
	}
	if oh := oh.(*MediaHypha); oh.Kind() == HyphaMedia {
		if h.binaryPath != "" {
			log.Println("There is a file collision for attachment of a hypha:", h.binaryPath, "and", oh.binaryPath, "-- going on with the latter")
		}
		h.binaryPath = oh.binaryPath
	}
	h.Unlock()
}
