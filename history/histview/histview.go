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
	"log"
	"net/http"
	"strings"
	"text/template"
)

func InitHandlers(rtr *mux.Router) {
	rtr.PathPrefix("/primitive-diff/").HandlerFunc(handlerPrimitiveDiff)
	chainPrimitiveDiff = viewutil.
		En(viewutil.CopyEnWith(fs, "view_primitive_diff.html")).
		Ru(template.Must(viewutil.CopyRuWith(fs, "view_primitive_diff.html").Parse(ruTranslation)))
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

var (
	//go:embed *.html
	fs            embed.FS
	ruTranslation = `
{{define "diff for at title"}}Разница для {{beautifulName .HyphaName}} для {{.Hash}}{{end}}
{{define "diff for at heading"}}Разница для <a href="/hypha/{{.HyphaName}}">{{beautifulName .HyphaName}}</a> для {{.Hash}}{{end}}

`
	chainPrimitiveDiff viewutil.Chain
)

type primitiveDiffData struct {
	viewutil.BaseData
	HyphaName string
	Hash      string
	Text      string
}

func primitiveDiff(meta viewutil.Meta, h hyphae.ExistingHypha, hash, text string) {
	if err := chainPrimitiveDiff.Get(meta).ExecuteTemplate(meta.W, "page", primitiveDiffData{
		BaseData: viewutil.BaseData{
			Meta:          meta,
			Addr:          "/primitive-diff/" + hash + "/" + h.CanonicalName(),
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		HyphaName: h.CanonicalName(),
		Hash:      hash,
		Text:      text,
	}); err != nil {
		log.Println(err)
	}
}
