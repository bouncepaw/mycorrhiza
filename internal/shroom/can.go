package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
)

func CanAttach(u *user.User, h hyphae.Hypha) error {
	if !u.CanProceed("upload-binary") {
		rejectUploadMediaLog(h, u, "no rights")
		return errors.New("ui.act_norights_media")
	}

	return nil
}

func CanEdit(u *user.User, h hyphae.Hypha) error {
	if !u.CanProceed("upload-text") {
		rejectEditLog(h, u, "no rights")
		return errors.New("ui.act_norights_edit")
	}

	return nil
}
