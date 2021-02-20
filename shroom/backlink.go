package shroom

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/link"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/util"
)

func BackLinkEntriesHTML(h *hyphae.Hypha) (html string) {
	for _, backlinkHypha := range h.BackLinks {
		_ = link.Link{}
		html += fmt.Sprintf(`<li class="backlinks__entry">
			<a class="backlinks__link" href="/hypha/%s">%s</a>`, backlinkHypha.Name, util.BeautifulName(backlinkHypha.Name))
	}
	return
}

// FindAllBacklinks iterates over all hyphae that have text parts, sets their outlinks and then sets backlinks.
func FindAllBacklinks() {
	for h := range hyphae.FilterTextHyphae(hyphae.YieldExistingHyphae()) {
		findBacklinkWorker(h)
	}
}

func findBacklinkWorker(h *hyphae.Hypha) {
	var (
		wg                sync.WaitGroup
		textContents, err = ioutil.ReadFile(h.TextPath)
	)
	if err == nil {
		for outlink := range markup.Doc(h.Name, string(textContents)).OutLinks() {
			go func() {
				wg.Add(1)
				outlinkHypha := hyphae.ByName(outlink)
				if outlinkHypha == h {
					return
				}

				outlinkHypha.AddBackLink(h)
				outlinkHypha.InsertIfNewKeepExistence()
				h.AddOutLink(outlinkHypha)
				wg.Done()
			}()
		}
		wg.Wait()
	} else {
		log.Println("Error when reading text contents of ‘%s’: %s", h.Name, err.Error())
	}
}
