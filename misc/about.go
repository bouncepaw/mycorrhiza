package misc

import (
	"log/slog"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/internal/version"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
)

var pageAbout *newtmpl.Page

type aboutData struct {
	newtmpl.BaseData
	Version           string
	Cfg               map[string]any
	Admins            []string
	UserCount         uint64
	RegistrationLimit uint64
}

func viewAbout(meta viewdata.Meta) {
	var data aboutData
	data.Version = version.Short
	data.Admins = usersvc.ListUsersWithGroup(userdomain.GroupAdmin)
	data.UserCount = usersvc.Count()
	data.RegistrationLimit = cfg.RegistrationLimit
	data.Cfg = map[string]any{
		"UseAuth":           cfg.UseAuth,
		"WikiName":          cfg.WikiName,
		"HomeHypha":         cfg.HomeHypha,
		"UserHypha":         cfg.UserHypha,
		"TelegramEnabled":   cfg.TelegramEnabled,
		"RegistrationLimit": cfg.RegistrationLimit,
	}
	if err := pageAbout.RenderDataTo(meta, &data); err != nil {
		slog.Error("Failed to execute page", "page", "about", "err", err)
	}
}
