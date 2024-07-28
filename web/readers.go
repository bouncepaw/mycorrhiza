package web

import (
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"github.com/bouncepaw/mycorrhiza/hypview"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/categories"
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
	"log/slog"
	"net/http"
	"os"
	"path"
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
		isMedia   = false

		mime     string
		fileSize int64
	)
	switch h := h.(type) {
	case *hyphae.MediaHypha:
		isMedia = true
		mime = mimetype.FromExtension(path.Ext(h.MediaFilePath()))

		fileinfo, err := os.Stat(h.MediaFilePath())
		if err != nil {
			slog.Error("failed to stat media file", "err", err)
			// no return
		}

		fileSize = fileinfo.Size()
	}
	_ = pageMedia.RenderTo(viewutil.MetaFrom(w, rq), map[string]any{
		"HyphaName":    h.CanonicalName(),
		"U":            u,
		"IsMediaHypha": isMedia,
		"MimeType":     mime,
		"FileSize":     fileSize,
	})
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
		contents     = template.HTML(fmt.Sprintf(`<p>%s</p>`, lc.Get("ui.revision_no_text")))
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
		contents = template.HTML(mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx)))
	}

	meta := viewutil.MetaFrom(w, rq)
	_ = pageRevision.RenderTo(meta, map[string]any{
		"ViewScripts": cfg.ViewScripts,
		"Contents":    contents,
		"RevHash":     revHash,
		"NaviTitle":   hypview.NaviTitle(meta, h.CanonicalName()),
		"HyphaName":   h.CanonicalName(),
	})
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
	case *hyphae.EmptyHypha, *hyphae.TextualHypha:
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
		hyphaName                               = util.HyphaNameFromRq(rq, "page", "hypha")
		h                                       = hyphae.ByName(hyphaName)
		contents                                template.HTML
		openGraph                               template.HTML
		lc                                      = l18n.FromRequest(rq)
		meta                                    = viewutil.MetaFrom(w, rq)
		subhyphae, prevHyphaName, nextHyphaName = tree.Tree(h.CanonicalName())
		cats                                    = categories.CategoriesWithHypha(h.CanonicalName())
		category_list                           = ":" + strings.Join(cats, ":") + ":"
		isMyProfile                             = cfg.UseAuth && util.IsProfileName(h.CanonicalName()) && meta.U.Name == strings.TrimPrefix(h.CanonicalName(), cfg.UserHypha+"/")

		data = map[string]any{
			"HyphaName":               h.CanonicalName(),
			"SubhyphaeHTML":           subhyphae,
			"PrevHyphaName":           prevHyphaName,
			"NextHyphaName":           nextHyphaName,
			"IsMyProfile":             isMyProfile,
			"NaviTitle":               hypview.NaviTitle(meta, h.CanonicalName()),
			"BacklinkCount":           backlinks.BacklinksCount(h.CanonicalName()),
			"GivenPermissionToModify": user.CanProceed(rq, "edit"),
			"Categories":              cats,
			"IsMediaHypha":            false,
		}
	)
	slog.Info("reading hypha", "name", h.CanonicalName(), "can edit", data["GivenPermissionToModify"])
	meta.BodyAttributes = map[string]string{
		"cats": category_list,
	}

	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		w.WriteHeader(http.StatusNotFound)
		data["Contents"] = ""
		_ = pageHypha.RenderTo(meta, data)
	case hyphae.ExistingHypha:
		fileContentsT, err := os.ReadFile(h.TextFilePath())
		if err == nil {
			ctx, _ := mycocontext.ContextFromStringInput(string(fileContentsT), mycoopts.MarkupOptions(hyphaName))
			getOpenGraph, descVisitor, imgVisitor := tools.OpenGraphVisitors(ctx)
			openGraph = template.HTML(getOpenGraph())
			ast := mycomarkup.BlockTree(ctx, descVisitor, imgVisitor)
			contents = template.HTML(mycomarkup.BlocksToHTML(ctx, ast))
		}
		switch h := h.(type) {
		case *hyphae.MediaHypha:
			contents = template.HTML(mycoopts.Media(h, lc)) + contents
			data["IsMediaHypha"] = true
		}

		data["Contents"] = contents
		meta.HeadElements = append(meta.HeadElements, openGraph)
		_ = pageHypha.RenderTo(meta, data)

		// TODO: check head cats
		// TODO: check opengraph
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
