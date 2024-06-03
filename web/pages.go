package web

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/newtmpl"
)

//go:embed views/*.html
var fs embed.FS

var pageOrphans, pageBacklinks, pageUserList *newtmpl.Page

func initPages() {
	pageOrphans = newtmpl.NewPage(fs, "views/orphans.html", map[string]string{
		"orphaned hyphae":    "Гифы-сироты",
		"orphan description": "Ниже перечислены гифы без ссылок на них.",
	})
	pageBacklinks = newtmpl.NewPage(fs, "views/backlinks.html", map[string]string{
		"backlinks to text": `Обратные ссылки на {{.}}`,
		"backlinks to link": `Обратные ссылки на <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"description":       `Ниже перечислены гифы, на которых есть ссылка на эту гифу, трансклюзия этой гифы или эта гифа вставлена как изображение.`,
	})
	pageUserList = newtmpl.NewPage(fs, "views/user-list.html", map[string]string{
		"title":          "Список пользователей",
		"administrators": "Администраторы",
		"moderators":     "Модераторы",
		"editors":        "Редакторы",
		"readers":        "Читатели",
	})
}
