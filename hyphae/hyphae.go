// Package hyphae is for the MediaHypha type, hypha storage and stuff like that. It shall not depend on mycorrhiza modules other than util.
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

var byNames = make(map[string]Hypher)
var byNamesMutex = sync.Mutex{}

// EmptyHypha returns an empty hypha struct with given name.
func EmptyHypha(hyphaName string) *MediaHypha {
	return &MediaHypha{
		name:       hyphaName,
		Exists:     false,
		TextPath:   "",
		binaryPath: "",
	}
}

// ByName returns a hypha by name. It may have been recorded to the storage.
func ByName(hyphaName string) (h Hypher) {
	h, recorded := byNames[hyphaName]
	if recorded {
		return h
	}
	return EmptyHypha(hyphaName)
}

func storeHypha(h Hypher) {
	byNamesMutex.Lock()
	byNames[h.CanonicalName()] = h
	byNamesMutex.Unlock()

	h.Lock()
	h.(*MediaHypha).Exists = true
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
