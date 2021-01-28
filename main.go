//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

// WikiDir is a rooted path to the wiki storage directory.
var WikiDir string

// HyphaPattern is a pattern which all hyphae must match.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"\'&%{}]+`)

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
	fmt.Fprint(
		w,
		base(
			title,
			fmt.Sprintf(
				`<main><p>%s. <a href="/page/%s">Go back to the hypha.<a></p></main>`,
				errMsg,
				name,
			),
			user.EmptyUser(),
		),
	)
}

// Show all hyphae
func handlerList(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		tbody     string
		pageCount = hyphae.Count()
		u         = user.FromRequest(rq)
	)
	for hyphaName, data := range HyphaStorage {
		tbody += templates.HyphaListRowHTML(hyphaName, mimetype.FromExtension(filepath.Ext(data.binaryPath)), data.binaryPath != "")
	}
	util.HTTP200Page(w, base("List of pages", templates.HyphaListHTML(tbody, pageCount), u))
}

// This part is present in all html documents.
var base = templates.BaseHTML

// Reindex all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		HttpErr(w, http.StatusForbidden, util.HomePage, "Not enough rights", "You must be an admin to reindex hyphae.")
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae.ResetCount()
	HyphaStorage = make(map[string]*HyphaData)
	log.Println("Wiki storage directory is", WikiDir)
	log.Println("Start indexing hyphae...")
	Index(WikiDir)
	log.Println("Indexed", hyphae.Count(), "hyphae")
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// Update header links by reading the configured hypha, if there is any, or resorting to default values.
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		HttpErr(w, http.StatusForbidden, util.HomePage, "Not enough rights", "You must be a moderator to update header links.")
		log.Println("Rejected", rq.URL)
		return
	}
	setHeaderLinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// Redirect to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var randomHyphaName string
	i := rand.Intn(hyphae.Count())
	for hyphaName := range HyphaStorage {
		if i == 0 {
			randomHyphaName = hyphaName
			break
		}
		i--
	}
	http.Redirect(w, rq, "/page/"+randomHyphaName, http.StatusSeeOther)
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

func handlerIcon(w http.ResponseWriter, rq *http.Request) {
	iconName := strings.TrimPrefix(rq.URL.Path, "/static/icon/")
	if iconName == "https" {
		iconName = "http"
	}
	files, err := ioutil.ReadDir(WikiDir + "/static/icon")
	if err == nil {
		for _, f := range files {
			if strings.HasPrefix(f.Name(), iconName+"-protocol-icon") {
				http.ServeFile(w, rq, WikiDir+"/static/icon/"+f.Name())
				return
			}
		}
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	switch iconName {
	case "gemini":
		w.Write([]byte(templates.IconGemini()))
	case "mailto":
		w.Write([]byte(templates.IconMailto()))
	case "gopher":
		w.Write([]byte(templates.IconGopher()))
	default:
		w.Write([]byte(templates.IconHTTP()))
	}
}

func handlerAbout(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base("About "+util.SiteName, templates.AboutHTML(), user.FromRequest(rq))))
}

func handlerRobotsTxt(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		`User-agent: *
Allow: /page/
Allow: /recent-changes
Disallow: /
Crawl-delay: 5`))
}

func main() {
	log.Println("Running MycorrhizaWiki β")
	parseCliArgs()
	if err := os.Chdir(WikiDir); err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki storage directory is", WikiDir)
	Index(WikiDir)
	log.Println("Indexed", hyphae.Count(), "hyphae")

	history.Start(WikiDir)
	setHeaderLinks()

	// See http_readers.go for /page/, /text/, /binary/
	// See http_mutators.go for /upload-binary/, /upload-text/, /edit/, /delete-ask/, /delete-confirm/, /rename-ask/, /rename-confirm/, /unattach-ask/, /unattach-confirm/
	// See http_auth.go for /login, /login-data, /logout, /logout-confirm
	// See http_history.go for /history/, /recent-changes
	http.HandleFunc("/list", handlerList)
	http.HandleFunc("/reindex", handlerReindex)
	http.HandleFunc("/update-header-links", handlerUpdateHeaderLinks)
	http.HandleFunc("/random", handlerRandom)
	http.HandleFunc("/about", handlerAbout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(WikiDir+"/static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, WikiDir+"/static/favicon.ico")
	})
	http.HandleFunc("/static/common.css", handlerStyle)
	http.HandleFunc("/static/icon/", handlerIcon)
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/page/"+util.HomePage, http.StatusSeeOther)
	})
	http.HandleFunc("/robots.txt", handlerRobotsTxt)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+util.ServerPort, nil))
}
