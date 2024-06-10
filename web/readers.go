package web

import (
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"github.com/bouncepaw/mycorrhiza/categories"
	"github.com/bouncepaw/mycorrhiza/hypview"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/mimetype"
	"github.com/bouncepaw/mycorrhiza/internal/tree"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/mycoopts"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/tools"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/util"
)

func initReaders(r *mux.Router) {
	r.PathPrefix("/page/").HandlerFunc(handlerHypha)
	r.PathPrefix("/hypha/").HandlerFunc(handlerHypha)
	r.PathPrefix("/text/").HandlerFunc(handlerText)
	r.PathPrefix("/binary/").HandlerFunc(handlerBinary)
	r.PathPrefix("/rev/").HandlerFunc(handlerRevision)
	r.PathPrefix("/rev-text/").HandlerFunc(handlerRevisionText)
	r.PathPrefix("/media/").HandlerFunc(handlerMedia)
	r.Path("/today").HandlerFunc(handlerToday)
	r.Path("/edit-today").HandlerFunc(handlerEditToday)

	// Backlinks
	r.PathPrefix("/backlinks/").HandlerFunc(handlerBacklinks)
	r.PathPrefix("/orphans").HandlerFunc(handlerOrphans)
}

func handlerEditToday(w http.ResponseWriter, rq *http.Request) {
	today := time.Now().Format(time.DateOnly)
	http.Redirect(w, rq, "/edit/"+today, http.StatusSeeOther)
}

func handlerToday(w http.ResponseWriter, rq *http.Request) {
	today := time.Now().Format(time.DateOnly)
	http.Redirect(w, rq, "/hypha/"+today, http.StatusSeeOther)
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
		viewutil.Base(
			viewutil.MetaFrom(w, rq),
			lc.Get("ui.media_title", &l18n.Replacements{"name": util.BeautifulName(hyphaName)}),
			hypview.MediaMenu(rq, h, u),
			map[string]string{},
		))
}

// handlerRevisionText sends Mycomarkup text of the hypha at the given revision. See also: handlerRevision, handlerText.
//
// /rev-text/<revHash>/<hyphaName>
func handlerRevisionText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	shorterURL := strings.TrimPrefix(rq.URL.Path, "/rev-text/")
	revHash, slug, found := strings.Cut(shorterURL, "/")
	if !found || !util.IsRevHash(revHash) || len(slug) < 1 {
		http.Error(w, "403 bad request", http.StatusBadRequest)
		return
	}
	var (
		hyphaName = util.CanonicalName(slug)
		h         = hyphae.ByName(hyphaName)
	)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		var mycoFilePath = filepath.Join(files.HyphaeDir(), h.CanonicalName()+".myco")
		var textContents, err = history.FileAtRevision(mycoFilePath, revHash)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("While serving text of ‘%s’ at revision ‘%s’: %s\n", hyphaName, revHash, err.Error())
			_, _ = io.WriteString(w, "Error: "+err.Error())
			return
		}
		log.Printf("Serving text of ‘%s’ from ‘%s’ at revision ‘%s’\n", hyphaName, mycoFilePath, revHash)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, textContents)
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
	lc := l18n.FromRequest(rq)
	shorterURL := strings.TrimPrefix(rq.URL.Path, "/rev/")
	revHash, slug, found := strings.Cut(shorterURL, "/")
	if !found || !util.IsRevHash(revHash) || len(slug) < 1 {
		http.Error(w, "403 bad request", http.StatusBadRequest)
		return
	}
	var (
		hyphaName    = util.CanonicalName(slug)
		h            = hyphae.ByName(hyphaName)
		contents     = fmt.Sprintf(`<p>%s</p>`, lc.Get("ui.revision_no_text"))
		textContents string
		err          error
		mycoFilePath string
	)
	switch h := h.(type) {
	case hyphae.ExistingHypha:
		mycoFilePath = h.TextFilePath()
	case *hyphae.EmptyHypha:
		mycoFilePath = filepath.Join(files.HyphaeDir(), h.CanonicalName()+".myco")
	}
	textContents, err = history.FileAtRevision(mycoFilePath, revHash)
	if err == nil {
		ctx, _ := mycocontext.ContextFromStringInput(textContents, mycoopts.MarkupOptions(hyphaName))
		contents = mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx))
	}

	page := hypview.Revision(
		viewutil.MetaFrom(w, rq),
		h,
		contents,
		revHash,
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(
		w,
		viewutil.Base(
			viewutil.MetaFrom(w, rq),
			lc.Get("ui.revision_title", &l18n.Replacements{"name": util.BeautifulName(hyphaName), "rev": revHash}),
			page,
			map[string]string{},
		),
	)
}

