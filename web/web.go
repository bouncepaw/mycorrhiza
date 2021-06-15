// Package web contains web handlers and initialization stuff.
//
// It exports just one function: Init. Call it if you want to have web capabilities.
package web

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/static"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

var stylesheets = []string{"default.css", "custom.css"}

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
				`<main class="main-width"><p>%s. <a href="/hypha/%s">Go back to the hypha.<a></p></main>`,
				errMsg,
				name,
			),
			user.EmptyUser(),
		),
	)
}

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)

	w.Header().Set("Content-Type", mime.TypeByExtension("css"))
	for _, name := range stylesheets {
		file, err := static.FS.Open(name)
		if err != nil {
			continue
		}
		io.Copy(w, file)
		file.Close()
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
Allow: /hypha/
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

	// Miscellaneous
	http.HandleFunc("/user-list/", handlerUserList)
	http.HandleFunc("/robots.txt", handlerRobotsTxt)

	// Static assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))
	http.HandleFunc("/static/style.css", handlerStyle)

	// Index page
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		addr, _ := url.Parse("/hypha/" + cfg.HomeHypha) // Let's pray it never fails
		rq.URL = addr
		handlerHypha(w, rq)
	})
}
