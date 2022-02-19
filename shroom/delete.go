package shroom

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
)

// DeleteHypha deletes hypha and makes a history record about that.
func DeleteHypha(u *user.User, h hyphae.Hypha, lc *l18n.Localizer) (hop *history.Op, errtitle string) {
	hop = history.
		Operation(history.TypeDeleteHypha).
		WithMsg(fmt.Sprintf("Delete ‘%s’", h.CanonicalName())).
		WithUser(u)

	if errtitle, err := CanDelete(u, h, lc); errtitle != "" {
		hop.WithErrAbort(err)
		return hop, errtitle
	}

	switch h := h.(type) {
	case *hyphae.MediaHypha:
		hop.WithFilesRemoved(h.MediaFilePath(), h.TextFilePath())
	case *hyphae.TextualHypha:
		hop.WithFilesRemoved(h.TextFilePath())
	default:
		panic("impossible")
	}
	originalText, _ := FetchTextFile(h)
	hop.Apply()
	if !hop.HasErrors() {
		backlinks.UpdateBacklinksAfterDelete(h, originalText)
		hyphae.DeleteHypha(h.(hyphae.ExistingHypha)) // we panicked before, so it's safe
	}
	return hop, ""
}
