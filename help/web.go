package help

// stuff.go is used for meta stuff about the wiki or all hyphae at once.
import (
	"github.com/bouncepaw/mycomarkup/v4"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v4/mycocontext"
)

var (
	chain         viewutil.Chain
	ruTranslation = `
{{define "title"}}Справка{{end}}
{{define "entry not found"}}Статья не найдена{{end}}
{{define "entry not found invitation"}}Если вы хотите написать эту статью сами, то будем рады вашим правкам <a class="wikilink wikilink_external wikilink_https" href="https://github.com/bouncepaw/mycorrhiza">в репозитории Миокризы</a>.{{end}}
`
)

func InitHandlers(r *mux.Router) {
	r.PathPrefix("/help").HandlerFunc(handlerHelp)
	chain = viewutil.
		En(viewutil.CopyEnWith(fs, "view_help.html")).
		Ru(template.Must(viewutil.CopyRuWith(fs, "view_help.html").Parse(ruTranslation)))
}

// handlerHelp gets the appropriate documentation or tells you where you (personally) have failed.
func handlerHelp(w http.ResponseWriter, rq *http.Request) {
	// See the history of this file to resurrect the old algorithm that supported multiple languages
	var (
		meta        = viewutil.MetaFrom(w, rq)
		articlePath = strings.TrimPrefix(strings.TrimPrefix(rq.URL.Path, "/help/"), "/help")
		lang        = "en"
	)
	if articlePath == "" {
		articlePath = "en"
	}

	if !strings.HasPrefix(articlePath, "en") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, "404 Not found")
		return
	}

	content, err := Get(articlePath)
	if err != nil && strings.HasPrefix(err.Error(), "open") {
		w.WriteHeader(http.StatusNotFound)
		viewHelp(meta, lang, "")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		viewHelp(meta, lang, err.Error())
		return
	}

	// TODO: change for the function that uses byte array when there is such function in mycomarkup.
	ctx, _ := mycocontext.ContextFromStringInput(string(content), shroom.MarkupOptions(articlePath))
	ast := mycomarkup.BlockTree(ctx)
	result := mycomarkup.BlocksToHTML(ctx, ast)
	w.WriteHeader(http.StatusOK)
	viewHelp(meta, lang, result)
}

type helpData struct {
	viewutil.BaseData
	ContentsHTML   string
	HelpTopicsHTML string
	Lang           string
}

func viewHelp(meta viewutil.Meta, lang, contentsHTML string) {
	if err := chain.Get(meta).ExecuteTemplate(meta.W, "page", helpData{
		BaseData: viewutil.BaseData{
			Meta:          meta,
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		ContentsHTML:   contentsHTML,
		HelpTopicsHTML: views.HelpTopics(lang, meta.Lc),
		Lang:           lang,
	}); err != nil {
		log.Println(err)
	}
}
