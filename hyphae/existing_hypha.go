package hyphae

import (
	"github.com/bouncepaw/mycorrhiza/util"
)

// ExistingHypha is not EmptyHypha.
type ExistingHypha interface {
	Hypha

	// DoesExist does nothing except marks that the type is an ExistingHypha.
	DoesExist()

	HasTextFile() bool
	TextFilePath() string
}

// RenameHyphaTo renames a hypha and renames stored filepaths as needed. The actual files are not moved, move them yourself.
func RenameHyphaTo(h ExistingHypha, newName string, replaceName func(string) string) {
	// TODO: that replaceName is suspicious.
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
