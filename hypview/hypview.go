package hypview

import (
	"embed"
	"html/template"
	"log/slog"
	"strings"
	texttemplate "text/template"

	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

var (
	//go:embed *.html
	fs              embed.FS
	tmplNaviTitle   *texttemplate.Template
	tmplRenameHypha *texttemplate.Template
)

func Init() {
	tmplNaviTitle = viewutil.CopyWith(fs, "view_navititle.html")
	tmplRenameHypha = viewutil.CopyWith(fs, "view_rename.html")
}

type renameData struct {
	*viewutil.BaseData
	HyphaName               string
	LeaveRedirectionDefault bool
}

func RenameHypha(meta viewdata.Meta, hyphaName string) {
	viewutil.ExecutePage(meta, tmplRenameHypha, renameData{
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

func NaviTitle(meta viewdata.Meta, hyphaName string) template.HTML {
	parts, partsWithParents := naviTitleify(hyphaName)
	var buf strings.Builder
	err := tmplNaviTitle.ExecuteTemplate(&buf, "navititle", naviTitleData{
		HyphaNameParts:            parts,
		HyphaNamePartsWithParents: partsWithParents,
		Icon:                      cfg.NaviTitleIcon,
		HomeHypha:                 cfg.HomeHypha,
	})
	if err != nil {
		slog.Error("Failed to render NaviTitle properly; using nevertheless", "err", err)
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
