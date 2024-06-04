package web

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	user2 "github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/util"
	viewutil2 "github.com/bouncepaw/mycorrhiza/web/viewutil"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"os"
	"sort"
)

const adminTranslationRu = `
{{define "panel title"}}Панель админстратора{{end}}
{{define "panel safe section title"}}Безопасная секция{{end}}
{{define "panel link about"}}Об этой вики{{end}}
{{define "panel update header"}}Обновить ссылки в верхней панели{{end}}
{{define "panel link user list"}}Список пользователей{{end}}
{{define "panel users"}}Управление пользователями{{end}}
{{define "panel unsafe section title"}}Опасная секция{{end}}
{{define "panel shutdown"}}Выключить вики{{end}}
{{define "panel reindex hyphae"}}Переиндексировать гифы{{end}}
{{define "panel interwiki"}}Интервики{{end}}

{{define "manage users"}}Управление пользователями{{end}}
{{define "create user"}}Создать пользователя{{end}}
{{define "reindex users"}}Переиндексировать пользователей{{end}}
{{define "name"}}Имя{{end}}
{{define "group"}}Группа{{end}}
{{define "registered at"}}Зарегистрирован{{end}}
{{define "actions"}}Действия{{end}}
{{define "edit"}}Изменить{{end}}

{{define "new user"}}Новый пользователь{{end}}
{{define "password"}}Пароль{{end}}
{{define "confirm password"}}Подтвердить пароль{{end}}
{{define "change password"}}Изменить пароль{{end}}
{{define "non local password change"}}Поменять пароль можно только у локальных пользователей.{{end}}
{{define "create"}}Создать{{end}}

{{define "change group"}}Изменить группу{{end}}
{{define "user x"}}Пользователь {{.}}{{end}}
{{define "update"}}Обновить{{end}}
{{define "delete user"}}Удалить пользователя{{end}}
{{define "delete user tip"}}Удаляет пользователя из базы данных. Правки пользователя будут сохранены. Имя пользователя освободится для повторной регистрации.{{end}}

{{define "delete user?"}}Удалить пользователя {{.}}?{{end}}
{{define "delete user warning"}}Вы уверены, что хотите удалить этого пользователя из базы данных? Это действие нельзя отменить.{{end}}
`

func viewPanel(meta viewutil2.Meta) {
	viewutil2.ExecutePage(meta, panelChain, &viewutil2.BaseData{})
}

type listData struct {
	*viewutil2.BaseData
	UserHypha string
	Users     []*user2.User
}

func viewList(meta viewutil2.Meta, users []*user2.User) {
	viewutil2.ExecutePage(meta, listChain, listData{
		BaseData:  &viewutil2.BaseData{},
		UserHypha: cfg.UserHypha,
		Users:     users,
	})
}

type newUserData struct {
	*viewutil2.BaseData
	Form util.FormData
}

func viewNewUser(meta viewutil2.Meta, form util.FormData) {
	viewutil2.ExecutePage(meta, newUserChain, newUserData{
		BaseData: &viewutil2.BaseData{},
		Form:     form,
	})
}

type editDeleteUserData struct {
	*viewutil2.BaseData
	Form util.FormData
	U    *user2.User
}

func viewEditUser(meta viewutil2.Meta, form util.FormData, u *user2.User) {
	viewutil2.ExecutePage(meta, editUserChain, editDeleteUserData{
		BaseData: &viewutil2.BaseData{},
		Form:     form,
		U:        u,
	})
}

func viewDeleteUser(meta viewutil2.Meta, form util.FormData, u *user2.User) {
	viewutil2.ExecutePage(meta, deleteUserChain, editDeleteUserData{
		BaseData: &viewutil2.BaseData{},
		Form:     form,
		U:        u,
	})
}

// handlerAdmin provides the admin panel.
func handlerAdmin(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	viewPanel(viewutil2.MetaFrom(w, rq))
}

// handlerAdminShutdown kills the wiki.
func handlerAdminShutdown(w http.ResponseWriter, rq *http.Request) {
	if user2.CanProceed(rq, "admin/shutdown") {
		log.Println("An admin commanded the wiki to shutdown")
		os.Exit(0)
	}
}

// handlerAdminReindexUsers reinitialises the user system.
func handlerAdminReindexUsers(w http.ResponseWriter, rq *http.Request) {
	user2.ReadUsersFromFilesystem()
	redirectTo := rq.Referer()
	if redirectTo == "" {
		redirectTo = "/hypha/" + cfg.UserHypha
	}
	http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
}

func handlerAdminUsers(w http.ResponseWriter, rq *http.Request) {
	// Get a sorted list of users
	var users []*user2.User
	for u := range user2.YieldUsers() {
		users = append(users, u)
	}

	sort.Slice(users, func(i, j int) bool {
		less := users[i].RegisteredAt.Before(users[j].RegisteredAt)
		return less
	})
	viewList(viewutil2.MetaFrom(w, rq), users)
}

func handlerAdminUserEdit(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := user2.ByName(vars["username"])
	if u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.FormDataFromRequest(rq, []string{"group"})

	if rq.Method == http.MethodPost {
		oldGroup := u.Group
		newGroup := f.Get("group")

		if user2.ValidGroup(newGroup) {
			u.Group = newGroup
			if err := user2.SaveUserDatabase(); err != nil {
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

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))

	viewEditUser(viewutil2.MetaFrom(w, rq), f, u)
}

func handlerAdminUserChangePassword(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := user2.ByName(vars["username"])
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
			if err := user2.SaveUserDatabase(); err != nil {
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

	viewEditUser(viewutil2.MetaFrom(w, rq), f, u)
}

func handlerAdminUserDelete(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	u := user2.ByName(vars["username"])
	if u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.NewFormData()

	if rq.Method == http.MethodPost {
		f = f.WithError(user2.DeleteUser(u.Name))
		if !f.HasError() {
			http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
		} else {
			log.Println(f.Error())
		}
	}

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	viewDeleteUser(viewutil2.MetaFrom(w, rq), f, u)
}

func handlerAdminUserNew(w http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodGet {
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		viewNewUser(viewutil2.MetaFrom(w, rq), util.NewFormData())
	} else if rq.Method == http.MethodPost {
		// Create a user
		f := util.FormDataFromRequest(rq, []string{"name", "password", "group"})

		err := user2.Register(f.Get("name"), f.Get("password"), f.Get("group"), "local", true)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
			viewNewUser(viewutil2.MetaFrom(w, rq), f.WithError(err))
		} else {
			http.Redirect(w, rq, "/admin/users/", http.StatusSeeOther)
		}
	}
}
