//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/util"
)

// WikiDir is a rooted path to the wiki storage directory.
var WikiDir string

// HyphaPattern is a pattern which all hyphae must match.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"\'&%]+`)

// HyphaStorage is a mapping between canonical hypha names and their meta information.
var HyphaStorage = make(map[string]*HyphaData)

// IterateHyphaNamesWith is a closure to be passed to subpackages to let them iterate all hypha names read-only.
func IterateHyphaNamesWith(f func(string)) {
	for hyphaName := range HyphaStorage {
		f(hyphaName)
	}
}

// HttpErr is used by many handlers to signal errors in a compact way.
func HttpErr(w http.ResponseWriter, status int, name, title, errMsg string) {
	log.Println(errMsg, "for", name)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, base(title, fmt.Sprintf(
		`<main><p>%s. <a href="/page/%s">Go back to the hypha.<a></p></main>`,
		errMsg, name)))
}

// Show all hyphae
func handlerList(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		tbody     string
		pageCount = len(HyphaStorage)
	)
	for hyphaName, data := range HyphaStorage {
		tbody += templates.HyphaListRowHTML(hyphaName, ExtensionToMime(filepath.Ext(data.binaryPath)), data.binaryPath != "")
	}
	util.HTTP200Page(w, base("List of pages", templates.HyphaListHTML(tbody, pageCount)))
}

// This part is present in all html documents.
var base = templates.BaseHTML

// Reindex all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	HyphaStorage = make(map[string]*HyphaData)
	log.Println("Wiki storage directory is", WikiDir)
	log.Println("Start indexing hyphae...")
	Index(WikiDir)
	log.Println("Indexed", len(HyphaStorage), "hyphae")
}

// Redirect to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var randomHyphaName string
	i := rand.Intn(len(HyphaStorage))
	for hyphaName := range HyphaStorage {
		if i == 0 {
			randomHyphaName = hyphaName
			break
		}
		i--
	}
	http.Redirect(w, rq, "/page/"+randomHyphaName, http.StatusSeeOther)
}

// Recent changes
func handlerRecentChanges(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		noPrefix = strings.TrimPrefix(rq.URL.String(), "/recent-changes/")
		n, err   = strconv.Atoi(noPrefix)
	)
	if err == nil {
		util.HTTP200Page(w, base(strconv.Itoa(n)+" recent changes", history.RecentChanges(n)))
	} else {
		http.Redirect(w, rq, "/recent-changes/20", http.StatusSeeOther)
	}
}

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if _, err := os.Stat(WikiDir + "/static/common.css"); err == nil {
		http.ServeFile(w, rq, WikiDir+"/static/common.css")
	} else {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		w.Write([]byte(templates.DefaultCSS()))
	}
}

func main() {
	log.Println("Running MycorrhizaWiki β")
	parseCliArgs()
	if err := os.Chdir(WikiDir); err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki storage directory is", WikiDir)
	log.Println("Start indexing hyphae...")
	Index(WikiDir)
	log.Println("Indexed", len(HyphaStorage), "hyphae")

	history.Start(WikiDir)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(WikiDir+"/static"))))
	// See http_readers.go for /page/, /text/, /binary/, /history/.
	// See http_mutators.go for /upload-binary/, /upload-text/, /edit/, /delete-ask/, /delete-confirm/, /rename-ask/, /rename-confirm/.
	http.HandleFunc("/list", handlerList)
	http.HandleFunc("/reindex", handlerReindex)
	http.HandleFunc("/random", handlerRandom)
	http.HandleFunc("/recent-changes/", handlerRecentChanges)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, WikiDir+"/static/favicon.ico")
	})
	http.HandleFunc("/static/common.css", handlerStyle)
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/page/"+util.HomePage, http.StatusSeeOther)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:"+util.ServerPort, nil))
}
