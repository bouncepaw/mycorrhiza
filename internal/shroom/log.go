package shroom

import (
	"log/slog"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
)

func rejectRenameLog(h hyphae.Hypha, u *user.User, errmsg string) {
	slog.Info("Reject rename",
		"hyphaName", h.CanonicalName(),
		"username", u.Name,
		"errmsg", errmsg)
}

func rejectRemoveMediaLog(h hyphae.Hypha, u *user.User, errmsg string) {
	slog.Info("Reject remove media",
		"hyphaName", h.CanonicalName(),
		"username", u.Name,
		"errmsg", errmsg)
}

func rejectEditLog(h hyphae.Hypha, u *user.User, errmsg string) {
	slog.Info("Reject edit",
		"hyphaName", h.CanonicalName(),
		"username", u.Name,
		"errmsg", errmsg)
}

func rejectUploadMediaLog(h hyphae.Hypha, u *user.User, errmsg string) {
	slog.Info("Reject upload media",
		"hyphaName", h.CanonicalName(),
		"username", u.Name,
		"errmsg", errmsg)
}
