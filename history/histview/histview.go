// Package histview provides web stuff for history
package histview

import (
	"embed"
	"encoding/hex"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
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

	chainPrimitiveDiff = viewutil.CopyEnRuWith(fs, "view_primitive_diff.html", ruTranslation)
	chainRecentChanges = viewutil.CopyEnRuWith(fs, "view_recent_changes.html", ruTranslation)
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

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "diff for at title"}}Разница для {{beautifulName .HyphaName}} для {{.Hash}}{{end}}
{{define "diff for at heading"}}Разница для <a href="/hypha/{{.HyphaName}}">{{beautifulName .HyphaName}}</a> для {{.Hash}}{{end}}

`
	// TODO: translate recent changes
	chainPrimitiveDiff viewutil.Chain
	chainRecentChanges viewutil.Chain
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
