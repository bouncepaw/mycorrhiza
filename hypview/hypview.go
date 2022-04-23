package hypview

import (
	"embed"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"strings"
	"text/template"
)

var (
	//go:embed *.html
	fs             embed.FS
	ruTranslation  = ``
	chainNaviTitle viewutil.Chain
)

func Init() {
	chainNaviTitle = viewutil.
		En(viewutil.CopyEnWith(fs, "view_navititle.html")).
		Ru(template.Must(viewutil.CopyRuWith(fs, "view_navititle.html").Parse(ruTranslation)))
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

	fmt.Printf("%q — %q\n", parts, partsWithParents)
	return parts, partsWithParents
}
