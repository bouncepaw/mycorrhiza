package admin

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"sort"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

// handlerAdmin provides the admin panel.
func handlerAdmin(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	viewPanel(viewutil.MetaFrom(w, rq))
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
	var users []*user.User
	for u := range user.YieldUsers() {
		users = append(users, u)
	}

	sort.Slice(users, func(i, j int) bool {
		less := users[i].RegisteredAt.Before(users[j].RegisteredAt)
		return less
	})
	viewList(viewutil.MetaFrom(w, rq), users)
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
	html := views.AdminUserEdit(u, f, lc)
	html = views.Base(viewutil.MetaFrom(w, rq), fmt.Sprintf(lc.Get("admin.user_title"), u.Name), html)

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	io.WriteString(w, html)
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
	html := views.AdminUserDelete(u, util.NewFormData(), lc)
	html = views.Base(viewutil.MetaFrom(w, rq), fmt.Sprintf(lc.Get("admin.user_title"), u.Name), html)

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	io.WriteString(w, html)
}

func handlerAdminUserNew(w http.ResponseWriter, rq *http.Request) {
	var lc = l18n.FromRequest(rq)
	if rq.Method == http.MethodGet {
		// New user form
		html := views.AdminUserNew(util.NewFormData(), lc)
		html = views.Base(viewutil.MetaFrom(w, rq), lc.Get("admin.newuser_title"), html)

		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		io.WriteString(w, html)
	} else if rq.Method == http.MethodPost {
		// Create a user
		f := util.FormDataFromRequest(rq, []string{"name", "password", "group"})

		err := user.Register(f.Get("name"), f.Get("password"), f.Get("group"), "local", true)

		if err != nil {
			html := views.AdminUserNew(f.WithError(err), lc)
			html = views.Base(viewutil.MetaFrom(w, rq), lc.Get("admin.newuser_title"), html)

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
			io.WriteString(w, html)
		} else {
			http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
		}
	}
}
