package web

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	viewutil2 "github.com/bouncepaw/mycorrhiza/web/viewutil"
)

//go:embed views/*.html
var fs embed.FS

var pageOrphans, pageBacklinks, pageUserList, pageChangePassword *newtmpl.Page

var panelChain, listChain, newUserChain, editUserChain, deleteUserChain viewutil2.Chain

func initPages() {

	panelChain = viewutil2.CopyEnRuWith(fs, "views/admin-panel.html", adminTranslationRu)
	listChain = viewutil2.CopyEnRuWith(fs, "views/admin-user-list.html", adminTranslationRu)
	newUserChain = viewutil2.CopyEnRuWith(fs, "views/admin-new-user.html", adminTranslationRu)
	editUserChain = viewutil2.CopyEnRuWith(fs, "views/admin-edit-user.html", adminTranslationRu)
	deleteUserChain = viewutil2.CopyEnRuWith(fs, "views/admin-delete-user.html", adminTranslationRu)

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
	pageChangePassword = newtmpl.NewPage(fs, "views/change-password.html", map[string]string{
		"change password":           "Сменить пароль",
		"confirm password":          "Повторите пароль",
		"current password":          "Текущий пароль",
		"non local password change": "Пароль можно поменять только местным аккаунтам. Telegram-аккаунтам нельзя.",
		"password":                  "Пароль",
		"submit":                    "Поменять",
	})
}
