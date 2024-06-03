package newtmpl

import (
	"embed"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"html/template"
	"strings"

	"github.com/bouncepaw/mycorrhiza/viewutil"
)

//go:embed *.html
var fs embed.FS

type Page struct {
	TemplateEnglish *template.Template
	TemplateRussian *template.Template
}

func NewPage(tmpl string, russianTranslation map[string]string) *Page {
	must := template.Must

	return &Page{
		TemplateEnglish: must(template.ParseFS(fs, "base.html")),
		TemplateRussian: must(must(template.ParseFS(fs, "base.html")).
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
