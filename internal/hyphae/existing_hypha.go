package hyphae

import (
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/util/fp"
)

// ExistingHypha is not EmptyHypha. *MediaHypha and *TextualHypha implement this interface.
type ExistingHypha interface {
	Hypha

	HasTextFile() bool
	TextFilePath() string
}

func FilePaths(hs ...ExistingHypha) []string {
	var paths []string
	for _, h := range hs {
		h.Lock()
		if h.HasTextFile() {
			paths = append(paths, h.TextFilePath())
		}
		if mh, ok := h.(*MediaHypha); ok {
			paths = append(paths, mh.MediaFilePath())
		}
		h.Unlock()
	}
	return paths
}

// RenameHyphaTo renames a hypha and updates its stored file paths per
// movedFiles. Paths not in movedFiles are
// kept. The actual files are not moved, move them yourself.
func RenameHyphaTo(h ExistingHypha, newName string, movedFiles map[string]string) {
	newName = util.CanonicalName(newName)
	moved := func(oldPath string) string {
		return fp.CoalesceOrEmpty(movedFiles[oldPath], oldPath)
	}
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.CanonicalName())

	switch h := h.(type) {
	case *TextualHypha:
		h.canonicalName = newName
		h.mycoFilePath = moved(h.mycoFilePath)
	case *MediaHypha:
		h.canonicalName = newName
		h.mycoFilePath = moved(h.mycoFilePath)
		h.mediaFilePath = moved(h.mediaFilePath)
	}

	byNames[h.CanonicalName()] = h
	byNamesMutex.Unlock()
	h.Unlock()
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

func WriteToMycoFile(h ExistingHypha, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(h.TextFilePath()), 0777); err != nil {
		return err
	}
	if err := os.WriteFile(h.TextFilePath(), data, 0666); err != nil {
		return err
	}
	switch h := h.(type) {
	case *MediaHypha:
		if !h.HasTextFile() {
			h.mycoFilePath = h.TextFilePath()
		}
	}
	return nil
}
