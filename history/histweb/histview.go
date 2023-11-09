// Package histweb provides web stuff for history
package histweb

import (
	"embed"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func InitHandlers(rtr *mux.Router) {
	rtr.PathPrefix("/primitive-diff/").HandlerFunc(handlerPrimitiveDiff)
	rtr.HandleFunc("/recent-changes/{count:[0-9]+}", handlerRecentChanges)
	rtr.HandleFunc("/recent-changes/", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/recent-changes/20", http.StatusSeeOther)
	})
	rtr.PathPrefix("/history/").HandlerFunc(handlerHistory)
	rtr.HandleFunc("/recent-changes-rss", handlerRecentChangesRSS)
	rtr.HandleFunc("/recent-changes-atom", handlerRecentChangesAtom)
	rtr.HandleFunc("/recent-changes-json", handlerRecentChangesJSON)

	chainPrimitiveDiff = viewutil.CopyEnRuWith(fs, "view_primitive_diff.html", ruTranslation)
	chainRecentChanges = viewutil.CopyEnRuWith(fs, "view_recent_changes.html", ruTranslation)
	chainHistory = viewutil.CopyEnRuWith(fs, "view_history.html", ruTranslation)
}

func handlerPrimitiveDiff(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	shorterURL := strings.TrimPrefix(rq.URL.Path, "/primitive-diff/")
	revHash, slug, found := strings.Cut(shorterURL, "/")
	if !found || !util.IsRevHash(revHash) || len(slug) < 1 {
		http.Error(w, "403 bad request", http.StatusBadRequest)
		return
	}
	var (
		mycoFilePath string
		h            = hyphae.ByName(util.CanonicalName(slug))
	)
	switch h := h.(type) {
	case hyphae.ExistingHypha:
		mycoFilePath = h.TextFilePath()
	case *hyphae.EmptyHypha:
		mycoFilePath = filepath.Join(files.HyphaeDir(), h.CanonicalName()+".myco")
	}
	text, err := history.PrimitiveDiffAtRevision(mycoFilePath, revHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	primitiveDiff(viewutil.MetaFrom(w, rq), h, revHash, text)
}

// handlerRecentChanges displays the /recent-changes/ page.
func handlerRecentChanges(w http.ResponseWriter, rq *http.Request) {
	// Error ignored: filtered by regex
	editCount, _ := strconv.Atoi(mux.Vars(rq)["count"])
	if editCount > 100 {
		return
	}
	recentChanges(viewutil.MetaFrom(w, rq), editCount, history.RecentChanges(editCount))
}

// handlerHistory lists all revisions of a hypha.
func handlerHistory(w http.ResponseWriter, rq *http.Request) {
	hyphaName := util.HyphaNameFromRq(rq, "history")
	var list string

	// History can be found for files that do not exist anymore.
	revs, err := history.Revisions(hyphaName)
	if err == nil {
		list = history.WithRevisions(hyphaName, revs)
	}
	log.Println("Found", len(revs), "revisions for", hyphaName)

	historyView(viewutil.MetaFrom(w, rq), hyphaName, list)
}

// genericHandlerOfFeeds is a helper function for the web feed handlers.
func genericHandlerOfFeeds(w http.ResponseWriter, rq *http.Request, f func(history.FeedOptions) (string, error), name string, contentType string) {
	opts, err := history.ParseFeedOptions(rq.URL.Query())
	var content string
	if err == nil {
		content, err = f(opts)
	}

	if err != nil {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "An error while generating "+name+": "+err.Error())
	} else {
		w.Header().Set("Content-Type", fmt.Sprintf("%s;charset=utf-8", contentType))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, content)
	}
}

func handlerRecentChangesRSS(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesRSS, "RSS", "application/rss+xml")
}

func handlerRecentChangesAtom(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesAtom, "Atom", "application/atom+xml")
}

func handlerRecentChangesJSON(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesJSON, "JSON feed", "application/feed+json")
}

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "history of title"}}История «{{.}}»{{end}}
{{define "history of heading"}}История <a href="/hypha/{{.}}">{{beautifulName .}}</a>{{end}}

{{define "diff for at title"}}Разница для {{beautifulName .HyphaName}} для {{.Hash}}{{end}}
{{define "diff for at heading"}}Разница для <a href="/hypha/{{.HyphaName}}">{{beautifulName .HyphaName}}</a> для {{.Hash}}{{end}}
{{define "no text diff available"}}Нет текстовой разницы.{{end}}

{{define "count pre"}}Отобразить{{end}}
{{define "count post"}}свежих правок.{{end}}
{{define "subscribe via"}}Подписаться через <a href="/recent-changes-rss">RSS</a>, <a href="/recent-changes-atom">Atom</a> или <a href="/recent-changes-json">JSON-ленту</a>.{{end}}
{{define "recent changes"}}Свежие правки{{end}}
{{define "n recent changes"}}{{.}} свеж{{if eq . 1}}ая правка{{else if le . 4}}их правок{{else}}их правок{{end}}{{end}}
{{define "recent empty"}}Правки не найдены.{{end}}
`
	chainPrimitiveDiff, chainRecentChanges, chainHistory viewutil.Chain
)

type recentChangesData struct {
	*viewutil.BaseData
	EditCount int
	Changes   []history.Revision
	UserHypha string
	Stops     []int
}

func recentChanges(meta viewutil.Meta, editCount int, changes []history.Revision) {
	viewutil.ExecutePage(meta, chainRecentChanges, recentChangesData{
		BaseData:  &viewutil.BaseData{},
		EditCount: editCount,
		Changes:   changes,
		UserHypha: cfg.UserHypha,
		Stops:     []int{20, 50, 100},
	})
}

type primitiveDiffData struct {
	*viewutil.BaseData
	HyphaName string
	Hash      string
	Text      template.HTML
}

func primitiveDiff(meta viewutil.Meta, h hyphae.Hypha, hash, text string) {
	hunks := history.SplitPrimitiveDiff(text)
	if len(hunks) > 0 {
		var buf strings.Builder
		for _, hunk := range hunks {
			lines := strings.Split(hunk, "\n")
			buf.WriteString(`<pre class="codeblock">`)
			for i, line := range lines {
				line = strings.Trim(line, "\r")
				var class string
				if len(line) > 0 {
					switch line[0] {
					case '+':
						class = "primitive-diff__addition"
					case '-':
						class = "primitive-diff__deletion"
					case '@':
						class = "primitive-diff__context"
					}
				}
				if i > 0 {
					buf.WriteString("\n")
				}
				line = template.HTMLEscapeString(line)
				fmt.Fprintf(&buf, `<code class="%s">%s</code>`,
					class, line)
			}
			buf.WriteString(`</pre>`)
		}
		text = buf.String()
	} else if text != "" {
		text = template.HTMLEscapeString(text)
		text = fmt.Sprintf(
			`<pre class="codeblock"><code>%s</code></pre>`, text)
	}
	viewutil.ExecutePage(meta, chainPrimitiveDiff, primitiveDiffData{
		BaseData:  &viewutil.BaseData{},
		HyphaName: h.CanonicalName(),
		Hash:      hash,
		Text:      template.HTML(text),
	})
}

type historyData struct {
	*viewutil.BaseData
	HyphaName string
	Contents  string
}

func historyView(meta viewutil.Meta, hyphaName, contents string) {
	viewutil.ExecutePage(meta, chainHistory, historyData{
		BaseData: &viewutil.BaseData{
			Addr: "/history/" + util.CanonicalName(hyphaName),
		},
		HyphaName: hyphaName,
		Contents:  contents,
	})
}
