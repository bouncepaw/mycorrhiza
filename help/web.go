package help

// stuff.go is used for meta stuff about the wiki or all hyphae at once.
import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bouncepaw/mycorrhiza/mycoopts"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"

	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"

	"github.com/gorilla/mux"
)

var (
	pageHelp *newtmpl.Page
)

func InitHandlers(r *mux.Router) {
	r.PathPrefix("/help").HandlerFunc(handlerHelp)
	pageHelp = newtmpl.NewPage(fs, "view_help.html")
}

// handlerHelp gets the appropriate documentation or tells you where you (personally) have failed.
func handlerHelp(w http.ResponseWriter, rq *http.Request) {
	// See the history of this file to resurrect the old algorithm that supported multiple languages
	var (
		meta        = viewdata.MetaFrom(w, rq)
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
		viewHelp(meta, lang, template.HTML(template.HTMLEscapeString(err.Error())), articlePath)
		return
	}

	// TODO: change for the function that uses byte array when there is such function in mycomarkup.
	ctx, _ := mycocontext.ContextFromStringInput(string(content), mycoopts.MarkupOptions(articlePath))
	ast := mycomarkup.BlockTree(ctx)
	result := mycomarkup.BlocksToHTML(ctx, ast)
	w.WriteHeader(http.StatusOK)
	viewHelp(meta, lang, template.HTML(result), articlePath)
}

type helpData struct {
	newtmpl.BaseData
	ContentsHTML template.HTML
	Lang         string
}

func viewHelp(meta viewdata.Meta, lang string, contentsHTML template.HTML, articlePath string) {
	if err := pageHelp.RenderDataTo(meta, &helpData{
		BaseData: newtmpl.BaseData{
			Addr: "/help/" + articlePath,
		},
		ContentsHTML: contentsHTML,
		Lang:         lang,
	}); err != nil {
		slog.Error("Failed to execute page", "page", "help", "err", err)
	}
}
