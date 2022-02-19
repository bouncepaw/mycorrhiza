package shroom

import (
	"fmt"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
)

// UnattachHypha unattaches hypha and makes a history record about that.
func UnattachHypha(u *user.User, h hyphae.Hypha, lc *l18n.Localizer) error {

	if err := CanUnattach(u, h, lc); err != nil {
		return err
	}
	H := h.(*hyphae.MediaHypha)

	hop := history.
		Operation(history.TypeUnattachHypha).
		WithFilesRemoved(H.MediaFilePath()).
		WithMsg(fmt.Sprintf("Unattach ‘%s’", h.CanonicalName())).
		WithUser(u).
		Apply()

	if len(hop.Errs) > 0 {
		rejectUnattachLog(h, u, "fail")
		// FIXME: something may be wrong here
		return fmt.Errorf("Could not unattach this hypha due to internal server errors: <code>%v</code>", hop.Errs)
	}

	if H.HasTextFile() {
		hyphae.Insert(hyphae.ShrinkMediaToTextual(H))
	} else {
		hyphae.DeleteHypha(H)
	}
	return nil
}
