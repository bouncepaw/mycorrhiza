package hyphae

import (
	"errors"
	"fmt"
	"log"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/user"
)

func rejectUnattachLog(h *Hypha, u *user.User, errmsg string) {
	log.Printf("Reject unattach ‘%s’ by @%s: %s\n", h.Name, u.Name, errmsg)
}

// CanUnattach checks if given user can unattach given hypha. If they can, `errtitle` is an empty string and `err` is nil. If they cannot, `errtitle` is not an empty string, and `err` is an error.
func (h *Hypha) CanUnattach(u *user.User) (err error, errtitle string) {
	if !u.CanProceed("unattach-confirm") {
		rejectUnattachLog(h, u, "no rights")
		return errors.New("Not enough rights to unattach, you must be a trusted editor"), "Not enough rights"
	}

	if !h.Exists {
		rejectUnattachLog(h, u, "does not exist")
		return errors.New("Cannot unattach this hypha because it does not exist"), "Does not exist"
	}

	if h.BinaryPath == "" {
		rejectUnattachLog(h, u, "no amnt")
		return errors.New("Cannot unattach this hypha because it has no attachment"), "No attachment"
	}

	return nil, ""
}

// UnattachHypha unattaches hypha and makes a history record about that.
func (h *Hypha) UnattachHypha(u *user.User) (hop *history.HistoryOp, errtitle string) {
	h.Lock()
	defer h.Unlock()
	hop = history.Operation(history.TypeUnattachHypha)

	if err, errtitle := h.CanUnattach(u); errtitle != "" {
		hop.WithError(err)
		return hop, errtitle
	}

	hop.
		WithFilesRemoved(h.BinaryPath).
		WithMsg(fmt.Sprintf("Unattach ‘%s’", h.Name)).
		WithUser(u).
		Apply()

	if len(hop.Errs) > 0 {
		rejectUnattachLog(h, u, "fail")
		return hop.WithError(errors.New(fmt.Sprintf("Could not unattach this hypha due to internal server errors: <code>%v</code>", hop.Errs))), "Error"
	}

	if h.BinaryPath != "" {
		h.BinaryPath = ""
	}
	// If nothing is left of the hypha
	if h.TextPath == "" {
		h.delete()
	}
	return hop, ""
}
