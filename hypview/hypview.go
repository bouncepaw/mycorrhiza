package hypview

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"strings"
)

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
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
{{define "rename tip"}}Если вы переименуете эту гифу, сломаются все ссылки, ведущие на неё, а также исходящие относительные ссылки. Также вы потеряете всю текущую историю для нового названия. Переименовывайте аккуратно.{{end}}
`
	chainNaviTitle   viewutil.Chain
	chainEmptyHypha  viewutil.Chain
	chainDeleteHypha viewutil.Chain
	chainRenameHypha viewutil.Chain
)

func Init() {
	chainNaviTitle = viewutil.CopyEnRuWith(fs, "view_navititle.html", "")
	chainEmptyHypha = viewutil.CopyEnRuWith(fs, "view_empty_hypha.html", ruTranslation)
	chainDeleteHypha = viewutil.CopyEnRuWith(fs, "view_delete.html", ruTranslation)
	chainRenameHypha = viewutil.CopyEnRuWith(fs, "view_rename.html", ruTranslation)
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
