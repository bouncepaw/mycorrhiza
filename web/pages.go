package web

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

//go:embed views/*.html
var fs embed.FS

var pageOrphans, pageBacklinks, pageUserList, pageChangePassword *newtmpl.Page
var pageHyphaDelete, pageHyphaEdit, pageHyphaEmpty, pageHypha *newtmpl.Page
var panelChain, listChain, newUserChain, editUserChain, deleteUserChain viewutil.Chain

func initPages() {

	panelChain = viewutil.CopyEnRuWith(fs, "views/admin-panel.html", adminTranslationRu)
	listChain = viewutil.CopyEnRuWith(fs, "views/admin-user-list.html", adminTranslationRu)
	newUserChain = viewutil.CopyEnRuWith(fs, "views/admin-new-user.html", adminTranslationRu)
	editUserChain = viewutil.CopyEnRuWith(fs, "views/admin-edit-user.html", adminTranslationRu)
	deleteUserChain = viewutil.CopyEnRuWith(fs, "views/admin-delete-user.html", adminTranslationRu)

	pageOrphans = newtmpl.NewPage(fs, map[string]string{
		"orphaned hyphae":    "Гифы-сироты",
		"orphan description": "Ниже перечислены гифы без ссылок на них.",
	}, "views/orphans.html", map[string]string{
		"orphaned hyphae":    "Гифы-сироты",
		"orphan description": "Ниже перечислены гифы без ссылок на них.",
	})
	pageBacklinks = newtmpl.NewPage(fs, map[string]string{
		"backlinks to text": `Обратные ссылки на {{.}}`,
		"backlinks to link": `Обратные ссылки на <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"description":       `Ниже перечислены гифы, на которых есть ссылка на эту гифу, трансклюзия этой гифы или эта гифа вставлена как изображение.`,
	}, "views/backlinks.html", map[string]string{
		"backlinks to text": `Обратные ссылки на {{.}}`,
		"backlinks to link": `Обратные ссылки на <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"description":       `Ниже перечислены гифы, на которых есть ссылка на эту гифу, трансклюзия этой гифы или эта гифа вставлена как изображение.`,
	})
	pageUserList = newtmpl.NewPage(fs, map[string]string{
		"title":          "Список пользователей",
		"administrators": "Администраторы",
		"moderators":     "Модераторы",
		"editors":        "Редакторы",
		"readers":        "Читатели",
	}, "views/user-list.html", map[string]string{
		"title":          "Список пользователей",
		"administrators": "Администраторы",
		"moderators":     "Модераторы",
		"editors":        "Редакторы",
		"readers":        "Читатели",
	})
	pageChangePassword = newtmpl.NewPage(fs, map[string]string{
		"change password":           "Сменить пароль",
		"confirm password":          "Повторите пароль",
		"current password":          "Текущий пароль",
		"non local password change": "Пароль можно поменять только местным аккаунтам. Telegram-аккаунтам нельзя.",
		"password":                  "Пароль",
		"submit":                    "Поменять",
	}, "views/change-password.html", map[string]string{
		"change password":           "Сменить пароль",
		"confirm password":          "Повторите пароль",
		"current password":          "Текущий пароль",
		"non local password change": "Пароль можно поменять только местным аккаунтам. Telegram-аккаунтам нельзя.",
		"password":                  "Пароль",
		"submit":                    "Поменять",
	})
	pageHyphaDelete = newtmpl.NewPage(fs, map[string]string{
		"delete hypha?":     "Удалить {{beautifulName .}}?",
		"delete [[hypha]]?": "Удалить <a href=\"/hypha/{{.}}\">{{beautifulName .}}</a>?",
		"want to delete?":   "Вы действительно хотите удалить эту гифу?",
		"delete tip":        "Нельзя отменить удаление гифы, но её история останется доступной.",
	}, "views/hypha-delete.html", map[string]string{
		"delete hypha?":     "Удалить {{beautifulName .}}?",
		"delete [[hypha]]?": "Удалить <a href=\"/hypha/{{.}}\">{{beautifulName .}}</a>?",
		"want to delete?":   "Вы действительно хотите удалить эту гифу?",
		"delete tip":        "Нельзя отменить удаление гифы, но её история останется доступной.",
	})
	pageHyphaEdit = newtmpl.NewPage(fs, map[string]string{
		"editing hypha":               `Редактирование {{beautifulName .}}`,
		"editing [[hypha]]":           `Редактирование <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"creating [[hypha]]":          `Создание <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"you're creating a new hypha": `Вы создаёте новую гифу.`,
		"describe your changes":       `Опишите ваши правки`,
		"save":                        `Сохранить`,
		"preview":                     `Предпросмотр`,
		"previewing hypha":            `Предпросмотр «{{beautifulName .}}»`,
		"preview tip":                 `Заметьте, эта гифа ещё не сохранена. Вот её предпросмотр:`,

		"markup":          `Разметка`,
		"link":            `Ссылка`,
		"link title":      `Текст`,
		"heading":         `Заголовок`,
		"bold":            `Жирный`,
		"italic":          `Курсив`,
		"highlight":       `Выделение`,
		"underline":       `Подчеркивание`,
		"mono":            `Моноширинный`,
		"super":           `Надстрочный`,
		"sub":             `Подстрочный`,
		"strike":          `Зачёркнутый`,
		"rocket":          `Ссылка-ракета`,
		"transclude":      `Трансклюзия`,
		"hr":              `Гориз. черта`,
		"code":            `Код-блок`,
		"bullets":         `Маркир. список`,
		"numbers":         `Нумер. список`,
		"mycomarkup help": `<a href="/help/en/mycomarkup" class="shy-link">Подробнее</a> о Микоразметке`,
		"actions":         `Действия`,
		"current date":    `Текущая дата`,
		"current time":    `Текущее время`,
		"selflink":        `Ссылка на вас`,
	}, "views/hypha-edit.html", map[string]string{
		"editing hypha":               `Редактирование {{beautifulName .}}`,
		"editing [[hypha]]":           `Редактирование <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"creating [[hypha]]":          `Создание <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"you're creating a new hypha": `Вы создаёте новую гифу.`,
		"describe your changes":       `Опишите ваши правки`,
		"save":                        `Сохранить`,
		"preview":                     `Предпросмотр`,
		"previewing hypha":            `Предпросмотр «{{beautifulName .}}»`,
		"preview tip":                 `Заметьте, эта гифа ещё не сохранена. Вот её предпросмотр:`,

		"markup":          `Разметка`,
		"link":            `Ссылка`,
		"link title":      `Текст`,
		"heading":         `Заголовок`,
		"bold":            `Жирный`,
		"italic":          `Курсив`,
		"highlight":       `Выделение`,
		"underline":       `Подчеркивание`,
		"mono":            `Моноширинный`,
		"super":           `Надстрочный`,
		"sub":             `Подстрочный`,
		"strike":          `Зачёркнутый`,
		"rocket":          `Ссылка-ракета`,
		"transclude":      `Трансклюзия`,
		"hr":              `Гориз. черта`,
		"code":            `Код-блок`,
		"bullets":         `Маркир. список`,
		"numbers":         `Нумер. список`,
		"mycomarkup help": `<a href="/help/en/mycomarkup" class="shy-link">Подробнее</a> о Микоразметке`,
		"actions":         `Действия`,
		"current date":    `Текущая дата`,
		"current time":    `Текущее время`,
		"selflink":        `Ссылка на вас`,
	})
	pageHyphaEmpty = newtmpl.NewPage(fs, map[string]string{
		"empty heading":                    `Эта гифа не существует`,
		"empty no rights":                  `У вас нет прав для создания новых гиф. Вы можете:`,
		"empty log in":                     `Войти в свою учётную запись, если она у вас есть`,
		"empty register":                   `Создать новую учётную запись`,
		"write a text":                     `Написать текст`,
		"write a text tip":                 `Напишите заметку, дневник, статью, рассказ или иной текст с помощью <a href="/help/en/mycomarkup" class="shy-link">Микоразметки</a>. Сохраняется полная история правок документа.`,
		"write a text writing conventions": `Не забывайте следовать правилам оформления этой вики, если они имеются.`,
		"write a text btn":                 `Создать`,
		"upload a media":                   `Загрузить медиа`,
		"upload a media tip":               `Загрузите изображение, видео или аудио. Распространённые форматы можно просматривать из браузера, остальные можно только скачать и просмотреть локально. Позже вы можете дописать пояснение к этому медиа.`,
		"upload a media btn":               `Загрузить`,
	}, "views/hypha-empty.html", map[string]string{
		"empty heading":                    `Эта гифа не существует`,
		"empty no rights":                  `У вас нет прав для создания новых гиф. Вы можете:`,
		"empty log in":                     `Войти в свою учётную запись, если она у вас есть`,
		"empty register":                   `Создать новую учётную запись`,
		"write a text":                     `Написать текст`,
		"write a text tip":                 `Напишите заметку, дневник, статью, рассказ или иной текст с помощью <a href="/help/en/mycomarkup" class="shy-link">Микоразметки</a>. Сохраняется полная история правок документа.`,
		"write a text writing conventions": `Не забывайте следовать правилам оформления этой вики, если они имеются.`,
		"write a text btn":                 `Создать`,
		"upload a media":                   `Загрузить медиа`,
		"upload a media tip":               `Загрузите изображение, видео или аудио. Распространённые форматы можно просматривать из браузера, остальные можно только скачать и просмотреть локально. Позже вы можете дописать пояснение к этому медиа.`,
		"upload a media btn":               `Загрузить`,
	})
	pageHypha = newtmpl.NewPage(fs, map[string]string{}, "views/hypha.html", map[string]string{})
}
