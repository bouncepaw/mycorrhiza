// Package hyphae is for the MediaHypha type, hypha storage and stuff like that. It shall not depend on mycorrhiza modules other than util.
package hyphae

import (
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

func (h *MediaHypha) Kind() HyphaKind { // sic!
	if h.HasAttachment() {
		return HyphaMedia
	}
	return HyphaText
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
