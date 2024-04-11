// Package misc provides miscellaneous informative views.
package misc

import (
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/backlinks"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/static"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
)

func InitAssetHandlers(rtr *mux.Router) {
	rtr.HandleFunc("/static/style.css", handlerStyle)
	rtr.HandleFunc("/robots.txt", handlerRobotsTxt)
	rtr.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))
	rtr.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/static/favicon.ico", http.StatusSeeOther)
	})
}

func InitHandlers(rtr *mux.Router) {
	rtr.HandleFunc("/list", handlerList)
	rtr.HandleFunc("/reindex", handlerReindex)
	rtr.HandleFunc("/update-header-links", handlerUpdateHeaderLinks)
	rtr.HandleFunc("/random", handlerRandom)
	rtr.HandleFunc("/about", handlerAbout)
	rtr.HandleFunc("/title-search/", handlerTitleSearch)
	initViews()
}

// handlerList shows a list of all hyphae in the wiki in random order.
func handlerList(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	// TODO: make this more effective, there are too many loops and vars
	var (
		hyphaNames  = make(chan string)
		sortedHypha = hyphae.PathographicSort(hyphaNames)
		entries     []listDatum
	)
	for hypha := range hyphae.YieldExistingHyphae() {
		hyphaNames <- hypha.CanonicalName()
	}
	close(hyphaNames)
	for hyphaName := range sortedHypha {
		switch h := hyphae.ByName(hyphaName).(type) {
		case *hyphae.TextualHypha:
			entries = append(entries, listDatum{h.CanonicalName(), ""})
		case *hyphae.MediaHypha:
			entries = append(entries, listDatum{h.CanonicalName(), filepath.Ext(h.MediaFilePath())[1:]})
		}
	}
	viewList(viewutil.MetaFrom(w, rq), entries)
}

// handlerReindex reindexes all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		var lc = l18n.FromRequest(rq)
		viewutil.HttpErr(viewutil.MetaFrom(w, rq), http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.reindex_no_rights"))
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae.ResetCount()
	log.Println("Reindexing hyphae in", files.HyphaeDir())
	hyphae.Index(files.HyphaeDir())
	backlinks.IndexBacklinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerUpdateHeaderLinks updates header links by reading the configured hypha, if there is any, or resorting to default values.
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		var lc = l18n.FromRequest(rq)
		viewutil.HttpErr(viewutil.MetaFrom(w, rq), http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.header_no_rights"))
		log.Println("Rejected", rq.URL)
		return
	}
	shroom.SetHeaderLinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerRandom redirects to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		randomHyphaName string
		amountOfHyphae  = hyphae.Count()
	)
	if amountOfHyphae == 0 {
		var lc = l18n.FromRequest(rq)
		viewutil.HttpErr(viewutil.MetaFrom(w, rq), http.StatusNotFound, cfg.HomeHypha, lc.Get("ui.random_no_hyphae_tip"))
		return
	}
	i := rand.Intn(amountOfHyphae)
	for h := range hyphae.YieldExistingHyphae() {
		if i == 0 {
			randomHyphaName = h.CanonicalName()
		}
		i--
	}
	http.Redirect(w, rq, "/hypha/"+randomHyphaName, http.StatusSeeOther)
}

// handlerAbout shows a summary of wiki's software.
func handlerAbout(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var (
		lc    = l18n.FromRequest(rq)
		title = lc.Get("ui.about_title", &l18n.Replacements{"name": cfg.WikiName})
	)
	_, err := io.WriteString(w, viewutil.Base(
		viewutil.MetaFrom(w, rq),
		title,
		AboutHTML(lc),
		map[string]string{},
	))
	if err != nil {
		log.Println(err)
	}
}

var stylesheets = []string{"default.css", "custom.css"}

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	var filesData []io.Reader

	var latestModTime time.Time

	// Step 1: Collect files data and determine the latest modification time
	for _, name := range stylesheets {
		file, err := static.FS.Open(name)
		if err != nil {
			continue
		}

		fileStats, err := file.Stat()
		if err != nil {
			continue
		}

		modTime := fileStats.ModTime()
		if modTime.After(latestModTime) {
			latestModTime = modTime
		}

		filesData = append(filesData, file)

		defer file.Close()
	}

	// Step 2: Check the "If-Modified-Since" header in the request
	if ifModifiedSince := rq.Header.Get("If-Modified-Since"); ifModifiedSince != "" {
		if ifModSinceTime, err := http.ParseTime(ifModifiedSince); err == nil && !latestModTime.UTC().After(ifModSinceTime) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Step 3: If content needs to be served, set Last-Modified header and serve the content
	w.Header().Set("Content-Type", mime.TypeByExtension(".css"))
	w.Header().Set("Last-Modified", latestModTime.UTC().Format(http.TimeFormat))

	for _, data := range filesData {
		if _, err := io.Copy(w, data); err != nil {
			log.Println(err)
		}
	}
}

func handlerRobotsTxt(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	file, err := static.FS.Open("robots.txt")
	if err != nil {
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println()
	}
	_ = file.Close()
}

func handlerTitleSearch(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	_ = rq.ParseForm()
	var (
		query       = rq.FormValue("q")
		hyphaName   = util.CanonicalName(query)
		_, nameFree = hyphae.AreFreeNames(hyphaName)
		results     []string
	)
	for hyphaName := range shroom.YieldHyphaNamesContainingString(query) {
		results = append(results, hyphaName)
	}
	w.WriteHeader(http.StatusOK)
	viewTitleSearch(viewutil.MetaFrom(w, rq), query, hyphaName, !nameFree, results)
}
