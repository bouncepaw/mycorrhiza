package admin

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
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
`

var (
	//go:embed *.html
	fs         embed.FS
	panelChain viewutil.Chain
)

func Init(rtr *mux.Router) {
	panelChain = viewutil.CopyEnRuWith(fs, "admin.html", adminTranslationRu)
}

func AdminPanel(meta viewutil.Meta) {
	viewutil.ExecutePage(meta, panelChain, &viewutil.BaseData{})
}
