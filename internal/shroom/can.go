package shroom

import (
	"errors"
	"log/slog"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
)

func CanAttach(u *user.User, h hyphae.Hypha) error {
	if !u.CanProceed("upload-binary") {
		slog.Info("Reject upload media",
			"hyphaName", h.CanonicalName(),
			"username", u.Name,
			"reason", "unauthorized")
		return errors.New("ui.act_norights_media")
	}

	return nil
}

func CanEdit(u *user.User, h hyphae.Hypha) error {
	if !u.CanProceed("upload-text") {
		slog.Info("Reject edit",
			"hyphaName", h.CanonicalName(),
			"username", u.Name,
			"reason", "unauthorized")
		return errors.New("ui.act_norights_edit")
	}

	return nil
}
