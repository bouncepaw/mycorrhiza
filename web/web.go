// Package web contains web handlers and initialization stuff.
//
// It exports just one function: Init. Call it if you want to have web capabilities.
package web

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/assets"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/views"
)

// httpErr is used by many handlers to signal errors in a compact way.
func httpErr(w http.ResponseWriter, status int, name, title, errMsg string) {
	log.Println(errMsg, "for", name)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(
		w,
		views.BaseHTML(
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

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
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
	util.PrepareRq(rq)
	w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
	w.Write([]byte(assets.ToolbarJS()))
}

// handlerIcon serves the requested icon. All icons are distributed as part of the Mycorrhiza binary.
//
// See assets/assets/icon/ for icons themselves, see assets/assets.qtpl for their sources.
func handlerIcon(w http.ResponseWriter, rq *http.Request) {
	iconName := strings.TrimPrefix(rq.URL.Path, "/assets/icon/")
	if iconName == "https" {
		iconName = "http"
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	icon := func() string {
		switch iconName {
		case "gemini":
			return assets.IconGemini()
		case "mailto":
			return assets.IconMailto()
		case "gopher":
			return assets.IconGopher()
		case "feed":
			return assets.IconFeed()
		default:
			return assets.IconHTTP()
		}
	}()
	_, err := io.WriteString(w, icon)
	if err != nil {
		log.Println(err)
	}

}

func handlerUserList(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(views.BaseHTML("User list", views.UserListHTML(), user.FromRequest(rq))))
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

func Init() {
	initAdmin()
	initReaders()
	initMutators()
	initAuth()
	initHistory()
	initStuff()

	http.HandleFunc("/user-list/", handlerUserList)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.WikiDir+"/static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, cfg.WikiDir+"/static/favicon.ico")
	})
	http.HandleFunc("/static/common.css", handlerStyle)
	http.HandleFunc("/static/toolbar.js", handlerToolbar)
	http.HandleFunc("/assets/icon/", handlerIcon)
	http.HandleFunc("/robots.txt", handlerRobotsTxt)
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		addr, _ := url.Parse("/hypha/" + cfg.HomeHypha) // Let's pray it never fails
		rq.URL = addr
		handlerHypha(w, rq)
	})
}
