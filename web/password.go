package web

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"mime"
	"net/http"
	"reflect"

	"github.com/bouncepaw/mycorrhiza/user"
)

func handlerUserChangePassword(w http.ResponseWriter, rq *http.Request) {
	u := user.FromRequest(rq)
	// TODO: is there a better way?
	if reflect.DeepEqual(u, user.EmptyUser()) || u == nil {
		util.HTTP404Page(w, "404 page not found")
		return
	}

	f := util.FormDataFromRequest(rq, []string{"current_password", "password", "password_confirm"})
	currentPassword := f.Get("current_password")

	if user.CredentialsOK(u.Name, currentPassword) {
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
				if err := user.SaveUserDatabase(); err != nil {
					u.Password = previousPassword
					f = f.WithError(err)
				} else {
					http.Redirect(w, rq, "/", http.StatusSeeOther)
					return
				}
			}
		} else {
			err := fmt.Errorf("passwords do not match")
			f = f.WithError(err)
		}
	} else {
		// TODO: handle first attempt different
		err := fmt.Errorf("incorrect password")
		f = f.WithError(err)
	}

	if f.HasError() {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))

	_ = pageChangePassword.RenderTo(
		viewutil.MetaFrom(w, rq),
		map[string]any{
			"Form": f,
			"U":    u,
		},
	)
}
