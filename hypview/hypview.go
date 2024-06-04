package hypview

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	viewutil2 "github.com/bouncepaw/mycorrhiza/web/viewutil"
	"html/template"
	"log"
	"strings"
)

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "editing hypha"}}Редактирование {{beautifulName .}}{{end}}
{{define "editing [[hypha]]"}}Редактирование <a href="/hypha/{{.}}">{{beautifulName .}}</a>{{end}}
{{define "creating [[hypha]]"}}Создание <a href="/hypha/{{.}}">{{beautifulName .}}</a>{{end}}
{{define "you're creating a new hypha"}}Вы создаёте новую гифу.{{end}}
{{define "describe your changes"}}Опишите ваши правки{{end}}
{{define "save"}}Сохранить{{end}}
{{define "preview"}}Предпросмотр{{end}}
{{define "previewing hypha"}}Предпросмотр «{{beautifulName .}}»{{end}}
{{define "preview tip"}}Заметьте, эта гифа ещё не сохранена. Вот её предпросмотр:{{end}}

{{define "markup"}}Разметка{{end}}
{{define "link"}}Ссылка{{end}}
{{define "link title"}}Текст{{end}}
{{define "heading"}}Заголовок{{end}}
{{define "bold"}}Жирный{{end}}
{{define "italic"}}Курсив{{end}}
{{define "highlight"}}Выделение{{end}}
{{define "underline"}}Подчеркивание{{end}}
{{define "mono"}}Моноширинный{{end}}
{{define "super"}}Надстрочный{{end}}
{{define "sub"}}Подстрочный{{end}}
{{define "strike"}}Зачёркнутый{{end}}
{{define "rocket"}}Ссылка-ракета{{end}}
{{define "transclude"}}Трансклюзия{{end}}
{{define "hr"}}Гориз. черта{{end}}
{{define "code"}}Код-блок{{end}}
{{define "bullets"}}Маркир. список{{end}}
{{define "numbers"}}Нумер. список{{end}}
{{define "mycomarkup help"}}<a href="/help/en/mycomarkup" class="shy-link">Подробнее</a> о Микоразметке{{end}}
{{define "actions"}}Действия{{end}}
{{define "current date"}}Текущая дата{{end}}
{{define "current time"}}Текущее время{{end}}
{{define "selflink"}}Ссылка на вас{{end}}

{{define "empty heading"}}Эта гифа не существует{{end}}
{{define "empty no rights"}}У вас нет прав для создания новых гиф. Вы можете:{{end}}
{{define "empty log in"}}Войти в свою учётную запись, если она у вас есть{{end}}
{{define "empty register"}}Создать новую учётную запись{{end}}
{{define "write a text"}}Написать текст{{end}}
{{define "write a text tip"}}Напишите заметку, дневник, статью, рассказ или иной текст с помощью <a href="/help/en/mycomarkup" class="shy-link">Микоразметки</a>. Сохраняется полная история правок документа.{{end}}
{{define "write a text writing conventions"}}Не забывайте следовать правилам оформления этой вики, если они имеются.{{end}}
{{define "write a text btn"}}Создать{{end}}
{{define "upload a media"}}Загрузить медиа{{end}}
{{define "upload a media tip"}}Загрузите изображение, видео или аудио. Распространённые форматы можно просматривать из браузера, остальные можно только скачать и просмотреть локально. Позже вы можете дописать пояснение к этому медиа.{{end}}
{{define "upload a media btn"}}Загрузить{{end}}

{{define "delete hypha?"}}Удалить {{beautifulName .}}?{{end}}
{{define "delete [[hypha]]?"}}Удалить <a href="/hypha/{{.}}">{{beautifulName .}}</a>?{{end}}
{{define "want to delete?"}}Вы действительно хотите удалить эту гифу?{{end}}
{{define "delete tip"}}Нельзя отменить удаление гифы, но её история останется доступной.{{end}}

{{define "rename hypha?"}}Переименовать {{beautifulName .}}?{{end}}
{{define "rename [[hypha]]?"}}Переименовать <a href="/hypha/{{.}}">{{beautifulName .}}</a>?{{end}}
{{define "new name"}}Новое название:{{end}}
{{define "rename recursively"}}Также переименовать подгифы{{end}}
{{define "rename tip"}}Переименовывайте аккуратно. <a href="/help/en/rename">Документация на английском.</a>{{end}}
{{define "leave redirection"}}Оставить перенаправление{{end}}

