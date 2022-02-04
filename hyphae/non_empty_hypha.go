// Package hyphae is for the NonEmptyHypha type, hypha storage and stuff like that. It shall not depend on mycorrhiza modules other than util.
package hyphae

import (
	"path/filepath"
	"sync"

	"github.com/bouncepaw/mycorrhiza/files"
)

// NonEmptyHypha keeps vital information about a media hypha
type NonEmptyHypha struct {
	sync.RWMutex

	name       string // Canonical name
	Exists     bool
	TextPath   string // == "" => no text part
	binaryPath string // == "" => no attachment
}

func (h *NonEmptyHypha) SetName(s string) { h.name = s }

func (h *NonEmptyHypha) BinaryPath() string     { return h.binaryPath }
func (h *NonEmptyHypha) SetBinaryPath(s string) { h.binaryPath = s }

func (h *NonEmptyHypha) CanonicalName() string {
	return h.name
}

func (h *NonEmptyHypha) Kind() HyphaKind { // sic!
	if h.HasAttachment() {
		return HyphaMedia
	}
	return HyphaText
}

func (h *NonEmptyHypha) HasTextPart() bool {
	return h.TextPath != ""
}

// TextPartPath returns rooted path to the file where the text part should be.
func (h *NonEmptyHypha) TextPartPath() string {
	if h.TextPath == "" {
		return filepath.Join(files.HyphaeDir(), h.name+".myco")
	}
	return h.TextPath
}

// HasAttachment is true if the hypha has an attachment.
func (h *NonEmptyHypha) HasAttachment() bool {
	return h.binaryPath != ""
}
