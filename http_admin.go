package main

import (
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

// This is not init(), because user.AuthUsed is not set at init-stage.
func initAdmin() {
	if user.AuthUsed {
		http.HandleFunc("/admin", handlerAdmin)
		http.HandleFunc("/admin/shutdown", handlerAdminShutdown)
		http.HandleFunc("/admin/reindex-users", handlerAdminReindexUsers)
	}
}

func handlerAdmin(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if user.CanProceed(rq, "admin") {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(base("Admin panel", views.AdminPanelHTML(), user.FromRequest(rq))))
	}
}

func handlerAdminShutdown(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if user.CanProceed(rq, "admin/shutdown") && rq.Method == "POST" {
		log.Fatal("An admin commanded the wiki to shutdown")
	}
}

func handlerAdminReindexUsers(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if user.CanProceed(rq, "admin") && rq.Method == "POST" {
		user.ReadUsersFromFilesystem()
		http.Redirect(w, rq, "/hypha/"+util.UserHypha, http.StatusSeeOther)
	}
}
