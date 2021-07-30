package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"github.com/bouncepaw/mycorrhiza/cfg"
)

func initReaders(r *mux.Router) {
	r.PathPrefix("/page/").HandlerFunc(handlerHypha)
	r.PathPrefix("/hypha/").HandlerFunc(handlerHypha)
	r.PathPrefix("/text/").HandlerFunc(handlerText)
	r.PathPrefix("/binary/").HandlerFunc(handlerBinary)
	r.PathPrefix("/rev/").HandlerFunc(handlerRevision)
	r.PathPrefix("/primitive-diff/").HandlerFunc(handlerPrimitiveDiff)
	r.PathPrefix("/attachment/").HandlerFunc(handlerAttachment)
}

func handlerAttachment(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		hyphaName = util.HyphaNameFromRq(rq, "attachment")
		h         = hyphae.ByName(hyphaName)
		u         = user.FromRequest(rq)
	)
	util.HTTP200Page(w,
		views.BaseHTML(
			fmt.Sprintf("Attachment of %s", util.BeautifulName(hyphaName)),
			views.AttachmentMenuHTML(rq, h, u),
			u))
}

func handlerPrimitiveDiff(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		shorterUrl      = strings.TrimPrefix(rq.URL.Path, "/primitive-diff/")
		firstSlashIndex = strings.IndexRune(shorterUrl, '/')
		revHash         = shorterUrl[:firstSlashIndex]
		hyphaName       = util.CanonicalName(shorterUrl[firstSlashIndex+1:])
		h               = hyphae.ByName(hyphaName)
		u               = user.FromRequest(rq)
	)
	util.HTTP200Page(w,
		views.BaseHTML(
			fmt.Sprintf("Diff of %s at %s", hyphaName, revHash),
			views.PrimitiveDiffHTML(rq, h, u, revHash),
			u))
}

// handlerRevision displays a specific revision of text part a page
func handlerRevision(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		shorterUrl        = strings.TrimPrefix(rq.URL.Path, "/rev/")
		firstSlashIndex   = strings.IndexRune(shorterUrl, '/')
		revHash           = shorterUrl[:firstSlashIndex]
		hyphaName         = util.CanonicalName(shorterUrl[firstSlashIndex+1:])
		h                 = hyphae.ByName(hyphaName)
		contents          = fmt.Sprintf(`<p>This hypha had no text at this revision.</p>`)
		textContents, err = history.FileAtRevision(h.TextPartPath(), revHash)
		u                 = user.FromRequest(rq)
	)
	if err == nil {
		ctx, _ := mycocontext.ContextFromStringInput(hyphaName, textContents)
		contents = mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx))
	}
	page := views.RevisionHTML(
		rq,
		h,
		contents,
		revHash,
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, views.BaseHTML(util.BeautifulName(hyphaName), page, u))
}

// handlerText serves raw source text of the hypha.
func handlerText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	hyphaName := util.HyphaNameFromRq(rq, "text")
	if h := hyphae.ByName(hyphaName); h.Exists {
		log.Println("Serving", h.TextPath)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, rq, h.TextPath)
	}
}

// handlerBinary serves attachment of the hypha.
func handlerBinary(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	hyphaName := util.HyphaNameFromRq(rq, "binary")
	if h := hyphae.ByName(hyphaName); h.Exists {
		log.Println("Serving", h.BinaryPath)
		w.Header().Set("Content-Type", mimetype.FromExtension(filepath.Ext(h.BinaryPath)))
		http.ServeFile(w, rq, h.BinaryPath)
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
	)
	if h.Exists {
		fileContentsT, errT := os.ReadFile(h.TextPath)
		_, errB := os.Stat(h.BinaryPath)
		if errT == nil {
			ctx, _ := mycocontext.ContextFromStringInput(hyphaName, string(fileContentsT))
			ctx = mycocontext.WithWebSiteURL(ctx, cfg.URL)
			getOpenGraph, descVisitor, imgVisitor := mycomarkup.OpenGraphVisitors(ctx)
			ast := mycomarkup.BlockTree(ctx, descVisitor, imgVisitor)
			contents = mycomarkup.BlocksToHTML(ctx, ast)
			openGraph = getOpenGraph()
		}
		if !os.IsNotExist(errB) {
			contents = views.AttachmentHTML(h) + contents
		}
	}
	if contents == "" {
		util.HTTP404Page(w,
			views.BaseHTML(
				util.BeautifulName(hyphaName),
				views.HyphaHTML(rq, h, contents),
				u,
				openGraph))
	} else {
		util.HTTP200Page(w,
			views.BaseHTML(
				util.BeautifulName(hyphaName),
				views.HyphaHTML(rq, h, contents),
				u,
				openGraph))
	}
}
