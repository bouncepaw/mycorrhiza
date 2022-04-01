package views

import (
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"html/template"
	"io"
	"log"
	"strings"
)

const adminTranslationRu = `
{{define "panel title"}}Панель админстратора{{end}}
{{define "panel safe section title"}}Безопасная секция{{end}}
{{define "panel link about"}}Об этой вики{{end}}
{{define "panel update header"}}Обновить ссылки в верхней панели{{end}}
{{define "panel link user list"}}Список пользователей{{end}}
{{define "panel users"}}Управление пользователями{{end}}
{{define "panel unsafe section title"}}Опасная секция{{end}}
{{define "panel shutdown"}}Выключить вики{{end}}
{{define "panel reindex hyphae"}}Переиндексировать гифы{{end}}
`

var (
	adminTemplatesEn *template.Template
	adminTemplatesRu *template.Template
)

func localizedAdminTemplates(meta viewutil.Meta) *template.Template {
	if meta.Lc.Locale == "ru" {
		return adminTemplatesRu
	}
	return adminTemplatesEn
}

func templateAsString(temp *template.Template, name string, datum ...interface{}) string {
	var buf strings.Builder
	var err error
	if len(datum) == 1 {
		err = temp.ExecuteTemplate(&buf, name, datum[0])
	} else {
		err = temp.ExecuteTemplate(&buf, name, nil)
	}
	if err != nil {
		log.Println(err)
		return ""
	}
	return buf.String()
}

func init() {
	adminTemplatesEn = template.Must(
		template.
			New("admin").
			Funcs(template.FuncMap{
				"beautifulName": util.BeautifulName,
			}).
			ParseFS(fs, "admin.html"))
	adminTemplatesRu = template.Must(
		template.Must(adminTemplatesEn.Clone()).Parse(adminTranslationRu))
}

func AdminPanel(meta viewutil.Meta) {
	var buf strings.Builder
	err := localizedAdminTemplates(meta).ExecuteTemplate(&buf, "panel", nil)
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(meta.W, Base(
		meta,
		templateAsString(localizedAdminTemplates(meta), "panel title"),
		buf.String(),
	))
}
