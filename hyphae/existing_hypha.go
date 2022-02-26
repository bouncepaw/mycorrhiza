package hyphae

import (
	"github.com/bouncepaw/mycorrhiza/util"
	"os"
	"path/filepath"
)

// ExistingHypha is not EmptyHypha. *MediaHypha and *TextualHypha implement this interface.
type ExistingHypha interface {
	Hypha

	HasTextFile() bool
	TextFilePath() string
}

// RenameHyphaTo renames a hypha and renames stored filepaths as needed. The actual files are not moved, move them yourself.
func RenameHyphaTo(h ExistingHypha, newName string, replaceName func(string) string) {
	// TODO: that replaceName is suspicious, get rid of it.
	newName = util.CanonicalName(newName)
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.CanonicalName())

	switch h := h.(type) {
	case *TextualHypha:
		h.canonicalName = newName
		h.mycoFilePath = replaceName(h.mycoFilePath)
	case *MediaHypha:
		h.canonicalName = newName
		h.mycoFilePath = replaceName(h.mediaFilePath)
		h.mediaFilePath = replaceName(h.mediaFilePath)
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
