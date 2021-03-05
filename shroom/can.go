package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
)

func canFactory(
	rejectLogger func(*hyphae.Hypha, *user.User, string),
	action string,
	dispatcher func(*hyphae.Hypha, *user.User) (string, string),
	noRightsMsg string,
	notExistsMsg string,
	careAboutExistince bool,
) func(*user.User, *hyphae.Hypha) (error, string) {
	return func(u *user.User, h *hyphae.Hypha) (error, string) {
		if !u.CanProceed(action) {
			rejectLogger(h, u, "no rights")
			return errors.New(noRightsMsg), "Not enough rights"
		}

		if careAboutExistince && !h.Exists {
			rejectLogger(h, u, "does not exist")
			return errors.New(notExistsMsg), "Does not exist"
		}

		if dispatcher == nil {
			return nil, ""
		}
		errmsg, errtitle := dispatcher(h, u)
		if errtitle == "" {
			return nil, ""
		}
		return errors.New(errmsg), errtitle
	}
}

var (
	CanDelete = canFactory(
		rejectDeleteLog,
		"delete-confirm",
		nil,
		"Not enough rights to delete, you must be a moderator",
		"Cannot delete this hypha because it does not exist",
		true,
	)

	CanRename = canFactory(
		rejectRenameLog,
		"rename-confirm",
		nil,
		"Not enough rights to rename, you must be a trusted editor",
		"Cannot rename this hypha because it does not exist",
		true,
	)

	CanUnattach = canFactory(
		rejectUnattachLog,
		"unattach-confirm",
		func(h *hyphae.Hypha, u *user.User) (errmsg, errtitle string) {
			if h.BinaryPath == "" {
				rejectUnattachLog(h, u, "no amnt")
				return "Cannot unattach this hypha because it has no attachment", "No attachment"
			}

			return "", ""
		},
		"Not enough rights to unattach, you must be a trusted editor",
		"Cannot unattach this hypha because it does not exist",
		true,
	)

	CanEdit = canFactory(
		rejectEditLog,
		"upload-text",
		nil,
		"You must be an editor to edit a hypha",
		"You cannot edit a hypha that does not exist",
		false,
	)

	CanAttach = canFactory(
		rejectAttachLog,
		"upload-binary",
		nil,
		"You must be an editor to attach a hypha",
		"You cannot attach a hypha that does not exist",
		false,
	)
)
