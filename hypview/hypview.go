package hypview

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"strings"
	"text/template"
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
`
	chainNaviTitle  viewutil.Chain
	chainEmptyHypha viewutil.Chain
)

func Init() {
	chainNaviTitle = viewutil.
		En(viewutil.CopyEnWith(fs, "view_navititle.html")).
		Ru(viewutil.CopyRuWith(fs, "view_navititle.html")) // no text inside
	chainEmptyHypha = viewutil.
		En(viewutil.CopyEnWith(fs, "view_empty_hypha.html")).
		Ru(template.Must(viewutil.CopyRuWith(fs, "view_empty_hypha.html").Parse(ruTranslation)))
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
