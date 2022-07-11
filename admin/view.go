package admin

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"net/http"
)

const adminTranslationRu = `
{{define "panel title"}}Панель админстратора{{end}}
{{define "panel safe section title"}}Безопасная секция{{end}}
{{define "panel link about"}}Об этой вики{{end}}
{{define "panel update header"}}Обновить ссылки в верхней панели{{end}}
{{define "panel link user list"}}Список пользователей{{end}}
{{define "panel users"}}Управление пользователями{{end}}
{{define "panel unsafe section title"}}Опасная секция{{end}}
{{define "panel shutdown"}}Выключить вики{{end}}
{{define "panel reindex hyphae"}}Переиндексировать гифы{{end}}

{{define "manage users"}}Управление пользователями{{end}}
{{define "create user"}}Создать пользователя{{end}}
{{define "reindex users"}}Переиндексировать пользователей{{end}}
{{define "name"}}Имя{{end}}
{{define "group"}}Группа{{end}}
{{define "registered at"}}Зарегистрирован{{end}}
{{define "actions"}}Действия{{end}}
{{define "edit"}}Изменить{{end}}
`

var (
	//go:embed *.html
	fs         embed.FS
	panelChain viewutil.Chain
	listChain  viewutil.Chain
)

func Init(rtr *mux.Router) {
	rtr.HandleFunc("/shutdown", handlerAdminShutdown).Methods(http.MethodPost)
	rtr.HandleFunc("/reindex-users", handlerAdminReindexUsers).Methods(http.MethodPost)

	rtr.HandleFunc("/new-user", handlerAdminUserNew).Methods(http.MethodGet, http.MethodPost)
	rtr.HandleFunc("/users/{username}/edit", handlerAdminUserEdit).Methods(http.MethodGet, http.MethodPost)
	rtr.HandleFunc("/users/{username}/delete", handlerAdminUserDelete).Methods(http.MethodGet, http.MethodPost)
	rtr.HandleFunc("/users", handlerAdminUsers)

	rtr.HandleFunc("/", handlerAdmin)

	panelChain = viewutil.CopyEnRuWith(fs, "view_panel.html", adminTranslationRu)
	listChain = viewutil.CopyEnRuWith(fs, "view_user_list.html", adminTranslationRu)
}

func viewPanel(meta viewutil.Meta) {
	viewutil.ExecutePage(meta, panelChain, &viewutil.BaseData{})
}

type listData struct {
	*viewutil.BaseData
	UserHypha string
	Users     []*user.User
}

func viewList(meta viewutil.Meta, users []*user.User) {
	viewutil.ExecutePage(meta, listChain, listData{
		BaseData:  &viewutil.BaseData{},
		UserHypha: cfg.UserHypha,
		Users:     users,
	})
}