// handlerText serves raw source text of the hypha.
func handlerText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	hyphaName := util.HyphaNameFromRq(rq, "text")
	switch h := hyphae.ByName(hyphaName).(type) {
	case hyphae.ExistingHypha:
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
		lc        = l18n.FromRequest(rq)
	)

	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		// contents = hypview.EmptyHypha()
		util.HTTP404Page(w,
			viewutil.Base(
				viewutil.MetaFrom(w, rq),
				util.BeautifulName(hyphaName),
				hypview.Hypha(viewutil.MetaFrom(w, rq), h, ""),
				map[string]string{},
				openGraph))
	case hyphae.ExistingHypha:
		fileContentsT, errT := os.ReadFile(h.TextFilePath())
		if errT == nil {
			ctx, _ := mycocontext.ContextFromStringInput(string(fileContentsT), mycoopts.MarkupOptions(hyphaName))
			getOpenGraph, descVisitor, imgVisitor := tools.OpenGraphVisitors(ctx)
			ast := mycomarkup.BlockTree(ctx, descVisitor, imgVisitor)
			contents = mycomarkup.BlocksToHTML(ctx, ast)
			openGraph = getOpenGraph()
		}
		switch h := h.(type) {
		case *hyphae.MediaHypha:
			contents = mycoopts.Media(h, lc) + contents
		}

		meta := viewutil.MetaFrom(w, rq)
		category_list := ":" + strings.Join(categories.CategoriesWithHypha(h.CanonicalName()), ":") + ":"
		subhyphae, prevHyphaName, nextHyphaName := tree.Tree(h.CanonicalName())
		isMyProfile := cfg.UseAuth && util.IsProfileName(h.CanonicalName()) && meta.U.Name == strings.TrimPrefix(h.CanonicalName(), cfg.UserHypha+"/")

		_ = pageHypha.RenderTo(
			meta,
			map[string]any{
				"SubhyphaeHTML": subhyphae,
				"PrevHyphaName": prevHyphaName,
				"NextHyphaName": nextHyphaName,
				"IsMyProfile":   isMyProfile,
				"NaviTitle":     hypview.NaviTitle(meta, h.CanonicalName()),
				"Contents":      template.HTML(contents),
			})
		util.HTTP200Page(w,
			viewutil.Base(
				viewutil.MetaFrom(w, rq),
				util.BeautifulName(hyphaName),
				hypview.Hypha(viewutil.MetaFrom(w, rq), h, contents),
				map[string]string{"cats": category_list},
				openGraph))
	}
}

// handlerBacklinks lists all backlinks to a hypha.
func handlerBacklinks(w http.ResponseWriter, rq *http.Request) {
	hyphaName := util.HyphaNameFromRq(rq, "backlinks")

	_ = pageBacklinks.RenderTo(viewutil.MetaFrom(w, rq),
		map[string]any{
			"Addr":      "/backlinks/" + hyphaName,
			"HyphaName": hyphaName,
			"Backlinks": backlinks.BacklinksFor(hyphaName),
		})
}

func handlerOrphans(w http.ResponseWriter, rq *http.Request) {
	_ = pageOrphans.RenderTo(viewutil.MetaFrom(w, rq),
		map[string]any{
			"Addr":    "/orphans",
			"Orphans": backlinks.Orphans(),
		})
}
