// Package viewutil provides utilities and common templates for views across all packages.
package viewutil

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"
	"text/template" // TODO: save the world

	"github.com/bouncepaw/mycorrhiza/internal/cfg"

	"github.com/bouncepaw/mycorrhiza/util"
)

var (
	//go:embed *.html
	fsys         embed.FS
	baseTemplate *template.Template
	m            = template.Must
)

func Init() {
	dataText := fmt.Sprintf(`
{{define "wiki name"}}%s{{end}}
{{define "user hypha"}}%s{{end}}
`, cfg.WikiName, cfg.UserHypha)
	baseTemplate = m(m(template.New("").
		Funcs(template.FuncMap{
			"beautifulName": util.BeautifulName,
			"inc":           func(i int) int { return i + 1 },
		}).ParseFS(fsys, "base.html")).
		Parse(dataText))
	if cfg.UseAuth {
		baseTemplate = m(baseTemplate.Parse(`
{{define "auth"}}
<ul class="top-bar__auth auth-links">
	<li class="auth-links__box auth-links__user-box">
		{{if .Meta.U.Group | eq "anon" }}
			<a href="/login" class="auth-links__link auth-links__login-link">
				Login
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
		m(baseTemplate.Parse(`{{define "registration"}}
{{if .Meta.U.Group | eq "anon"}}
	 <li class="auth-links__box auth-links__register-box">
		 <a href="/register" class="auth-links__link auth-links__register-link">
			 Register
		 </a>
	 </li>
{{end}}
{{end}}`))
	}
}

func baseWithWeirdBody() *template.Template {
	return m(m(baseTemplate.Clone()).Parse(`
{{define "body"}}{{.Body}}{{end}}
{{define "title"}}{{.Title}}{{end}}
`))
}

type BaseData struct {
	Meta           Meta
	HeadElements   []string
	HeaderLinks    []HeaderLink
	CommonScripts  []string
	EditScripts    []string
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

// Base is a temporary wrapper around the base template, meant to facilitate the migration from qtpl.
// TODO: get rid of this
func Base(meta Meta, title, body string, bodyAttributes map[string]string, headElements ...string) string {
	var w strings.Builder
	meta.W = &w
	t := baseWithWeirdBody()
	err := t.ExecuteTemplate(&w, "page", BaseData{
		Meta:           meta,
		Title:          title,
		HeadElements:   headElements,
		HeaderLinks:    HeaderLinks,
		CommonScripts:  cfg.CommonScripts,
		EditScripts:    cfg.EditScripts,
		Body:           body,
		BodyAttributes: bodyAttributes,
	})
	if err != nil {
		slog.Info("Failed to execute the legacy Base template; proceeding anyway", "err", err)
	}
	return w.String()
}

func CopyWith(fsys fs.FS, filename string) *template.Template {
	return m(m(baseTemplate.Clone()).ParseFS(fsys, filename))
}

// ExecutePage executes template page in the given template with the given data that has BaseData nested. It also sets some common BaseData fields
func ExecutePage(meta Meta, tmpl *template.Template, data interface {
	withBaseValues(meta Meta, headerLinks []HeaderLink, commonScripts []string)
}) {
	data.withBaseValues(meta, HeaderLinks, cfg.CommonScripts)
	if err := tmpl.ExecuteTemplate(meta.W, "page", data); err != nil {
		slog.Info("Failed to execute page; proceeding anyway", "err", err)
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
