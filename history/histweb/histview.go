// Package histweb provides web stuff for history
package histweb

import (
	"embed"
	"encoding/hex"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	if !found || len(revHash) < 7 || len(slug) < 1 {
		http.Error(w, "403 bad request", http.StatusBadRequest)
		return
	}
	paddedRevHash := revHash
	if len(paddedRevHash)%2 != 0 {
		paddedRevHash = paddedRevHash[:len(paddedRevHash)-1]
	}
	if _, err := hex.DecodeString(paddedRevHash); err != nil {
		http.Error(w, "403 bad request", http.StatusBadRequest)
		return
	}
	switch h := hyphae.ByName(util.CanonicalName(slug)).(type) {
	case *hyphae.EmptyHypha:
		http.Error(w, "404 not found", http.StatusNotFound)
	case hyphae.ExistingHypha:
		text, err := history.PrimitiveDiffAtRevision(h.TextFilePath(), revHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		primitiveDiff(viewutil.MetaFrom(w, rq), h, revHash, text)
	}
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
	genericHandlerOfFeeds(w, rq, history.RecentChangesJSON, "JSON feed", "application/json")
}

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "history of title"}}История «{{.}}»{{end}}
{{define "history of heading"}}История <a href="/hypha/{{.}}">{{beautifulName .}}</a>{{end}}

{{define "diff for at title"}}Разница для {{beautifulName .HyphaName}} для {{.Hash}}{{end}}
{{define "diff for at heading"}}Разница для <a href="/hypha/{{.HyphaName}}">{{beautifulName .HyphaName}}</a> для {{.Hash}}{{end}}

{{define "count pre"}}Отобразить{{end}}
{{define "count post"}}свежих правок{{end}}
{{define "subscribe via"}}Подписаться через <a href="/recent-changes-rss">RSS</a>, <a href="/recent-changes-atom">Atom</a> или <a href="/recent-changes-json">JSON-ленту</a>.{{end}}
{{define "recent changes"}}Свежие правки{{end}}
{{define "n recent changes"}}{{.}} недавн{{if eq . 1}}ее изменение{{else if le . 4}}их изменения{{else}}их изменений{{end}}{{end}}
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
	Text      string
}

func primitiveDiff(meta viewutil.Meta, h hyphae.ExistingHypha, hash, text string) {
	viewutil.ExecutePage(meta, chainPrimitiveDiff, primitiveDiffData{
		BaseData:  &viewutil.BaseData{},
		HyphaName: h.CanonicalName(),
		Hash:      hash,
		Text:      text,
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
