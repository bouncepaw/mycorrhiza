package shroom

import (
	"io/ioutil"
	"log"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/markup"
)

// FindAllBacklinks iterates over all hyphae that have text parts, sets their outlinks and then sets backlinks.
func FindAllBacklinks() {
	for h := range hyphae.FilterTextHyphae(hyphae.YieldExistingHyphae()) {
		findBacklinkWorker(h)
	}
}

func findBacklinkWorker(h *hyphae.Hypha) {
	var (
		textContents, err = ioutil.ReadFile(h.TextPath)
	)
	if err == nil {
		for outlink := range markup.Doc(h.Name, string(textContents)).OutLinks() {
			outlinkHypha := hyphae.ByName(outlink)
			if outlinkHypha == h {
				break
			}

			outlinkHypha.AddBackLink(h)
			outlinkHypha.InsertIfNewKeepExistence()
			h.AddOutLink(outlinkHypha)
		}
	} else {
		log.Println("Error when reading text contents of ‘%s’: %s", h.Name, err.Error())
	}
}
