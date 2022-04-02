package backlinks

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

func InitHandlers(rtr *mux.Router) {
	rtr.PathPrefix("/backlinks/").HandlerFunc(handlerBacklinks)
	chain = viewutil.
		En(viewutil.CopyEnWith(fs, "view_backlinks.html")).
		Ru(template.Must(viewutil.CopyRuWith(fs, "view_backlinks.html").Parse(ruTranslation)))
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

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "backlinks to text"}}Обратные ссылки на {{.}}{{end}}
{{define "backlinks to link"}}Обратные ссылки на <a href="/hypha/{{.}}">{{beautifulName .}}</a>{{end}}
{{define "description"}}Ниже перечислены гифы, на которых есть ссылка на эту гифу, трансклюзия этой гифы или эта гифа вставлена как изображение.{{end}}
`
	chain viewutil.Chain
)

type backlinksData struct {
	viewutil.BaseData
	HyphaName string
	Backlinks []string
}

func viewBacklinks(meta viewutil.Meta, hyphaName string, backlinks []string) {
	if err := chain.Get(meta).ExecuteTemplate(meta.W, "page", backlinksData{
		BaseData: viewutil.BaseData{
			Meta:          meta,
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		HyphaName: hyphaName,
		Backlinks: backlinks,
	}); err != nil {
		log.Println(err)
	}
}
