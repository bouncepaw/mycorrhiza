package web

import (
	"fmt"
	"mime"
	"net/http"
	"reflect"

	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
)

func handlerUserChangePassword(w http.ResponseWriter, rq *http.Request) {
	u := usersvc.FromRequest(rq)
	// TODO: is there a better way?
	if reflect.DeepEqual(u, usersvc.EmptyUser()) || u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.NewFormData()
	if rq.Method == http.MethodPost {
		f = util.FormDataFromRequest(rq, []string{"current_password", "password", "password_confirm"})
		password := f.Get("password")
		switch {
		case !usersvc.CredentialsOK(u.Name, f.Get("current_password")):
			f = f.WithError(fmt.Errorf("incorrect password"))
		case password == "":
			f = f.WithError(fmt.Errorf("passwords should not be empty"))
		case password != f.Get("password_confirm"):
			f = f.WithError(fmt.Errorf("passwords do not match"))
		default:
			previousPassword := u.Password
			if err := u.ChangePassword(password); err != nil {
				f = f.WithError(err)
			} else if err := usersvc.SaveUserDatabase(); err != nil {
				u.Password = previousPassword
				f = f.WithError(err)
			} else {
				http.Redirect(w, rq, "/", http.StatusSeeOther)
				return
			}
		}
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}

	_ = pageChangePassword.RenderTo(
		viewdata.MetaFrom(w, rq),
		map[string]any{
			"Form": f,
			"U":    u,
		},
	)
}
