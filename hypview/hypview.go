package hypview

import (
	"embed"
	"html/template"
	"log"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "rename hypha?"}}Переименовать {{beautifulName .}}?{{end}}
{{define "rename [[hypha]]?"}}Переименовать <a href="/hypha/{{.}}">{{beautifulName .}}</a>?{{end}}
{{define "new name"}}Новое название:{{end}}
{{define "rename recursively"}}Также переименовать подгифы{{end}}
{{define "rename tip"}}Переименовывайте аккуратно. <a href="/help/en/rename">Документация на английском.</a>{{end}}
{{define "leave redirection"}}Оставить перенаправление{{end}}


`
	chainNaviTitle   viewutil.Chain
	chainRenameHypha viewutil.Chain
)

func Init() {
	chainNaviTitle = viewutil.CopyEnRuWith(fs, "view_navititle.html", "")
	chainRenameHypha = viewutil.CopyEnRuWith(fs, "view_rename.html", ruTranslation)
}

type renameData struct {
	*viewutil.BaseData
	HyphaName               string
	LeaveRedirectionDefault bool
}

func RenameHypha(meta viewutil.Meta, hyphaName string) {
	viewutil.ExecutePage(meta, chainRenameHypha, renameData{
		BaseData: &viewutil.BaseData{
			Addr: "/rename/" + hyphaName,
		},
		HyphaName:               hyphaName,
		LeaveRedirectionDefault: backlinks.BacklinksCount(hyphaName) != 0,
	})
}

type naviTitleData struct {
	HyphaNameParts            []string
	HyphaNamePartsWithParents []string
	Icon                      string
	HomeHypha                 string
}

func NaviTitle(meta viewutil.Meta, hyphaName string) template.HTML {
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
	return template.HTML(buf.String())
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
