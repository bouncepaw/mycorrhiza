package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
)

func canFactory(
	rejectLogger func(*hyphae.Hypha, *user.User, string),
	action string,
	dispatcher func(*hyphae.Hypha, *user.User, *l18n.Localizer) (string, string),
	noRightsMsg string,
	notExistsMsg string,
	mustExist bool,
) func(*user.User, *hyphae.Hypha, *l18n.Localizer) (string, error) {
	return func(u *user.User, h *hyphae.Hypha, lc *l18n.Localizer) (string, error) {
		if !u.CanProceed(action) {
			rejectLogger(h, u, "no rights")
			return lc.Get("ui.act_no_rights"), errors.New(lc.Get(noRightsMsg))
		}

		if mustExist && !h.Exists {
			rejectLogger(h, u, "does not exist")
			return lc.Get("ui.act_notexist"), errors.New(lc.Get(notExistsMsg))
		}

		if dispatcher == nil {
			return "", nil
		}
		errmsg, errtitle := dispatcher(h, u, lc)
		if errtitle == "" {
			return "", nil
		}
		return errtitle, errors.New(errmsg)
	}
}

// CanDelete and etc are hyphae operation checkers based on user rights and hyphae existence.
var (
	CanDelete = canFactory(
		rejectDeleteLog,
		"delete-confirm",
		nil,
		"ui.act_norights_delete",
		"ui.act_notexist_delete",
		true,
	)

	CanRename = canFactory(
		rejectRenameLog,
		"rename-confirm",
		nil,
		"ui.act_norights_rename",
		"ui.act_notexist_rename",
		true,
	)

	CanUnattach = canFactory(
		rejectUnattachLog,
		"unattach-confirm",
		func(h *hyphae.Hypha, u *user.User, lc *l18n.Localizer) (errmsg, errtitle string) {
			if h.BinaryPath == "" {
				rejectUnattachLog(h, u, "no amnt")
				return lc.Get("ui.act_noattachment_tip"), lc.Get("ui.act_noattachment")
			}

			return "", ""
		},
		"ui.act_norights_unattach",
		"ui.act_notexist_unattach",
		true,
	)

	CanEdit = canFactory(
		rejectEditLog,
		"upload-text",
		nil,
		"ui.act_norights_edit",
		"You cannot edit a hypha that does not exist",
		false,
	)

	CanAttach = canFactory(
		rejectAttachLog,
		"upload-binary",
		nil,
		"ui.act_norights_attach",
		"You cannot attach a hypha that does not exist",
		false,
	)
)

/* I've left 'not exists' messages for edit and attach out of translation as they are not used -- chekoopa */
