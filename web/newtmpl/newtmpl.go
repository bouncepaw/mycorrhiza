package newtmpl

import (
	"embed"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"html/template"
	"strings"

	"github.com/bouncepaw/mycorrhiza/viewutil"
)

//go:embed *.html
var fs embed.FS

var base = template.Must(template.ParseFS(fs, "base.html"))

type Page struct {
	TemplateEnglish *template.Template
	TemplateRussian *template.Template
}

func NewPage(fs embed.FS, tmpl string, russianTranslation map[string]string) *Page {
	must := template.Must
	en := must(must(must(
		base.Clone()).
		Funcs(template.FuncMap{
			"beautifulName": util.BeautifulName,
			"inc":           func(i int) int { return i + 1 },
		}).
		Parse(fmt.Sprintf(`
{{define "wiki name"}}%s{{end}}
{{define "user hypha"}}%s{{end}}
`, cfg.WikiName, cfg.UserHypha))).
		ParseFS(fs, tmpl))

	if cfg.UseAuth {
		en = must(en.Parse(`
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
		must(en.Parse(`{{define "registration"}}
{{if .Meta.U.Group | eq "anon"}}
	 <li class="auth-links__box auth-links__register-box">
		 <a href="/register" class="auth-links__link auth-links__register-link">
			 {{block "register" .}}Register{{end}}
		 </a>
	 </li>
{{end}}
{{end}}`))
	}

	russianTranslation["search by title"] = "Поиск по названию"

	return &Page{
		TemplateEnglish: en,
		TemplateRussian: must(must(en.Clone()).
			Parse(translationsIntoTemplates(russianTranslation))),
	}
}

func translationsIntoTemplates(m map[string]string) string {
	var sb strings.Builder
	for k, v := range m {
		sb.WriteString(fmt.Sprintf(`{{define "%s"}}%s{{end}}
`, k, v))
	}
	return sb.String()
}

func (p *Page) RenderTo(meta viewutil.Meta, data map[string]any) error {
	data["Meta"] = meta
	data["HeadElements"] = meta.HeadElements
	data["BodyAttributes"] = meta.BodyAttributes
	data["CommonScripts"] = cfg.CommonScripts
	data["EditScripts"] = cfg.EditScripts
	data["HeaderLinks"] = viewutil.HeaderLinks

	tmpl := p.TemplateEnglish
	if meta.LocaleIsRussian() {
		tmpl = p.TemplateRussian
	}

	return tmpl.ExecuteTemplate(meta.W, "page", data)
}