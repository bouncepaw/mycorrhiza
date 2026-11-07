package newtmpl

import (
	"embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
)

//go:embed *.html
var fs embed.FS

var base = template.Must(template.ParseFS(fs, "base.html"))

type Page struct {
	Template *template.Template
}

func NewPage(fs embed.FS, tmpls ...string) *Page {
	must := template.Must
	tmpl := must(must(must(
		base.Clone()).
		Funcs(template.FuncMap{
			"beautifulName": util.BeautifulName,
			"inc":           func(i int) int { return i + 1 },
			"base": func(hyphaName string) string {
				parts := strings.Split(hyphaName, "/")
				return parts[len(parts)-1]
			},
			"beautifulLink": func(hyphaName string) template.HTML {
				return template.HTML(
					fmt.Sprintf(
						`<a href="/hypha/%s">%s</a>`, hyphaName, hyphaName))
			},
		}).
		Parse(fmt.Sprintf(`
{{define "wiki name"}}%s{{end}}
{{define "user hypha"}}%s{{end}}
`, cfg.WikiName, cfg.UserHypha))).
		ParseFS(fs, tmpls...))

	if cfg.UseAuth {
		tmpl = must(tmpl.Parse(`
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
		must(tmpl.Parse(`{{define "registration"}}
{{if .Meta.U.Group | eq "anon"}}
	 <li class="auth-links__box auth-links__register-box">
		 <a href="/register" class="auth-links__link auth-links__register-link">
			 Register
		 </a>
	 </li>
{{end}}
{{end}}`))
	}

	return &Page{Template: tmpl}
}

func (p *Page) RenderTo(meta viewdata.Meta, data map[string]any) error {
	var common BaseData
	common.setBaseValuesFrom(meta)
	data["Meta"] = common.Meta
	data["HeadElements"] = common.HeadElements
	data["BodyAttributes"] = common.BodyAttributes
	data["CommonScripts"] = common.CommonScripts
	data["EditScripts"] = common.EditScripts
	data["HeaderLinks"] = common.HeaderLinks
	data["UseAuth"] = common.UseAuth

	return p.Template.ExecuteTemplate(meta.W, "page", data)
}

type BaseData struct {
	Meta           viewdata.Meta
	HeadElements   []template.HTML
	BodyAttributes map[string]string
	CommonScripts  []string
	EditScripts    []string
	HeaderLinks    []viewdata.HeaderLink
	UseAuth        bool
	Addr           string
}

func (bd *BaseData) setBaseValuesFrom(meta viewdata.Meta) {
	bd.Meta = meta
	bd.HeadElements = meta.HeadElements
	bd.BodyAttributes = meta.BodyAttributes
	bd.CommonScripts = cfg.CommonScripts
	bd.EditScripts = cfg.EditScripts
	bd.HeaderLinks = viewdata.HeaderLinks
	bd.UseAuth = cfg.UseAuth
}

func (p *Page) RenderDataTo(meta viewdata.Meta, data interface {
	setBaseValuesFrom(viewdata.Meta)
}) error {
	data.setBaseValuesFrom(meta)
	return p.Template.ExecuteTemplate(meta.W, "page", data)
}
