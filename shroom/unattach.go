package shroom

import (
	"fmt"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
)

// UnattachHypha unattaches hypha and makes a history record about that.
func UnattachHypha(u *user.User, h hyphae.Hypher, lc *l18n.Localizer) (hop *history.Op, errtitle string) {
	hop = history.Operation(history.TypeUnattachHypha)

	if errtitle, err := CanUnattach(u, h, lc); errtitle != "" {
		hop.WithErrAbort(err)
		return hop, errtitle
	}
	H := h.(*hyphae.NonEmptyHypha)

	hop.
		WithFilesRemoved(H.BinaryPath()).
		WithMsg(fmt.Sprintf("Unattach ‘%s’", h.CanonicalName())).
		WithUser(u).
		Apply()

	if len(hop.Errs) > 0 {
		rejectUnattachLog(h, u, "fail")
		// FIXME: something may be wrong here
		return hop.WithErrAbort(fmt.Errorf("Could not unattach this hypha due to internal server errors: <code>%v</code>", hop.Errs)), "Error"
	}

	if H.BinaryPath() != "" {
		H.SetBinaryPath("")
	}
	// If nothing is left of the hypha
	if H.TextPath == "" {
		hyphae.DeleteHypha(h)
	}
	return hop, ""
}
