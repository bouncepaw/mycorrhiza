//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=assets
//go:generate qtc -dir=views
//go:generate qtc -dir=tree
package main

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bouncepaw/mycorrhiza/assets"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
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

// This part is present in all html documents.
var base = views.BaseHTML

// Stop the wiki

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	if _, err := os.Stat(cfg.WikiDir + "/static/common.css"); err == nil {
		http.ServeFile(w, rq, cfg.WikiDir+"/static/common.css")
	} else {
		w.Header().Set("Content-Type", "text/css;charset=utf-8")
		w.Write([]byte(assets.DefaultCSS()))
	}
	if bytes, err := ioutil.ReadFile(cfg.WikiDir + "/static/custom.css"); err == nil {
		w.Write(bytes)
	}
}

func handlerToolbar(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
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

func prepareRq(rq *http.Request) {
	log.Println(rq.RequestURI)
	rq.URL.Path = strings.TrimSuffix(rq.URL.Path, "/")
}

func main() {
	parseCliArgs()

	// It is ok if the path is ""
	cfg.ReadConfigFile(cfg.ConfigFilePath)

	if err := files.CalculatePaths(); err != nil {
		log.Fatal(err)
	}

	log.Println("Running MycorrhizaWiki")
	if err := os.Chdir(WikiDir); err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki storage directory is", WikiDir)
	hyphae.Index(WikiDir)
	log.Println("Indexed", hyphae.Count(), "hyphae")

	// Initialize user database
	user.InitUserDatabase()

	history.Start(WikiDir)
	shroom.SetHeaderLinks()

	go handleGemini()

	// See http_admin.go for /admin, /admin/*
	initAdmin()
	// See http_readers.go for /page/, /hypha/, /text/, /binary/, /attachment/
	// See http_mutators.go for /upload-binary/, /upload-text/, /edit/, /delete-ask/, /delete-confirm/, /rename-ask/, /rename-confirm/, /unattach-ask/, /unattach-confirm/
	// See http_auth.go for /login, /login-data, /logout, /logout-confirm
	http.HandleFunc("/user-list/", handlerUserList)
	// See http_history.go for /history/, /recent-changes
	// See http_stuff.go for list, reindex, update-header-links, random, about
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(WikiDir+"/static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, WikiDir+"/static/favicon.ico")
	})
	http.HandleFunc("/static/common.css", handlerStyle)
	http.HandleFunc("/static/toolbar.js", handlerToolbar)
	http.HandleFunc("/static/icon/", handlerIcon)
	http.HandleFunc("/robots.txt", handlerRobotsTxt)
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		addr, _ := url.Parse("/hypha/" + cfg.HomeHypha) // Let's pray it never fails
		rq.URL = addr
		handlerHypha(w, rq)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:"+cfg.HTTPPort, nil))
}
