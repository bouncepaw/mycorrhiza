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
var pageRevision, pageMedia *newtmpl.Page
var pageAuthLock, pageAuthLogin, pageAuthLogout, pageAuthRegister, pageAuthTelegram *newtmpl.Page

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
	}, "views/orphans.html")
	pageBacklinks = newtmpl.NewPage(fs, map[string]string{
		"backlinks to text": `Обратные ссылки на {{.}}`,
		"backlinks to link": `Обратные ссылки на <a href="/hypha/{{.}}">{{beautifulName .}}</a>`,
		"description":       `Ниже перечислены гифы, на которых есть ссылка на эту гифу, трансклюзия этой гифы или эта гифа вставлена как изображение.`,
	}, "views/backlinks.html")
	pageUserList = newtmpl.NewPage(fs, map[string]string{
		"title":          "Список пользователей",
		"administrators": "Администраторы",
		"moderators":     "Модераторы",
		"editors":        "Редакторы",
		"readers":        "Читатели",
	}, "views/user-list.html")
	pageChangePassword = newtmpl.NewPage(fs, map[string]string{
		"change password":           "Сменить пароль",
		"confirm password":          "Повторите пароль",
		"current password":          "Текущий пароль",
		"non local password change": "Пароль можно поменять только местным аккаунтам. Telegram-аккаунтам нельзя.",
		"password":                  "Пароль",
		"submit":                    "Поменять",
	}, "views/change-password.html")
	pageHyphaDelete = newtmpl.NewPage(fs, map[string]string{
		"delete hypha?":     "Удалить {{beautifulName .}}?",
		"delete [[hypha]]?": "Удалить <a href=\"/hypha/{{.}}\">{{beautifulName .}}</a>?",
		"want to delete?":   "Вы действительно хотите удалить эту гифу?",
		"delete tip":        "Нельзя отменить удаление гифы, но её история останется доступной.",
	}, "views/hypha-delete.html")
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

		"markup":             `Разметка`,
		"link":               `Ссылка`,
		"link title":         `Текст`,
		"heading":            `Заголовок`,
		"bold":               `Жирный`,
		"italic":             `Курсив`,
		"highlight":          `Выделение`,
		"underline":          `Подчеркивание`,
		"mono":               `Моноширинный`,
		"super":              `Надстрочный`,
		"sub":                `Подстрочный`,
		"strike":             `Зачёркнутый`,
		"rocket":             `Ссылка-ракета`,
		"transclude":         `Трансклюзия`,
		"hr":                 `Гориз. черта`,
		"code":               `Код-блок`,
		"bullets":            `Маркир. список`,
		"numbers":            `Нумер. список`,
		"mycomarkup help":    `<a href="/help/en/mycomarkup" class="shy-link">Подробнее</a> о Микоразметке`,
		"actions":            `Действия`,
		"current date local": `Местная дата`,
		"current time local": `Местное время`,
		"current date utc":   "Дата UTC",
		"current time utc":   "Время UTC",
		"selflink":           `Ссылка на вас`,
	}, "views/hypha-edit.html")
	pageHypha = newtmpl.NewPage(fs, map[string]string{
		"edit text":     "Редактировать",
		"log out":       "Выйти",
		"admin panel":   "Админка",
		"subhyphae":     "Подгифы",
		"history":       "История",
		"rename":        "Переименовать",
		"delete":        "Удалить",
		"view markup":   "Посмотреть разметку",
		"manage media":  "Медиа",
		"turn to media": "Превратить в медиа-гифу",
		"backlinks":     "{{.BacklinkCount}} обратн{{if eq .BacklinkCount 1}}ая ссылка{{else if and (le .BacklinkCount 4) (gt .BacklinkCount 1)}}ые ссылки{{else}}ых ссылок{{end}}",

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
	}, "views/hypha.html")
	pageRevision = newtmpl.NewPage(fs, map[string]string{
		"revision warning": "Обратите внимание, просмотр медиа в истории пока что недоступен.",
		"revision link":    "Посмотреть Микоразметку для этой ревизии",
		"hypha at rev":     "{{.HyphaName}} на {{.RevHash}}",
	}, "views/hypha-revision.html")
	pageMedia = newtmpl.NewPage(fs, map[string]string{ // TODO: сделать новый перевод
		"media title":    "Медиа «{{.HyphaName | beautifulLink}}»",
		"tip":            "На этой странице вы можете управлять медиа.",
		"empty":          "Эта гифа не имеет медиа, здесь вы можете его загрузить.",
		"what is media?": "Что такое медиа?",
		"stat":           "Свойства",
		"stat size":      "Размер файла:",
		"stat mime":      "MIME-тип:",

		"upload title": "Прикрепить",
		"upload tip":   "Вы можете загрузить новое медиа. Пожалуйста, не загружайте слишком большие изображения без необходимости, чтобы впоследствии не ждать её долгую загрузку.",
		"upload btn":   "Загрузить",

		"remove title": "Открепить",
		"remove tip":   "Заметьте, чтобы заменить медиа, вам не нужно его перед этим откреплять.",
		"remove btn":   "Открепить",
	}, "views/hypha-media.html")

	pageAuthLock = newtmpl.NewPage(fs, map[string]string{
		"lock title": "Доступ закрыт",
		"username":   "Логин",
		"password":   "Пароль",
		"log in":     "Войти",
	}, "views/auth-telegram.html", "views/auth-lock.html")

	pageAuthLogin = newtmpl.NewPage(fs, map[string]string{
		"username":       "Логин",
		"password":       "Пароль",
		"log in":         "Войти",
		"cookie tip":     "Отправляя эту форму, вы разрешаете вики хранить cookie в вашем браузере. Это позволит движку связывать ваши правки с вашей учётной записью. Вы будете авторизованы, пока не выйдете из учётной записи.",
		"log in to x":    "Войти в {{.}}",
		"auth disabled":  "Аутентификация отключена. Вы можете делать правки анонимно.",
		"error username": "Неизвестное имя пользователя.",
		"error password": "Неправильный пароль.",
		"error telegram": "Не удалось войти через Телеграм.",
		"go home":        "Домой",
	}, "views/auth-telegram.html", "views/auth-login.html")

	pageAuthLogout = newtmpl.NewPage(fs, map[string]string{
		"log out?":            "Выйти?",
		"log out":             "Выйти",
		"cannot log out anon": "Вы не можете выйти, потому что ещё не вошли.",
		"log in":              "Войти",
		"go home":             "Домой",
	}, "views/auth-logout.html")

	pageAuthRegister = newtmpl.NewPage(fs, map[string]string{
		"username":      "Логин",
		"password":      "Пароль",
		"cookie tip":    "Отправляя эту форму, вы разрешаете вики хранить cookie в вашем браузере. Это позволит движку связывать ваши правки с вашей учётной записью. Вы будете авторизованы, пока не выйдете из учётной записи.",
		"password tip":  "Сервер хранит ваш пароль в зашифрованном виде, даже администраторы не смогут его прочесть.",
		"register btn":  "Зарегистрироваться",
		"register on x": "Регистрация на {{.}}",
	}, "views/auth-telegram.html", "views/auth-register.html")

}
