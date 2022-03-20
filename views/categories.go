package views

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/hyphae/categories"
	"github.com/bouncepaw/mycorrhiza/util"
	"html/template"
	"io"
	"log"
	"strings"
)

//go:embed categories.html
var fs embed.FS

const categoriesRu = `
{{define "empty cat"}}Эта категория пуста.{{end}}
{{define "add hypha"}}Добавить в категорию{{end}}
{{define "cat"}}Категория{{end}}
{{define "hypha name"}}Имя гифы{{end}}
`

var (
	categoryT *template.Template
)

func init() {
	categoryT = template.Must(template.
		New("category").
		Funcs(
			template.FuncMap{
				"beautifulName": util.BeautifulName,
			}).
		ParseFS(fs, "*"))
}

func categoryCardHTML(hyphaName string) string {
	var buf strings.Builder
	err := categoryT.ExecuteTemplate(&buf, "category card", struct {
		HyphaName  string
		Categories []string
	}{
		hyphaName,
		categories.WithHypha(hyphaName),
	})
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

func CategoryPageHTML(meta Meta, catName string) {
	var buf strings.Builder
	var t, _ = categoryT.Clone()
	if meta.Lc.Locale == "ru" {
		_, _ = t.Parse(categoriesRu)
	}
	err := t.ExecuteTemplate(&buf, "category page", struct {
		CatName string
		Hyphae  []string
	}{
		catName,
		categories.Contents(catName),
	})
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(meta.W, BaseHTML(
		"Category "+util.BeautifulName(catName),
		buf.String(),
		meta.Lc,
		meta.U,
	))
	if err != nil {
		log.Println(err)
	}
}

func CategoryListHTML(meta Meta) {
	var buf strings.Builder
	err := categoryT.ExecuteTemplate(&buf, "category list", struct {
		Categories []string
	}{
		categories.List(),
	})
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(meta.W, BaseHTML(
		"Category list",
		buf.String(),
		meta.Lc,
		meta.U,
	))
	if err != nil {
		log.Println(err)
	}
}
