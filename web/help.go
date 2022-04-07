package web

// stuff.go is used for meta stuff about the wiki or all hyphae at once.
import (
	"github.com/bouncepaw/mycomarkup/v4"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/help"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v4/mycocontext"
)

func initHelp(r *mux.Router) {
	r.PathPrefix("/help").HandlerFunc(handlerHelp)
}

// handlerHelp gets the appropriate documentation or tells you where you (personally) have failed.
func handlerHelp(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	articlePath := strings.TrimPrefix(strings.TrimPrefix(rq.URL.Path, "/help/"), "/help")
	// See the history of this file to resurrect the old algorithm that supported multiple languages
	lang := "en"
	if articlePath == "" {
		articlePath = "en"
	}

	if !strings.HasPrefix(articlePath, "en") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, "404 Not found")
		return
	}

	content, err := help.Get(articlePath)
	if err != nil && strings.HasPrefix(err.Error(), "open") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("help.entry_not_found"),
				views.Help(views.HelpEmptyError(lc), lang, lc),
			),
		)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				err.Error(),
				views.Help(err.Error(), lang, lc),
			),
		)
		return
	}

	// TODO: change for the function that uses byte array when there is such function in mycomarkup.
	ctx, _ := mycocontext.ContextFromStringInput(string(content), shroom.MarkupOptions(articlePath))
	ast := mycomarkup.BlockTree(ctx)
	result := mycomarkup.BlocksToHTML(ctx, ast)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(
		w,
		views.Base(
			viewutil.MetaFrom(w, rq),
			lc.Get("help.title"),
			views.Help(result, lang, lc),
		),
	)
}
