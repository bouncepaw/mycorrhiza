// Package viewutil provides utilities and common templates for views across all packages.
package viewutil

import (
	"embed"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"html/template"
	"log"
	"strings"
)

var (
	//go:embed viewutil.go
	fs     embed.FS
	BaseEn *template.Template
	BaseRu *template.Template
	m      = template.Must
)

const ruText = `
{{define "search by title"}}Поиск по названию{{end}}
{{define "close this dialog"}}Закрыть этот диалог{{end}}
{{define "login"}}Войти{{end}}
{{define "Register"}}Регистрация{{end}}
`

func Init() {
	dataText := fmt.Sprintf(`
{{define "wiki name"}}%s{{end}}
`, cfg.WikiName)
	BaseEn = m(m(template.ParseFS(fs, "viewutil.go")).Parse(dataText))
	if !cfg.UseAuth {
		m(BaseEn.Parse(`{{define "auth"}}{{end}}`))
	}
	if !cfg.AllowRegistration {
		m(BaseEn.Parse(`{{define "registration"}}{{end}}`))
	}
	BaseRu = m(m(BaseEn.Clone()).Parse(ruText))
}

// TODO: get rid of this
func localizedBaseWithWeirdBody(meta Meta) *template.Template {
	t := func() *template.Template {
		if meta.Locale() == "ru" {
			return BaseRu
		}
		return BaseEn
	}()
	return m(m(t.Clone()).Parse(`{{define "body"}}.Body{{end}}`))
}

type baseData struct {
	Meta          Meta
	Title         string
	HeadElements  []string
	HeaderLinks   []cfg.HeaderLink
	CommonScripts []string
	Body          string // TODO: remove
}

// Base is a temporary wrapper around BaseEn and BaseRu, meant to facilitate the migration from qtpl.
func Base(title, body string, lc *l18n.Localizer, u *user.User, headElements ...string) string {
	var w strings.Builder
	meta := Meta{
		Lc: lc,
		U:  u,
		W:  &w,
	}
	t := localizedBaseWithWeirdBody(meta)
	err := t.Execute(&w, baseData{
		Meta:          meta,
		Title:         title,
		HeadElements:  headElements,
		HeaderLinks:   cfg.HeaderLinks,
		CommonScripts: cfg.CommonScripts,
		Body:          body,
	})
	if err != nil {
		log.Println(err)
	}
	return w.String()
}
