package hyphae

import (
	"errors"
	"fmt"
	"log"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/user"
)

func rejectDeleteLog(h *Hypha, u *user.User, errmsg string) {
	log.Printf("Reject delete ‘%s’ by @%s: %s\n", h.Name, u.Name, errmsg)
}

// CanDelete checks if given user can delete given hypha.
func (h *Hypha) CanDelete(u *user.User) (err error, errtitle string) {
	// First, check if can unattach at all
	if !u.CanProceed("delete-confirm") {
		rejectDeleteLog(h, u, "no rights")
		return errors.New("Not enough rights to delete, you must be a moderator"), "Not enough rights"
	}

	if !h.Exists {
		rejectDeleteLog(h, u, "does not exist")
		return errors.New("Cannot delete this hypha because it does not exist"), "Does not exist"
	}

	return nil, ""
}

// DeleteHypha deletes hypha and makes a history record about that.
func (h *Hypha) DeleteHypha(u *user.User) (hop *history.HistoryOp, errtitle string) {
	h.Lock()
	defer h.Unlock()
	hop = history.Operation(history.TypeDeleteHypha)

	if err, errtitle := h.CanDelete(u); errtitle != "" {
		hop.WithError(err)
		return hop, errtitle
	}

	hop.
		WithFilesRemoved(h.TextPath, h.BinaryPath).
		WithMsg(fmt.Sprintf("Delete ‘%s’", h.Name)).
		WithUser(u).
		Apply()
	if len(hop.Errs) == 0 {
		h.delete()
	}
	return hop, ""
}
