package views

import (
	"github.com/bouncepaw/mycorrhiza/hyphae/categories"
	"github.com/bouncepaw/mycorrhiza/util"
	"html/template"
	"io"
	"log"
	"strings"
)

const categoriesRu = `
{{define "empty cat"}}Эта категория пуста.{{end}}
{{define "add hypha"}}Добавить в категорию{{end}}
{{define "cat"}}Категория{{end}}
{{define "hypha name"}}Имя гифы{{end}}
{{define "categories"}}Категории{{end}}
{{define "placeholder"}}Имя категории{{end}}
{{define "remove from category title"}}Убрать гифу из этой категории{{end}}
{{define "add to category"}}Добавить гифу в эту категорию{{end}}
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
		ParseFS(fs, "categories.html"))
}

func categoryCard(meta Meta, hyphaName string) string {
	var buf strings.Builder
	t, err := categoryT.Clone()
	if err != nil {
		log.Println(err)
		return ""
	}
	if meta.Lc.Locale == "ru" {
		_, err = t.Parse(categoriesRu)
		if err != nil {
			log.Println(err)
			return ""
		}
	}
	err = t.ExecuteTemplate(&buf, "category card", struct {
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

func CategoryPage(meta Meta, catName string) {
	var buf strings.Builder
	var t, err = categoryT.Clone()
	if err != nil {
		log.Println(err)
		return
	}
	if meta.Lc.Locale == "ru" {
		_, err = t.Parse(categoriesRu)
		if err != nil {
			log.Println(err)
			return
		}
	}
	err = t.ExecuteTemplate(&buf, "category page", struct {
		CatName string
		Hyphae  []string
	}{
		catName,
		categories.Contents(catName),
	})
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(meta.W, Base(
		"Category "+util.BeautifulName(catName),
		buf.String(),
		meta.Lc,
		meta.U,
	))
	if err != nil {
		log.Println(err)
	}
}

func CategoryList(meta Meta) {
	var buf strings.Builder
	err := categoryT.ExecuteTemplate(&buf, "category list", struct {
		Categories []string
	}{
		categories.List(),
	})
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(meta.W, Base(
		"Category list",
		buf.String(),
		meta.Lc,
		meta.U,
	))
	if err != nil {
		log.Println(err)
	}
}
