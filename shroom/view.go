package shroom

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"io/ioutil"
	"os"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/util"
)

// FetchTextPart tries to read text file of the given hypha. If there is no file, empty string is returned.
func FetchTextPart(h *hyphae.Hypha) (string, error) {
	if h.TextPath == "" {
		return "", nil
	}
	text, err := ioutil.ReadFile(h.TextPath)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return string(text), nil
}

func SetHeaderLinks() {
	if userLinksHypha := hyphae.ByName(cfg.HeaderLinksHypha); !userLinksHypha.Exists {
		util.SetDefaultHeaderLinks()
	} else {
		contents, err := ioutil.ReadFile(userLinksHypha.TextPath)
		if err != nil || len(contents) == 0 {
			util.SetDefaultHeaderLinks()
		} else {
			text := string(contents)
			util.ParseHeaderLinks(text, markup.Rocketlink)
		}
	}
}
