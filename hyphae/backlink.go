package hyphae

import (
	"fmt"
	"io/ioutil"

	"github.com/bouncepaw/mycorrhiza/link"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/util"
)

func (h *Hypha) BackLinkEntriesHTML() (html string) {
	for _, backlinkHypha := range h.BackLinks {
		_ = link.Link{}
		html += fmt.Sprintf(`<li class="backlinks__entry">
			<a class="backlinks__link" href="/hypha/%s">%s</a>`, backlinkHypha.Name, util.BeautifulName(backlinkHypha.Name))
	}
	return
}

func (h *Hypha) outlinksThis(oh *Hypha) bool {
	for _, outlink := range h.OutLinks {
		if outlink == oh {
			return true
		}
	}
	return false
}

func (h *Hypha) backlinkedBy(oh *Hypha) bool {
	for _, backlink := range h.BackLinks {
		if backlink == oh {
			return true
		}
	}
	return false
}

// FindAllBacklinks iterates over all hyphae that have text parts, sets their outlinks and then sets backlinks.
func FindAllBacklinks() {
	for h := range FilterTextHyphae(YieldExistingHyphae()) {
		findBacklinkWorker(h)
	}
}

func findBacklinkWorker(h *Hypha) {
	h.Lock()
	defer h.Unlock()

	textContents, err := ioutil.ReadFile(h.TextPath)
	if err == nil {
		for outlink := range markup.Doc(h.Name, string(textContents)).OutLinks() {
			outlink := outlink
			outlinkHypha := ByName(outlink)
			if outlinkHypha == h {
				continue
			}

			outlinkHypha.Lock()
			if !outlinkHypha.backlinkedBy(h) {
				outlinkHypha.BackLinks = append(outlinkHypha.BackLinks, h)
				outlinkHypha.InsertIfNewKeepExistence()
			}
			outlinkHypha.Unlock()

			// Insert outlinkHypha if unique
			if !h.outlinksThis(outlinkHypha) {
				h.OutLinks = append(h.OutLinks, outlinkHypha)
			}
		}
	}
}
