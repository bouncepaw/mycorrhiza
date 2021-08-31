package shroom

import (
	"fmt"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
)

// DeleteHypha deletes hypha and makes a history record about that.
func DeleteHypha(u *user.User, h *hyphae.Hypha) (hop *history.HistoryOp, errtitle string) {
	hop = history.Operation(history.TypeDeleteHypha)

	if err, errtitle := CanDelete(u, h); errtitle != "" {
		hop.WithErrAbort(err)
		return hop, errtitle
	}

	originalText, _ := FetchTextPart(h)
	hop.
		WithFilesRemoved(h.TextPath, h.BinaryPath).
		WithMsg(fmt.Sprintf("Delete ‘%s’", h.Name)).
		WithUser(u).
		Apply()
	if !hop.HasErrors() {
		hyphae.BacklinksOnDelete(h, originalText)
		h.Delete()
	}
	return hop, ""
}
