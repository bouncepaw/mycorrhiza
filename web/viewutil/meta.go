package viewutil

import (
	"html/template"
	"io"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/internal/user"
)

// Meta is a bundle of common stuffs used by views, templates.
type Meta struct {
	U    *usersvc.User
	W    io.Writer
	Addr string

	// New template additions
	HeadElements   []template.HTML
	BodyAttributes map[string]string
}

// MetaFrom makes a Meta from the given data. You are meant to further modify it.
func MetaFrom(w http.ResponseWriter, rq *http.Request) Meta {
	return Meta{
		U:    usersvc.FromRequest(rq),
		W:    w,
		Addr: rq.URL.Path,
	}
}
