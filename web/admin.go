package web

import (
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
)

func viewPanel(meta viewdata.Meta) {
	if err := pageAdminPanel.RenderDataTo(meta, &newtmpl.BaseData{}); err != nil {
		slog.Info("Failed to execute page; proceeding anyway", "page", "admin-panel", "err", err)
	}
}

type listData struct {
	newtmpl.BaseData
	UserHypha string
	Users     []*usersvc.User
}

func viewList(meta viewdata.Meta, users []*usersvc.User) {
	if err := pageAdminUserList.RenderDataTo(meta, &listData{
		UserHypha: cfg.UserHypha,
		Users:     users,
	}); err != nil {
		slog.Info("Failed to execute page; proceeding anyway", "page", "admin-user-list", "err", err)
	}
}

type newUserData struct {
	newtmpl.BaseData
	Form util.FormData
}

func viewNewUser(meta viewdata.Meta, form util.FormData) {
	if err := pageAdminNewUser.RenderDataTo(meta, &newUserData{
		Form: form,
	}); err != nil {
		slog.Info("Failed to execute page; proceeding anyway", "page", "admin-new-user", "err", err)
	}
}

type editDeleteUserData struct {
	newtmpl.BaseData
	Form util.FormData
	U    *usersvc.User
}

func viewEditUser(meta viewdata.Meta, form util.FormData, u *usersvc.User) {
	if err := pageAdminEditUser.RenderDataTo(meta, &editDeleteUserData{
		Form: form,
		U:    u,
	}); err != nil {
		slog.Info("Failed to execute page; proceeding anyway", "page", "admin-edit-user", "err", err)
	}
}

func viewDeleteUser(meta viewdata.Meta, form util.FormData, u *usersvc.User) {
	if err := pageAdminDeleteUser.RenderDataTo(meta, &editDeleteUserData{
		Form: form,
		U:    u,
	}); err != nil {
		slog.Info("Failed to execute page; proceeding anyway", "page", "admin-delete-user", "err", err)
	}
}

// handlerAdmin provides the admin panel.
func handlerAdmin(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	viewPanel(viewdata.MetaFrom(w, rq))
}

// handlerAdminShutdown kills the wiki.
func handlerAdminShutdown(w http.ResponseWriter, rq *http.Request) {
	if usersvc.CanProceed(rq, userdomain.ActionAdminShutdown) {
		slog.Info("An admin commanded the wiki to shutdown")
		os.Exit(0)
	}
}

// handlerAdminReindexUsers reinitialises the user system.
func handlerAdminReindexUsers(w http.ResponseWriter, rq *http.Request) {
	usersvc.ReadUsersFromFilesystem()
	redirectTo := rq.Referer()
	if redirectTo == "" {
		redirectTo = "/hypha/" + cfg.UserHypha
	}
	http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
}

func handlerAdminUsers(w http.ResponseWriter, rq *http.Request) {
	// Get a sorted list of users
	var users []*usersvc.User
	for u := range usersvc.YieldUsers() {
		users = append(users, u)
	}

	sort.Slice(users, func(i, j int) bool {
		less := users[i].RegisteredAt.Before(users[j].RegisteredAt)
		return less
	})
	viewList(viewdata.MetaFrom(w, rq), users)
}

func handlerAdminUserEdit(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := usersvc.ByName(vars["username"])
	if u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.FormDataFromRequest(rq, []string{"group"})

	if rq.Method == http.MethodPost {
		oldGroup := u.Group
		newGroup := userdomain.GroupKind(f.Get("group"))

		if newGroup.IsValid() {
			u.Group = newGroup
			if err := usersvc.SaveUserDatabase(); err != nil {
				u.Group = oldGroup
				slog.Info("Failed to save user database", "err", err)
				f = f.WithError(err)
			} else {
				http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
				return
			}
		} else {
			f = f.WithError(fmt.Errorf("invalid group ‘%s’", newGroup))
		}
	}

	f.Put("group", u.Group.String())

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))

	viewEditUser(viewdata.MetaFrom(w, rq), f, u)
}

func handlerAdminUserChangePassword(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := usersvc.ByName(vars["username"])
	if u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.FormDataFromRequest(rq, []string{"password", "password_confirm"})

	password := f.Get("password")
	passwordConfirm := f.Get("password_confirm")
	// server side validation
	if password == "" {
		err := fmt.Errorf("passwords should not be empty")
		f = f.WithError(err)
	}
	if password == passwordConfirm {
		previousPassword := u.Password // for rollback
		if err := u.ChangePassword(password); err != nil {
			f = f.WithError(err)
		} else {
			if err := usersvc.SaveUserDatabase(); err != nil {
				u.Password = previousPassword
				f = f.WithError(err)
			} else {
				http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
				return
			}
		}
	} else {
		err := fmt.Errorf("passwords do not match")
		f = f.WithError(err)
	}

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))

	viewEditUser(viewdata.MetaFrom(w, rq), f, u)
}

func handlerAdminUserDelete(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := usersvc.ByName(vars["username"])
	if u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.NewFormData()

	if rq.Method == http.MethodPost {
		f = f.WithError(usersvc.DeleteUser(u.Name))
		if !f.HasError() {
			http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
		} else {
			slog.Info("Failed to delete user", "err", f.Error())
		}
	}

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	viewDeleteUser(viewdata.MetaFrom(w, rq), f, u)
}

func handlerAdminUserNew(w http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodGet {
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		viewNewUser(viewdata.MetaFrom(w, rq), util.NewFormData())
	} else if rq.Method == http.MethodPost {
		// Create a user
		f := util.FormDataFromRequest(rq, []string{"name", "password", "group"})

		err := usersvc.Register(
			f.Get("name"),
			f.Get("password"),
			userdomain.GroupKind(f.Get("group")),
			userdomain.SourceLocal,
			true,
		)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
			viewNewUser(viewdata.MetaFrom(w, rq), f.WithError(err))
		} else {
			http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
		}
	}
}
