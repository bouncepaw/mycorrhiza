package hypview

import (
	"embed"
	"html/template"
	"log"
	"strings"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/viewutil"
)

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "editing hypha"}}Редактирование {{beautifulName .}}{{end}}
{{define "editing [[hypha]]"}}Редактирование <a href="/hypha/{{.}}">{{beautifulName .}}</a>{{end}}
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
{{define "mycomarkup help"}}<a href="/help/en/mycomarkup" class="shy-link">Подробнее</a> о микоразметке{{end}}
{{define "actions"}}Действия{{end}}
{{define "current date"}}Текущая дата{{end}}
{{define "current time"}}Текущее время{{end}}
{{define "selflink"}}Ссылка на вас{{end}}

{{define "empty heading"}}Эта гифа не существует{{end}}
{{define "empty no rights"}}У вас нет прав для создания новых гиф. Вы можете:{{end}}
{{define "empty log in"}}Войти в свою учётную запись, если она у вас есть{{end}}
{{define "empty register"}}Создать новую учётную запись{{end}}
{{define "write a text"}}Написать текст{{end}}
{{define "write a text tip"}}Напишите заметку, дневник, статью, рассказ или иной текст с помощью <a href="/help/en/mycomarkup" class="shy-link">микоразметки</a>. Сохраняется полная история правок документа.{{end}}
{{define "write a text writing conventions"}}Не забывайте следовать правилам оформления этой вики, если они имеются.{{end}}
{{define "write a text btn"}}Создать{{end}}
{{define "upload a media"}}Загрузить медиа{{end}}
{{define "upload a media tip"}}Загрузите изображение, видео или аудио. Распространённые форматы можно просматривать из браузера, остальные – просто скачать. Позже вы можете дописать пояснение к этому медиа.{{end}}
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
{{define "leave redirections"}}Оставить перенаправления{{end}}
`
	chainNaviTitle   viewutil.Chain
	chainEditHypha   viewutil.Chain
	chainEmptyHypha  viewutil.Chain
	chainDeleteHypha viewutil.Chain
	chainRenameHypha viewutil.Chain
)

func Init() {
	chainNaviTitle = viewutil.CopyEnRuWith(fs, "view_navititle.html", "")
	chainEditHypha = viewutil.CopyEnRuWith(fs, "view_edit.html", ruTranslation)
	chainEmptyHypha = viewutil.CopyEnRuWith(fs, "view_empty_hypha.html", ruTranslation)
	chainDeleteHypha = viewutil.CopyEnRuWith(fs, "view_delete.html", ruTranslation)
	chainRenameHypha = viewutil.CopyEnRuWith(fs, "view_rename.html", ruTranslation)
}

type editData struct {
	*viewutil.BaseData
	HyphaName string
	IsNew     bool
	Content   string
	Message   string
	Preview   template.HTML
}

func EditHypha(meta viewutil.Meta, hyphaName string, isNew bool, content string, message string, preview template.HTML) {
	viewutil.ExecutePage(meta, chainEditHypha, editData{
		BaseData: &viewutil.BaseData{
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

type deleteRenameData struct {
	*viewutil.BaseData
	HyphaName string
}

func RenameHypha(meta viewutil.Meta, hyphaName string) {
	viewutil.ExecutePage(meta, chainRenameHypha, deleteRenameData{
		BaseData: &viewutil.BaseData{
			Addr: "/rename/" + hyphaName,
		},
		HyphaName: hyphaName,
	})
}

func DeleteHypha(meta viewutil.Meta, hyphaName string) {
	viewutil.ExecutePage(meta, chainDeleteHypha, deleteRenameData{
		BaseData: &viewutil.BaseData{
			Addr: "/delete/" + hyphaName,
		},
		HyphaName: hyphaName,
	})
}

type emptyHyphaData struct {
	Meta              viewutil.Meta
	HyphaName         string
	AllowRegistration bool
	UseAuth           bool
}

func EmptyHypha(meta viewutil.Meta, hyphaName string) string {
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

func NaviTitle(meta viewutil.Meta, hyphaName string) string {
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
