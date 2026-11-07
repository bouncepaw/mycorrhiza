package newtmpl

import (
	"embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
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

func (p *Page) RenderTo(meta viewutil.Meta, data map[string]any) error {
	data["Meta"] = meta
	data["HeadElements"] = meta.HeadElements
	data["BodyAttributes"] = meta.BodyAttributes
	data["CommonScripts"] = cfg.CommonScripts
	data["EditScripts"] = cfg.EditScripts
	data["HeaderLinks"] = viewutil.HeaderLinks
	data["UseAuth"] = cfg.UseAuth

	return p.Template.ExecuteTemplate(meta.W, "page", data)
}
