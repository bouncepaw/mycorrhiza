//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/util"
)

// WikiDir is a rooted path to the wiki storage directory.
var WikiDir string

// HyphaPattern is a pattern which all hyphae must match. Not used currently.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"\'&%]+`)

// HyphaStorage is a mapping between canonical hypha names and their meta information.
var HyphaStorage = make(map[string]*HyphaData)

// IterateHyphaNamesWith is a closure to be passed to subpackages to let them iterate all hypha names read-only.
func IterateHyphaNamesWith(f func(string)) {
	for hyphaName, _ := range HyphaStorage {
		f(hyphaName)
	}
}

// HttpErr is used by many handlers to signal errors in a compact way.
func HttpErr(w http.ResponseWriter, status int, name, title, errMsg string) {
	log.Println(errMsg, "for", name)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, base(title, fmt.Sprintf(
		`<p>%s. <a href="/page/%s">Go back to the hypha.<a></p>`,
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
		tbody += templates.HyphaListRowHTML(hyphaName, data.binaryType.Mime(), data.binaryPath != "")
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

func main() {
	log.Println("Running MycorrhizaWiki β")

	var err error
	WikiDir, err = filepath.Abs(os.Args[1])
	util.WikiDir = WikiDir
	if err != nil {
		log.Fatal(err)
	}
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
	// See http_mutators.go for /upload-binary/, /upload-text/, /edit/.
	http.HandleFunc("/list", handlerList)
	http.HandleFunc("/reindex", handlerReindex)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, WikiDir+"/static/favicon.ico")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/page/home", http.StatusSeeOther)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:1737", nil))
}
