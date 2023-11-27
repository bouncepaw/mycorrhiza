package settings

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/user"
)

// TODO: translate untranslated strings
const settingsTranslationRu = `
{{define "change password"}}Change password{{end}}
{{define "confirm password"}}Confirm password{{end}}
{{define "current password"}}Current password{{end}}
{{define "non local password change"}}Non-local accounts cannot have their passwords changed.{{end}}
{{define "password"}}Password{{end}}
{{define "submit"}}Submit{{end}}
`

var (
	//go:embed *.html
	fs                                                                  embed.FS
	changePassowrdChain viewutil.Chain
)

func Init(rtr *mux.Router) {
	rtr.HandleFunc("/change-password", handlerUserChangePassword).Methods(http.MethodGet, http.MethodPost)

	changePassowrdChain = viewutil.CopyEnRuWith(fs, "view_change_password.html", settingsTranslationRu)
}

func changePasswordPage(meta viewutil.Meta, form util.FormData, u *user.User) {
	viewutil.ExecutePage(meta, changePassowrdChain, changePasswordData{
		BaseData: &viewutil.BaseData{},
		Form: form,
		U: u,
	})
}

type changePasswordData struct {
	*viewutil.BaseData
	Form util.FormData
	U *user.User
}
