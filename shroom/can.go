package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
)

// TODO: get rid of this abomination

func canFactory(
	rejectLogger func(hyphae.Hypha, *user.User, string),
	action string,
	dispatcher func(hyphae.Hypha, *user.User, *l18n.Localizer) (string, string),
	noRightsMsg string,
	notExistsMsg string,
	mustExist bool,
) func(*user.User, hyphae.Hypha, *l18n.Localizer) error {
	return func(u *user.User, h hyphae.Hypha, lc *l18n.Localizer) error {
		if !u.CanProceed(action) {
			rejectLogger(h, u, "no rights")
			return errors.New(noRightsMsg)
		}

		if mustExist {
			switch h.(type) {
			case *hyphae.EmptyHypha:
				rejectLogger(h, u, "does not exist")
				return errors.New(notExistsMsg)
			}
		}

		if dispatcher == nil {
			return nil
		}
		errmsg, errtitle := dispatcher(h, u, lc)
		if errtitle == "" {
			return nil
		}
		return errors.New(errmsg)
	}
}

// CanDelete and etc are hyphae operation checkers based on user rights and hyphae existence.
var (
	CanEdit = canFactory(
		rejectEditLog,
		"upload-text",
		nil,
		"ui.act_norights_edit",
		"You cannot edit a hypha that does not exist",
		false,
	)

	CanAttach = canFactory(
		rejectUploadMediaLog,
		"upload-binary",
		nil,
		"ui.act_norights_media",
		"You cannot attach a hypha that does not exist",
		false,
	)
)

/* I've left 'not exists' messages for edit and attach out of translation as they are not used -- chekoopa */
