package web

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

// initAdmin sets up /admin routes if auth is used. Call it after you have decided if you want to use auth.
func initAdmin(r *mux.Router) {
	r.HandleFunc("/shutdown", handlerAdminShutdown).Methods(http.MethodPost)
	r.HandleFunc("/reindex-users", handlerAdminReindexUsers).Methods(http.MethodPost)

	r.HandleFunc("/user/new", handlerAdminUserNew).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/users/{username}/edit", handlerAdminUserEdit).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/users/{username}/delete", handlerAdminUserDelete).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/users", handlerAdminUsers)

	r.HandleFunc("/", handlerAdmin)
}

// handlerAdmin provides the admin panel.
func handlerAdmin(w http.ResponseWriter, rq *http.Request) {
	var lc = l18n.FromRequest(rq)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, views.BaseHTML(lc.Get("admin.panel_title"), views.AdminPanelHTML(lc), lc, user.FromRequest(rq)))
	if err != nil {
		log.Println("an error occurred in handlerAdmin function:", err)
	}
}

// handlerAdminShutdown kills the wiki.
func handlerAdminShutdown(w http.ResponseWriter, rq *http.Request) {
	if user.CanProceed(rq, "admin/shutdown") {
		log.Println("An admin commanded the wiki to shutdown")
		os.Exit(0)
	}
}

// handlerAdminReindexUsers reinitialises the user system.
func handlerAdminReindexUsers(w http.ResponseWriter, rq *http.Request) {
	user.ReadUsersFromFilesystem()
	redirectTo := rq.Referer()
	if redirectTo == "" {
		redirectTo = "/hypha/" + cfg.UserHypha
	}
	http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
}

func handlerAdminUsers(w http.ResponseWriter, rq *http.Request) {
	// Get a sorted list of users
	var userList []*user.User
	for u := range user.YieldUsers() {
		userList = append(userList, u)
	}

	sort.Slice(userList, func(i, j int) bool {
		less := userList[i].RegisteredAt.Before(userList[j].RegisteredAt)
		return less
	})

	var lc = l18n.FromRequest(rq)
	html := views.AdminUsersPanelHTML(userList, lc)
	html = views.BaseHTML(lc.Get("admin.users_title"), html, lc, user.FromRequest(rq))

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	_, err := io.WriteString(w, html)
	if err != nil {
		log.Println("an error occurred in handlerAdminUsers function:", err)
	}
}

func handlerAdminUserEdit(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := user.ByName(vars["username"])
	if u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.FormDataFromRequest(rq, []string{"group"})

	if rq.Method == http.MethodPost {
		oldGroup := u.Group
		newGroup := f.Get("group")

		if user.ValidGroup(newGroup) {
			u.Group = newGroup
			if err := user.SaveUserDatabase(); err != nil {
				u.Group = oldGroup
				log.Println(err)
				f = f.WithError(err)
			} else {
				http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
				return
			}
		} else {
			f = f.WithError(fmt.Errorf("invalid group ‘%s’", newGroup))
		}
	}

	f.Put("group", u.Group)

	var lc = l18n.FromRequest(rq)
	html := views.AdminUserEditHTML(u, f, lc)
	html = views.BaseHTML(fmt.Sprintf(lc.Get("admin.user_title"), u.Name), html, lc, user.FromRequest(rq))

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	_, err := io.WriteString(w, html)
	if err != nil {
		log.Println("an error occurred in handlerAdminUserEdit function:", err)
	}
}

func handlerAdminUserDelete(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := user.ByName(vars["username"])
	if u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.NewFormData()

	if rq.Method == http.MethodPost {
		f = f.WithError(user.DeleteUser(u.Name))
		if !f.HasError() {
			http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
		} else {
			log.Println(f.Error())
		}
	}

	var lc = l18n.FromRequest(rq)
	html := views.AdminUserDeleteHTML(u, util.NewFormData(), lc)
	html = views.BaseHTML(fmt.Sprintf(lc.Get("admin.user_title"), u.Name), html, l18n.FromRequest(rq), user.FromRequest(rq))

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	_, err := io.WriteString(w, html)
	if err != nil {
		log.Println("an error occurred in handlerAdminUSerDelete function:", err)
	}
}

func handlerAdminUserNew(w http.ResponseWriter, rq *http.Request) {
	var lc = l18n.FromRequest(rq)
	if rq.Method == http.MethodGet {
		// New user form
		html := views.AdminUserNewHTML(util.NewFormData(), lc)
		html = views.BaseHTML(lc.Get("admin.newuser_title"), html, lc, user.FromRequest(rq))

		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		_, err := io.WriteString(w, html)
		if err != nil {
			log.Println("an error occurred in handlerAdminUserNew function, in get method:", err)
		}
	} else if rq.Method == http.MethodPost {
		// Create a user
		f := util.FormDataFromRequest(rq, []string{"name", "password", "group"})

		err := user.Register(f.Get("name"), f.Get("password"), f.Get("group"), "local", true)

		if err != nil {
			html := views.AdminUserNewHTML(f.WithError(err), lc)
			html = views.BaseHTML(lc.Get("admin.newuser_title"), html, lc, user.FromRequest(rq))

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
			_, err := io.WriteString(w, html)
			if err != nil {
				log.Println("an error occurred in handlerAdminUSerNew function, in post method:", err)
			}
		} else {
			http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
		}
	}
}
