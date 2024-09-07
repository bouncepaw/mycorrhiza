package shroom

import (
	"fmt"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
)

// RemoveMedia removes media from the media hypha and makes a history record about that. If it only had media, the hypha will be deleted. If it also had text, the hypha will become textual.
func RemoveMedia(u *user.User, h *hyphae.MediaHypha) error {
	hop := history.
		Operation(history.TypeRemoveMedia).
		WithFilesRemoved(h.MediaFilePath()).
		WithMsg(fmt.Sprintf("Remove media from ‘%s’", h.CanonicalName())).
		WithUser(u).
		Apply()

	if len(hop.Errs) > 0 {
		rejectRemoveMediaLog(h, u, "fail")
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
