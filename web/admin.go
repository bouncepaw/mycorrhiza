package web

import (
	"fmt"
	"mime"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

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

		http.HandleFunc("/admin/users/", handlerAdminUsers)
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
		redirectTo := rq.Referer()
		if redirectTo == "" {
			redirectTo = "/hypha/"+cfg.UserHypha
		}
		http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
	}
}

func handlerAdminUsers(w http.ResponseWriter, r *http.Request) {
	util.PrepareRq(r)
	if user.CanProceed(r, "admin") {
		path := strings.TrimPrefix(r.URL.Path, "/admin/users")
		parts := strings.Split(path, "/")[1:]

		// Users dashboard
		if len(parts) == 0 {
			// Get a sorted list of users
			var userList []*user.User
			for u := range user.YieldUsers() {
				userList = append(userList, u)
			}

			sort.Slice(userList, func(i, j int) bool {
				less := userList[i].RegisteredAt.Before(userList[j].RegisteredAt)
				return less
			})

			html := views.AdminUsersPanelHTML(userList)
			html = views.BaseHTML("Manage users", html, user.FromRequest(r))

			w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
			if _, err := io.WriteString(w, html); err != nil {
				log.Println(err)
			}
			return
		}

		// User edit page
		if len(parts) == 2 && parts[1] == "edit" {
			u := user.UserByName(parts[0])

			if u != nil && u.Name != "anon" {
				if r.Method == http.MethodGet {
					html := views.AdminUsersUserHTML(u)
					html = views.BaseHTML(fmt.Sprintf("User %s", u.Name), html, user.FromRequest(r))

					w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
					if _, err := io.WriteString(w, html); err != nil {
						log.Println(err)
					}
					return
				} else if r.Method == http.MethodPost {
					oldGroup := u.Group
					newGroup := r.PostFormValue("group")
					if user.ValidGroup(newGroup) {
						u.Group = newGroup
						if err := user.SaveUserDatabase(); err != nil {
							u.Group = oldGroup
							log.Println(err)
							w.WriteHeader(http.StatusInternalServerError)
							io.WriteString(w, err.Error())
						} else {
							http.Redirect(w, r, "/admin/users/", http.StatusSeeOther)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						io.WriteString(w, "invalid group")
					}
					return
				}
			}
		}

		util.HTTP404Page(w, "404 page not found")
	}
}
