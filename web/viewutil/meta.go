package viewutil

import (
	user2 "github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"io"
	"net/http"
)

// Meta is a bundle of common stuffs used by views, templates.
type Meta struct {
	Lc   *l18n.Localizer
	U    *user2.User
	W    io.Writer
	Addr string

	// New template additions
	HeadElements   []string
	BodyAttributes map[string]string
}

// MetaFrom makes a Meta from the given data. You are meant to further modify it.
func MetaFrom(w http.ResponseWriter, rq *http.Request) Meta {
	return Meta{
		Lc:   l18n.FromRequest(rq),
		U:    user2.FromRequest(rq),
		W:    w,
		Addr: rq.URL.Path,
	}
}

func (m Meta) Locale() string {
	return m.Lc.Locale
}

func (m Meta) LocaleIsRussian() bool {
	return m.Locale() == "ru"
}
