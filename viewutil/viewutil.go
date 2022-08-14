// Package viewutil provides utilities and common templates for views across all packages.
package viewutil

import (
	"embed"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"io/fs"
	"log"
	"strings"
	"text/template" // TODO: save the world
)

var (
	//go:embed *.html
	fsys   embed.FS
	BaseEn *template.Template
	BaseRu *template.Template
	m      = template.Must
)

const ruText = `
{{define "search by title"}}Поиск по названию{{end}}
{{define "close this dialog"}}Закрыть этот диалог{{end}}
{{define "login"}}Войти{{end}}
{{define "register"}}Регистрация{{end}}
{{define "confirm"}}Подтвердить{{end}}
{{define "cancel"}}Отмена{{end}}
{{define "save"}}Сохранить{{end}}
{{define "error"}}Ошибка{{end}}
{{define "delete"}}Удалить{{end}}
`

func Init() {
	dataText := fmt.Sprintf(`
{{define "wiki name"}}%s{{end}}
{{define "user hypha"}}%s{{end}}
`, cfg.WikiName, cfg.UserHypha)
	BaseEn = m(m(template.New("").
		Funcs(template.FuncMap{
			"beautifulName": util.BeautifulName,
			"inc":           func(i int) int { return i + 1 },
		}).ParseFS(fsys, "base.html")).
		Parse(dataText))
	if cfg.UseAuth {
		BaseEn = m(BaseEn.Parse(`
{{define "auth"}}
<ul class="top-bar__auth auth-links">
	<li class="auth-links__box auth-links__user-box">
		{{if .Meta.U.Group | eq "anon" }}
			<a href="/login" class="auth-links__link auth-links__login-link">
				{{block "login" .}}Login{{end}}
			</a>
		{{else}}
			<a href="/hypha/{{block "user hypha" .}}{{end}}/{{.Meta.U.Name}}" class="auth-links__link auth-links__user-link">
				{{beautifulName .Meta.U.Name}}
			</a>
		{{end}}
	</li>
	{{block "registration" .}}{{end}}
</ul>
{{end}}
`))
	}
	if cfg.AllowRegistration {
		m(BaseEn.Parse(`{{define "registration"}}
{{if .Meta.U.Group | eq "anon"}}
	 <li class="auth-links__box auth-links__register-box">
		 <a href="/register" class="auth-links__link auth-links__register-link">
			 {{block "register" .}}Register{{end}}
		 </a>
	 </li>
{{end}}
{{end}}`))
	}
	BaseRu = m(m(BaseEn.Clone()).Parse(ruText))
}

func localizedBaseWithWeirdBody(meta Meta) *template.Template {
	t := func() *template.Template {
		if meta.Locale() == "ru" {
			return BaseRu
		}
		return BaseEn
	}()
	return m(m(t.Clone()).Parse(`
{{define "body"}}{{.Body}}{{end}}
{{define "title"}}{{.Title}}{{end}}
`))
}

type BaseData struct {
	Meta           Meta
	HeadElements   []string
	HeaderLinks    []HeaderLink
	CommonScripts  []string
	Addr           string
	Title          string // TODO: remove
	Body           string // TODO: remove
	BodyAttributes map[string]string
}

func (bd *BaseData) withBaseValues(meta Meta, headerLinks []HeaderLink, commonScripts []string) {
	bd.Meta = meta
	bd.HeaderLinks = headerLinks
	bd.CommonScripts = commonScripts
}

// Base is a temporary wrapper around BaseEn and BaseRu, meant to facilitate the migration from qtpl.
// TODO: get rid of this
func Base(meta Meta, title, body string, bodyAttributes map[string]string, headElements ...string) string {
	var w strings.Builder
	meta.W = &w
	t := localizedBaseWithWeirdBody(meta)
	err := t.ExecuteTemplate(&w, "page", BaseData{
		Meta:           meta,
		Title:          title,
		HeadElements:   headElements,
		HeaderLinks:    HeaderLinks,
		CommonScripts:  cfg.CommonScripts,
		Body:           body,
		BodyAttributes: bodyAttributes,
	})
	if err != nil {
		log.Println(err)
	}
	return w.String()
}

func CopyEnRuWith(fsys fs.FS, filename, ruTranslation string) Chain {
	return en(copyEnWith(fsys, filename)).
		ru(template.Must(copyRuWith(fsys, filename).Parse(ruTranslation)))
}

func copyEnWith(fsys fs.FS, f string) *template.Template {
	return m(m(BaseEn.Clone()).ParseFS(fsys, f))
}

func copyRuWith(fsys fs.FS, f string) *template.Template {
	return m(m(BaseRu.Clone()).ParseFS(fsys, f))
}

// ExecutePage executes template page in the given chain with the given data that has BaseData nested. It also sets some common BaseData fields
func ExecutePage(meta Meta, chain Chain, data interface {
	withBaseValues(meta Meta, headerLinks []HeaderLink, commonScripts []string)
}) {
	data.withBaseValues(meta, HeaderLinks, cfg.CommonScripts)
	if err := chain.Get(meta).ExecuteTemplate(meta.W, "page", data); err != nil {
		log.Println(err)
	}
}

// HeaderLinks is a list off current header links. Feel free to iterate it directly but do not modify it by yourself. Call ParseHeaderLinks if you need to set new header links.
var HeaderLinks []HeaderLink

// HeaderLink represents a header link. Header links are the links shown in the top gray bar.
type HeaderLink struct {
	// Href is the URL of the link. It goes <a href="here">...</a>.
	Href string
	// Display is what is shown when the link is rendered. It goes <a href="...">here</a>.
	Display string
}
