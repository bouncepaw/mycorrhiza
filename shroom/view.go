package shroom

import (
	"errors"
	"os"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
)

// FetchTextPart tries to read text file of the given hypha. If there is no file, empty string is returned.
func FetchTextPart(h hyphae.Hypher) (string, error) {
	switch h.(type) {
	case *hyphae.EmptyHypha:
		return "", errors.New("empty hyphae have no text")
	}
	if !h.HasTextPart() {
		return "", nil
	}
	text, err := os.ReadFile(h.TextPartPath())
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return string(text), nil
}

// SetHeaderLinks initializes header links by reading the configured hypha, if there is any, or resorting to default values.
func SetHeaderLinks() {
	switch userLinksHypha := hyphae.ByName(cfg.HeaderLinksHypha).(type) {
	case *hyphae.EmptyHypha:
		cfg.SetDefaultHeaderLinks()
	default:
		contents, err := os.ReadFile(userLinksHypha.TextPartPath())
		if err != nil || len(contents) == 0 {
			cfg.SetDefaultHeaderLinks()
		} else {
			text := string(contents)
			cfg.ParseHeaderLinks(text)
		}
	}
}
