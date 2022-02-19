package shroom

import (
	"errors"
	"os"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
)

// FetchTextFile tries to read text file of the given hypha. If there is no file, empty string is returned.
func FetchTextFile(h hyphae.Hypha) (string, error) {
	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		return "", errors.New("empty hyphae have no text")
	case *hyphae.MediaHypha:
		if !h.HasTextFile() {
			return "", nil
		}
		text, err := os.ReadFile(h.TextFilePath())
		if os.IsNotExist(err) {
			return "", nil
		} else if err != nil {
			return "", err
		}
		return string(text), nil
	case *hyphae.TextualHypha:
		text, err := os.ReadFile(h.TextFilePath())
		if os.IsNotExist(err) {
			return "", nil
		} else if err != nil {
			return "", err
		}
		return string(text), nil
	}
	panic("unreachable")
}

// SetHeaderLinks initializes header links by reading the configured hypha, if there is any, or resorting to default values.
func SetHeaderLinks() {
	switch userLinksHypha := hyphae.ByName(cfg.HeaderLinksHypha).(type) {
	case *hyphae.EmptyHypha:
		cfg.SetDefaultHeaderLinks()
	case hyphae.ExistingHypha:
		contents, err := os.ReadFile(userLinksHypha.TextFilePath())
		if err != nil || len(contents) == 0 {
			cfg.SetDefaultHeaderLinks()
		} else {
			text := string(contents)
			cfg.ParseHeaderLinks(text)
		}
	}
}
