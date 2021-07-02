package shroom

import (
	"os"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
)

// FetchTextPart tries to read text file of the given hypha. If there is no file, empty string is returned.
func FetchTextPart(h *hyphae.Hypha) (string, error) {
	if h.TextPath == "" {
		return "", nil
	}
	text, err := os.ReadFile(h.TextPath)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return string(text), nil
}

func SetHeaderLinks() {
	if userLinksHypha := hyphae.ByName(cfg.HeaderLinksHypha); !userLinksHypha.Exists {
		cfg.SetDefaultHeaderLinks()
	} else {
		contents, err := os.ReadFile(userLinksHypha.TextPath)
		if err != nil || len(contents) == 0 {
			cfg.SetDefaultHeaderLinks()
		} else {
			text := string(contents)
			cfg.ParseHeaderLinks(text)
		}
	}
}