{{define "remove media from x?"}}Убрать медиа у {{beautifulName .}}?{{end}}
{{define "remove media from [[x]]?"}}Убрать медиа у <a href="/hypha/{{.MatchedHyphaName}}">{{beautifulName .MatchedHyphaName}}</a>?{{end}}
{{define "remove media for real?"}}Вы точно хотите убрать медиа у гифы «{{beautifulName .MatchedHyphaName}}»?{{end}}
`
	chainNaviTitle   viewutil2.Chain
	chainEditHypha   viewutil2.Chain
	chainEmptyHypha  viewutil2.Chain
	chainDeleteHypha viewutil2.Chain
	chainRenameHypha viewutil2.Chain
	chainRemoveMedia viewutil2.Chain
)

func Init() {
	chainNaviTitle = viewutil2.CopyEnRuWith(fs, "view_navititle.html", "")
	chainEditHypha = viewutil2.CopyEnRuWith(fs, "view_edit.html", ruTranslation)
	chainEmptyHypha = viewutil2.CopyEnRuWith(fs, "view_empty_hypha.html", ruTranslation)
	chainDeleteHypha = viewutil2.CopyEnRuWith(fs, "view_delete.html", ruTranslation)
	chainRenameHypha = viewutil2.CopyEnRuWith(fs, "view_rename.html", ruTranslation)
	chainRemoveMedia = viewutil2.CopyEnRuWith(fs, "view_remove_media.html", ruTranslation)
}

type editData struct {
	*viewutil2.BaseData
	HyphaName string
	IsNew     bool
	Content   string
	Message   string
	Preview   template.HTML
}

func EditHypha(meta viewutil2.Meta, hyphaName string, isNew bool, content string, message string, preview template.HTML) {
	viewutil2.ExecutePage(meta, chainEditHypha, editData{
		BaseData: &viewutil2.BaseData{
			Addr:        "/edit/" + hyphaName,
			EditScripts: cfg.EditScripts,
		},
		HyphaName: hyphaName,
		IsNew:     isNew,
		Content:   content,
		Message:   message,
		Preview:   preview,
	})
}

type renameData struct {
	*viewutil2.BaseData
	HyphaName               string
	LeaveRedirectionDefault bool
}

func RenameHypha(meta viewutil2.Meta, hyphaName string) {
	viewutil2.ExecutePage(meta, chainRenameHypha, renameData{
		BaseData: &viewutil2.BaseData{
			Addr: "/rename/" + hyphaName,
		},
		HyphaName:               hyphaName,
		LeaveRedirectionDefault: backlinks.BacklinksCount(hyphaName) != 0,
	})
}

type deleteRemoveMediaData struct {
	*viewutil2.BaseData
	HyphaName string
}

func DeleteHypha(meta viewutil2.Meta, hyphaName string) {
	viewutil2.ExecutePage(meta, chainDeleteHypha, deleteRemoveMediaData{
		BaseData: &viewutil2.BaseData{
			Addr: "/delete/" + hyphaName,
		},
		HyphaName: hyphaName,
	})
}

func RemoveMedia(meta viewutil2.Meta, hyphaName string) {
	viewutil2.ExecutePage(meta, chainRemoveMedia, deleteRemoveMediaData{
		BaseData: &viewutil2.BaseData{
			Addr: "/remove-media/" + hyphaName,
		},
		HyphaName: hyphaName,
	})
}

type emptyHyphaData struct {
	Meta              viewutil2.Meta
	HyphaName         string
	AllowRegistration bool
	UseAuth           bool
}

func EmptyHypha(meta viewutil2.Meta, hyphaName string) string {
	var buf strings.Builder
	if err := chainEmptyHypha.Get(meta).ExecuteTemplate(&buf, "empty hypha card", emptyHyphaData{
		Meta:              meta,
		HyphaName:         hyphaName,
		AllowRegistration: cfg.AllowRegistration,
		UseAuth:           cfg.UseAuth,
	}); err != nil {
		log.Println(err)
	}
	return buf.String()
}

type naviTitleData struct {
	HyphaNameParts            []string
	HyphaNamePartsWithParents []string
	Icon                      string
	HomeHypha                 string
}

func NaviTitle(meta viewutil2.Meta, hyphaName string) string {
	parts, partsWithParents := naviTitleify(hyphaName)
	var buf strings.Builder
	err := chainNaviTitle.Get(meta).ExecuteTemplate(&buf, "navititle", naviTitleData{
		HyphaNameParts:            parts,
		HyphaNamePartsWithParents: partsWithParents,
		Icon:                      cfg.NaviTitleIcon,
		HomeHypha:                 cfg.HomeHypha,
	})
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

func naviTitleify(hyphaName string) ([]string, []string) {
	var (
		prevAcc          = "/hypha"
		parts            = strings.Split(hyphaName, "/")
		partsWithParents []string
	)

	for _, part := range parts {
		prevAcc += "/" + part
		partsWithParents = append(partsWithParents, prevAcc)
	}

	return parts, partsWithParents
}
