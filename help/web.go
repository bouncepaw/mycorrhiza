package help

// stuff.go is used for meta stuff about the wiki or all hyphae at once.
import (
	"github.com/bouncepaw/mycomarkup/v5"
	"github.com/bouncepaw/mycorrhiza/mycoopts"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"

	"github.com/bouncepaw/mycomarkup/v5/mycocontext"
)

var (
	chain         viewutil.Chain
	ruTranslation = `
{{define "title"}}Справка{{end}}
{{define "entry not found"}}Статья не найдена{{end}}
{{define "entry not found invitation"}}Если вы хотите написать эту статью сами, то будем рады вашим правкам <a class="wikilink wikilink_external wikilink_https" href="https://github.com/bouncepaw/mycorrhiza">в репозитории Миокризы</a>.{{end}}

{{define "topics"}}Темы справки{{end}}
{{define "main"}}Введение{{end}}
{{define "hypha"}}Гифа{{end}}
{{define "media"}}Медиа{{end}}
{{define "mycomarkup"}}Микоразметка{{end}}
{{define "category"}}Категории{{end}}
{{define "interface"}}Интерфейс{{end}}
{{define "prevnext"}}Пред/след{{end}}
{{define "top_bar"}}Верхняя панель{{end}}
{{define "sibling_hyphae"}}Гифы-сиблинги{{end}}
{{define "special pages"}}Специальные страницы{{end}}
{{define "recent_changes"}}Недавние изменения{{end}}
{{define "feeds"}}Ленты{{end}}
{{define "orphans"}}Гифы-сироты{{end}}
{{define "configuration"}}Конфигурация (для администраторов){{end}}
{{define "config_file"}}Файл конфигурации{{end}}
{{define "lock"}}Замок{{end}}
{{define "whitelist"}}Белый список{{end}}
{{define "telegram"}}Вход через Телеграм{{end}}
`
)

func InitHandlers(r *mux.Router) {
	r.PathPrefix("/help").HandlerFunc(handlerHelp)
	chain = viewutil.CopyEnRuWith(fs, "view_help.html", ruTranslation)
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
		viewHelp(meta, lang, "", articlePath)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		viewHelp(meta, lang, err.Error(), articlePath)
		return
	}

	// TODO: change for the function that uses byte array when there is such function in mycomarkup.
	ctx, _ := mycocontext.ContextFromStringInput(string(content), mycoopts.MarkupOptions(articlePath))
	ast := mycomarkup.BlockTree(ctx)
	result := mycomarkup.BlocksToHTML(ctx, ast)
	w.WriteHeader(http.StatusOK)
	viewHelp(meta, lang, result, articlePath)
}

type helpData struct {
	*viewutil.BaseData
	ContentsHTML string
	Lang         string
}

func viewHelp(meta viewutil.Meta, lang, contentsHTML, articlePath string) {
	viewutil.ExecutePage(meta, chain, helpData{
		BaseData: &viewutil.BaseData{
			Addr: "/help/" + articlePath,
		},
		ContentsHTML: contentsHTML,
		Lang:         lang,
	})
}
