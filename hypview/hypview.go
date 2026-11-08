package hypview

import (
	"html/template"
	"log/slog"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
)

var (
	tmplNaviTitle   *template.Template
	pageRenameHypha *newtmpl.Page
)

func Init() {
	tmplNaviTitle = newtmpl.NewPage("navititle.html").Template
	pageRenameHypha = newtmpl.NewPage("hypha-rename.html")
}

type renameData struct {
	newtmpl.BaseData
	HyphaName               string
	LeaveRedirectionDefault bool
}

func RenameHypha(meta viewdata.Meta, hyphaName string) {
	if err := pageRenameHypha.RenderDataTo(meta, &renameData{
		BaseData: newtmpl.BaseData{
			Addr: "/rename/" + hyphaName,
		},
		HyphaName:               hyphaName,
		LeaveRedirectionDefault: backlinks.BacklinksCount(hyphaName) != 0,
	}); err != nil {
		slog.Error("Failed to execute page", "page", "rename", "err", err)
	}
}

type naviTitleData struct {
	HyphaNameParts            []string
	HyphaNamePartsWithParents []string
	Icon                      template.HTML
	HomeHypha                 string
}

func NaviTitle(meta viewdata.Meta, hyphaName string) template.HTML {
	parts, partsWithParents := naviTitleify(hyphaName)
	var buf strings.Builder
	err := tmplNaviTitle.ExecuteTemplate(&buf, "navititle", naviTitleData{
		HyphaNameParts:            parts,
		HyphaNamePartsWithParents: partsWithParents,
		Icon:                      template.HTML(cfg.NaviTitleIcon),
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
