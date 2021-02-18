package main

import (
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/user"
)

// This is not init(), because user.AuthUsed is not set at init-stage.
func initAdmin() {
	if user.AuthUsed {
		http.HandleFunc("/admin", handlerAdmin)
		http.HandleFunc("/admin/shutdown", handlerAdminShutdown)
	}
}

func handlerAdmin(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if user.CanProceed(rq, "admin") {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(base("Admin panel", templates.AdminPanelHTML(), user.FromRequest(rq))))
	}
}

func handlerAdminShutdown(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	if user.CanProceed(rq, "admin/shutdown") && rq.Method == "POST" {
		log.Fatal("An admin commanded the wiki to shutdown")
	}
}
