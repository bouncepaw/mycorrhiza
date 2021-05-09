package web

import (
	"io"
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

// initAdmin sets up /admin routes if auth is used. Call it after you have decided if you want to use auth.
func initAdmin() {
	if user.AuthUsed {
		http.HandleFunc("/admin", handlerAdmin)
		http.HandleFunc("/admin/shutdown", handlerAdminShutdown)
		http.HandleFunc("/admin/reindex-users", handlerAdminReindexUsers)
	}
}

// handlerAdmin provides the admin panel.
func handlerAdmin(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if user.CanProceed(rq, "admin") {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, views.BaseHTML("Admin panel", views.AdminPanelHTML(), user.FromRequest(rq)))
		if err != nil {
			log.Println(err)
		}
	}
}

// handlerAdminShutdown kills the wiki.
func handlerAdminShutdown(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if user.CanProceed(rq, "admin/shutdown") && rq.Method == "POST" {
		log.Fatal("An admin commanded the wiki to shutdown")
	}
}

// handlerAdminReindexUsers reinitialises the user system.
func handlerAdminReindexUsers(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if user.CanProceed(rq, "admin") && rq.Method == "POST" {
		user.ReadUsersFromFilesystem()
		http.Redirect(w, rq, "/hypha/"+cfg.UserHypha, http.StatusSeeOther)
	}
}
