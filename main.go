//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=assets
//go:generate qtc -dir=views
//go:generate qtc -dir=tree
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/bouncepaw/mycorrhiza/assets"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

// WikiDir is a rooted path to the wiki storage directory.
var WikiDir string

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
				`<main class="main-width"><p>%s. <a href="/page/%s">Go back to the hypha.<a></p></main>`,
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
	util.HTTP200Page(w, base("List of pages", views.HyphaListHTML(), user.FromRequest(rq)))
}

// This part is present in all html documents.
var base = views.BaseHTML

// Reindex all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		HttpErr(w, http.StatusForbidden, util.HomePage, "Not enough rights", "You must be an admin to reindex hyphae.")
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae.ResetCount()
	log.Println("Wiki storage directory is", WikiDir)
	log.Println("Start indexing hyphae...")
	hyphae.Index(WikiDir)
	log.Println("Indexed", hyphae.Count(), "hyphae")
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// Stop the wiki

// Update header links by reading the configured hypha, if there is any, or resorting to default values.
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		HttpErr(w, http.StatusForbidden, util.HomePage, "Not enough rights", "You must be a moderator to update header links.")
		log.Println("Rejected", rq.URL)
		return
	}
	shroom.SetHeaderLinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// Redirect to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		randomHyphaName string
		amountOfHyphae  int = hyphae.Count()
	)
	if amountOfHyphae == 0 {
		HttpErr(w, http.StatusNotFound, util.HomePage, "There are no hyphae",
			"It is not possible to display a random hypha because the wiki does not contain any hyphae")
		return
	}
	i := rand.Intn(amountOfHyphae)
	for h := range hyphae.YieldExistingHyphae() {
		if i == 0 {
			randomHyphaName = h.Name
		}
		i--
	}
	http.Redirect(w, rq, "/hypha/"+randomHyphaName, http.StatusSeeOther)
}

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if _, err := os.Stat(util.WikiDir + "/static/common.css"); err == nil {
		http.ServeFile(w, rq, util.WikiDir+"/static/common.css")
	} else {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		w.Write([]byte(assets.DefaultCSS()))
	}
	if bytes, err := ioutil.ReadFile(util.WikiDir + "/static/custom.css"); err == nil {
		w.Write(bytes)
	}
}

func handlerToolbar(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
	w.Write([]byte(assets.ToolbarJS()))
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
		w.Write([]byte(assets.IconGemini()))
	case "mailto":
		w.Write([]byte(assets.IconMailto()))
	case "gopher":
		w.Write([]byte(assets.IconGopher()))
	case "feed":
		w.Write([]byte(assets.IconFeed()))
	default:
		w.Write([]byte(assets.IconHTTP()))
	}
}

func handlerAbout(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base("About "+util.SiteName, views.AboutHTML(), user.FromRequest(rq))))
}

func handlerUserList(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base("User list", views.UserListHTML(), user.FromRequest(rq))))
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
	parseCliArgs()
	log.Println("Running MycorrhizaWiki")
	if err := os.Chdir(WikiDir); err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki storage directory is", WikiDir)
	hyphae.Index(WikiDir)
	log.Println("Indexed", hyphae.Count(), "hyphae")

	if user.AuthUsed && (util.FixedCredentialsPath != "" || util.RegistrationCredentialsPath != "") {
		user.ReadUsersFromFilesystem()
	}
	history.Start(WikiDir)
	shroom.SetHeaderLinks()

	go handleGemini()

	// See http_admin.go for /admin, /admin/*
	initAdmin()
	// See http_readers.go for /page/, /hypha/, /text/, /binary/, /attachment/
	// See http_mutators.go for /upload-binary/, /upload-text/, /edit/, /delete-ask/, /delete-confirm/, /rename-ask/, /rename-confirm/, /unattach-ask/, /unattach-confirm/
	// See http_auth.go for /login, /login-data, /logout, /logout-confirm
	// See http_history.go for /history/, /recent-changes
	http.HandleFunc("/list", handlerList)
	http.HandleFunc("/reindex", handlerReindex)
	http.HandleFunc("/update-header-links", handlerUpdateHeaderLinks)
	http.HandleFunc("/random", handlerRandom)
	http.HandleFunc("/about", handlerAbout)
	http.HandleFunc("/user-list", handlerUserList)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(WikiDir+"/static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, WikiDir+"/static/favicon.ico")
	})
	http.HandleFunc("/static/common.css", handlerStyle)
	http.HandleFunc("/static/toolbar.js", handlerToolbar)
	http.HandleFunc("/static/icon/", handlerIcon)
	http.HandleFunc("/robots.txt", handlerRobotsTxt)
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/hypha/"+util.HomePage, http.StatusSeeOther)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:"+util.ServerPort, nil))
}
