package viewutil

import (
	"html/template"
	"io"
	"net/http"

	"golang.org/x/text/language"

	"github.com/bouncepaw/mycorrhiza/internal/user"
)

// Meta is a bundle of common stuffs used by views, templates.
type Meta struct {
	locale string
	U      *usersvc.User
	W      io.Writer
	Addr   string

	// New template additions
	HeadElements   []template.HTML
	BodyAttributes map[string]string
}

// MetaFrom makes a Meta from the given data. You are meant to further modify it.
func MetaFrom(w http.ResponseWriter, rq *http.Request) Meta {
	return Meta{
		locale: localeFromRequest(rq),
		U:      usersvc.FromRequest(rq),
		W:      w,
		Addr:   rq.URL.Path,
	}
}

func (m Meta) Locale() string {
	return m.locale
}

func localeFromRequest(r *http.Request) string {
	tags, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	tag, _, _ := supportedLocales.Match(tags...)
	base, _ := tag.Base()
	return base.String()
}

var supportedLocales = language.NewMatcher([]language.Tag{language.English, language.Russian})
