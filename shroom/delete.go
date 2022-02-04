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
func DeleteHypha(u *user.User, h hyphae.Hypher, lc *l18n.Localizer) (hop *history.Op, errtitle string) {
	hop = history.Operation(history.TypeDeleteHypha)

	if errtitle, err := CanDelete(u, h, lc); errtitle != "" {
		hop.WithErrAbort(err)
		return hop, errtitle
	}

	originalText, _ := FetchTextPart(h)
	hop.
		WithFilesRemoved(h.TextPartPath(), h.(*hyphae.NonEmptyHypha).BinaryPath()).
		WithMsg(fmt.Sprintf("Delete ‘%s’", h.CanonicalName())).
		WithUser(u).
		Apply()
	if !hop.HasErrors() {
		backlinks.UpdateBacklinksAfterDelete(h, originalText)
		hyphae.DeleteHypha(h)
	}
	return hop, ""
}
