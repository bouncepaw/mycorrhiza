package backlinks

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/newtmpl"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

func InitHandlers(rtr *mux.Router) {
	rtr.PathPrefix("/backlinks/").HandlerFunc(handlerBacklinks)
	rtr.PathPrefix("/orphans").HandlerFunc(handlerOrphans)
	chainBacklinks = viewutil.CopyEnRuWith(fs, "view_backlinks.html", ruTranslation)
	pageOrphans = newtmpl.NewPage(fs, "view_orphans.html", map[string]string{
		"orphaned hyphae":    "Гифы-сироты",
		"orphan description": "Ниже перечислены гифы без ссылок на них.",
	})
}

// handlerBacklinks lists all backlinks to a hypha.
func handlerBacklinks(w http.ResponseWriter, rq *http.Request) {
	var (
		hyphaName = util.HyphaNameFromRq(rq, "backlinks")
		backlinks []string
	)
	for b := range yieldHyphaBacklinks(hyphaName) {
		backlinks = append(backlinks, b)
	}
	viewBacklinks(viewutil.MetaFrom(w, rq), hyphaName, backlinks)
}

func handlerOrphans(w http.ResponseWriter, rq *http.Request) {
	var orphans []string
	for h := range hyphae.YieldExistingHyphae() {
		if BacklinksCount(h.CanonicalName()) == 0 {
			orphans = append(orphans, h.CanonicalName())
		}
	}
	sort.Strings(orphans)

	_ = pageOrphans.RenderTo(viewutil.MetaFrom(w, rq),
		map[string]any{
			"Addr":    "/orphans",
			"Orphans": orphans,
		})
}

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "backlinks to text"}}Обратные ссылки на {{.}}{{end}}
{{define "backlinks to link"}}Обратные ссылки на <a href="/hypha/{{.}}">{{beautifulName .}}</a>{{end}}
{{define "description"}}Ниже перечислены гифы, на которых есть ссылка на эту гифу, трансклюзия этой гифы или эта гифа вставлена как изображение.{{end}}
`
	chainBacklinks viewutil.Chain

	pageOrphans *newtmpl.Page
)

type backlinksData struct {
	*viewutil.BaseData
	HyphaName string
	Backlinks []string
}

func viewBacklinks(meta viewutil.Meta, hyphaName string, backlinks []string) {
	viewutil.ExecutePage(meta, chainBacklinks, backlinksData{
		BaseData: &viewutil.BaseData{
			Addr: "/backlinks/" + hyphaName,
		},
		HyphaName: hyphaName,
		Backlinks: backlinks,
	})
}
