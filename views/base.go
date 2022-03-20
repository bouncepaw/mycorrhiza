package views

import (
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"io"
	"net/http"
)

// Meta is a bundle of common stuffs used by views, templates.
type Meta struct {
	Lc        *l18n.Localizer
	U         *user.User
	W         io.Writer
	PageTitle string
}

// MetaFrom makes a Meta from the given data. You are meant to further modify it.
func MetaFrom(w http.ResponseWriter, rq *http.Request) Meta {
	return Meta{
		Lc: l18n.FromRequest(rq),
		U:  user.FromRequest(rq),
		W:  w,
	}
}
