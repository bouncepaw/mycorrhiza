package web

import (
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

// InitAdmin sets up /admin routes if auth is used. Call it after you have decided if you want to use auth.
func InitAdmin() {
	if user.AuthUsed {
		http.HandleFunc("/admin", HandlerAdmin)
		http.HandleFunc("/admin/shutdown", HandlerAdminShutdown)
		http.HandleFunc("/admin/reindex-users", HandlerAdminReindexUsers)
	}
}

func HandlerAdmin(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if user.CanProceed(rq, "admin") {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(views.BaseHTML("Admin panel", views.AdminPanelHTML(), user.FromRequest(rq))))
	}
}

func HandlerAdminShutdown(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if user.CanProceed(rq, "admin/shutdown") && rq.Method == "POST" {
		log.Fatal("An admin commanded the wiki to shutdown")
	}
}

func HandlerAdminReindexUsers(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if user.CanProceed(rq, "admin") && rq.Method == "POST" {
		user.ReadUsersFromFilesystem()
		http.Redirect(w, rq, "/hypha/"+cfg.UserHypha, http.StatusSeeOther)
	}
}
