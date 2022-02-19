package shroom

import (
	"fmt"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
)

// RemoveMedia unattaches hypha and makes a history record about that.
func RemoveMedia(u *user.User, h *hyphae.MediaHypha, lc *l18n.Localizer) error {
	hop := history.
		Operation(history.TypeUnattachHypha).
		WithFilesRemoved(h.MediaFilePath()).
		WithMsg(fmt.Sprintf("Unattach ‘%s’", h.CanonicalName())).
		WithUser(u).
		Apply()

	if len(hop.Errs) > 0 {
		rejectUnattachLog(h, u, "fail")
		// FIXME: something may be wrong here
		return fmt.Errorf("Could not unattach this hypha due to internal server errors: <code>%v</code>", hop.Errs)
	}

	if h.HasTextFile() {
		hyphae.Insert(hyphae.ShrinkMediaToTextual(h))
	} else {
		hyphae.DeleteHypha(h)
	}
	return nil
}
