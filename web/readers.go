package web

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v3"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"github.com/bouncepaw/mycomarkup/v3/tools"
)

func initReaders(r *mux.Router) {
	r.PathPrefix("/page/").HandlerFunc(handlerHypha)
	r.PathPrefix("/hypha/").HandlerFunc(handlerHypha)
	r.PathPrefix("/text/").HandlerFunc(handlerText)
	r.PathPrefix("/binary/").HandlerFunc(handlerBinary)
	r.PathPrefix("/rev/").HandlerFunc(handlerRevision)
	r.PathPrefix("/rev-text/").HandlerFunc(handlerRevisionText)
	r.PathPrefix("/primitive-diff/").HandlerFunc(handlerPrimitiveDiff)
	r.PathPrefix("/media/").HandlerFunc(handlerMedia)
}

func handlerMedia(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		hyphaName = util.HyphaNameFromRq(rq, "media")
		h         = hyphae.ByName(hyphaName)
		u         = user.FromRequest(rq)
		lc        = l18n.FromRequest(rq)
	)
	util.HTTP200Page(w,
		views.Base(
			lc.Get("ui.media_title", &l18n.Replacements{"name": util.BeautifulName(hyphaName)}),
			views.MediaMenu(rq, h, u),
			lc,
			u))
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
	var (
		hyphaName = util.CanonicalName(slug)
		h         = hyphae.ByName(hyphaName)
		user      = user.FromRequest(rq)
		locale    = l18n.FromRequest(rq)
	)
	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "404 not found")
	case hyphae.ExistingHypha:
		util.HTTP200Page(w, views.Base(
			locale.Get("ui.diff_title", &l18n.Replacements{"name": util.BeautifulName(hyphaName), "rev": revHash}),
			views.PrimitiveDiff(rq, h, user, revHash), locale, user))
	}
}

// handlerRevisionText sends Mycomarkup text of the hypha at the given revision. See also: handlerRevision, handlerText.
//
// /rev-text/<revHash>/<hyphaName>
func handlerRevisionText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		shorterURL      = strings.TrimPrefix(rq.URL.Path, "/rev-text/")
		firstSlashIndex = strings.IndexRune(shorterURL, '/')
		revHash         = shorterURL[:firstSlashIndex]
		hyphaName       = util.CanonicalName(shorterURL[firstSlashIndex+1:])
		h               = hyphae.ByName(hyphaName)
	)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		log.Printf(`Hypha ‘%s’ does not exist`)
		w.WriteHeader(http.StatusNotFound)
	case hyphae.ExistingHypha:
		if !h.HasTextFile() {
			log.Printf(`Media hypha ‘%s’ has no text`)
			w.WriteHeader(http.StatusNotFound)
		}
		var textContents, err = history.FileAtRevision(h.TextFilePath(), revHash)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("While serving text of ‘%s’ at revision ‘%s’: %s\n", hyphaName, revHash, err.Error())
			_, _ = io.WriteString(w, "Error: "+err.Error())
			return
		}
		log.Printf("Serving text of ‘%s’ from ‘%s’ at revision ‘%s’\n", hyphaName, h.TextFilePath(), revHash)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, textContents)
	}
}

// handlerRevision displays a specific revision of the text part the hypha
func handlerRevision(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		lc              = l18n.FromRequest(rq)
		shorterURL      = strings.TrimPrefix(rq.URL.Path, "/rev/")
		firstSlashIndex = strings.IndexRune(shorterURL, '/')
		revHash         = shorterURL[:firstSlashIndex]
		hyphaName       = util.CanonicalName(shorterURL[firstSlashIndex+1:])
		h               = hyphae.ByName(hyphaName)
		contents        = fmt.Sprintf(`<p>%s</p>`, lc.Get("ui.revision_no_text"))
		u               = user.FromRequest(rq)
	)
	switch h := h.(type) {
	case hyphae.ExistingHypha:
		var textContents, err = history.FileAtRevision(h.TextFilePath(), revHash)

		if err == nil {
			ctx, _ := mycocontext.ContextFromStringInput(hyphaName, textContents)
			contents = mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx))
		}
	}
	page := views.Revision(
		rq,
		lc,
		h,
		contents,
		revHash,
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(
		w,
		views.Base(
			lc.Get("ui.revision_title", &l18n.Replacements{"name": util.BeautifulName(hyphaName), "rev": revHash}),
			page,
			lc,
			u,
		),
	)
}

// handlerText serves raw source text of the hypha.
func handlerText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	hyphaName := util.HyphaNameFromRq(rq, "text")
	switch h := hyphae.ByName(hyphaName).(type) {
	case *hyphae.TextualHypha:
		log.Println("Serving", h.TextFilePath())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, rq, h.TextFilePath())
	}
}

// handlerBinary serves attachment of the hypha.
func handlerBinary(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	hyphaName := util.HyphaNameFromRq(rq, "binary")
	switch h := hyphae.ByName(hyphaName).(type) {
	case *hyphae.EmptyHypha:
	case *hyphae.TextualHypha:
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Textual hypha ‘%s’ has no media, cannot serve\n", h.CanonicalName())
	case *hyphae.MediaHypha:
		log.Println("Serving", h.MediaFilePath())
		w.Header().Set("Content-Type", mimetype.FromExtension(filepath.Ext(h.MediaFilePath())))
		http.ServeFile(w, rq, h.MediaFilePath())
	}
}

// handlerHypha is the main hypha action that displays the hypha and the binary upload form along with some navigation.
func handlerHypha(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		hyphaName = util.HyphaNameFromRq(rq, "page", "hypha")
		h         = hyphae.ByName(hyphaName)
		contents  string
		openGraph string
		u         = user.FromRequest(rq)
		lc        = l18n.FromRequest(rq)
	)

	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		util.HTTP404Page(w,
			views.Base(
				util.BeautifulName(hyphaName),
				views.Hypha(views.MetaFrom(w, rq), h, contents),
				lc,
				u,
				openGraph))
	case hyphae.ExistingHypha:
		fileContentsT, errT := os.ReadFile(h.TextFilePath())
		if errT == nil {
			ctx, _ := mycocontext.ContextFromStringInput(hyphaName, string(fileContentsT))
			ctx = mycocontext.WithWebSiteURL(ctx, cfg.URL)
			getOpenGraph, descVisitor, imgVisitor := tools.OpenGraphVisitors(ctx)
			ast := mycomarkup.BlockTree(ctx, descVisitor, imgVisitor)
			contents = mycomarkup.BlocksToHTML(ctx, ast)
			openGraph = getOpenGraph()
		}
		switch h := h.(type) {
		case *hyphae.MediaHypha:
			contents = views.Media(h, lc) + contents
		}

		util.HTTP200Page(w,
			views.Base(
				util.BeautifulName(hyphaName),
				views.Hypha(views.MetaFrom(w, rq), h, contents),
				lc,
				u,
				openGraph))
	}
}
