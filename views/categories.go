package views

import (
	"github.com/bouncepaw/mycorrhiza/hyphae/categories"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
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
{{define "placeholder"}}Имя категории...{{end}}
{{define "remove from category title"}}Убрать гифу из этой категории{{end}}
{{define "add to category title"}}Добавить гифу в эту категорию{{end}}
{{define "category list heading"}}Список категорий{{end}}
{{define "no categories"}}В этой вики нет категорий.{{end}}
{{define "category x"}}Категория {{. | beautifulName}}{{end}}
`

var (
	categoryTemplatesEn *template.Template
	categoryTemplatesRu *template.Template
)

func init() {
	categoryTemplatesEn = template.Must(template.
		New("category").
		Funcs(
			template.FuncMap{
				"beautifulName": util.BeautifulName,
			}).
		ParseFS(fs, "categories.html"))
	categoryTemplatesRu = template.Must(template.Must(categoryTemplatesEn.Clone()).Parse(categoriesRu))
}

func localizedCatTemplates(meta viewutil.Meta) *template.Template {
	if meta.Lc.Locale == "ru" {
		return categoryTemplatesRu
	}
	return categoryTemplatesEn
}

func localizedCatTemplateAsString(meta viewutil.Meta, name string, datum ...interface{}) string {
	var buf strings.Builder
	var err error
	if len(datum) == 1 {
		err = localizedCatTemplates(meta).ExecuteTemplate(&buf, name, datum[0])
	} else {
		err = localizedCatTemplates(meta).ExecuteTemplate(&buf, name, nil)
	}
	if err != nil {
		log.Println(err)
		return ""
	}
	return buf.String()
}

func categoryCard(meta viewutil.Meta, hyphaName string) string {
	var buf strings.Builder
	err := localizedCatTemplates(meta).ExecuteTemplate(&buf, "category card", struct {
		HyphaName               string
		Categories              []string
		GivenPermissionToModify bool
	}{
		hyphaName,
		categories.WithHypha(hyphaName),
		meta.U.CanProceed("add-to-category"),
	})
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

func CategoryPage(meta viewutil.Meta, catName string) {
	var buf strings.Builder
	err := localizedCatTemplates(meta).ExecuteTemplate(&buf, "category page", struct {
		CatName                 string
		Hyphae                  []string
		GivenPermissionToModify bool
	}{
		catName,
		categories.Contents(catName),
		meta.U.CanProceed("add-to-category"),
	})
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(meta.W, Base(
		meta,
		localizedCatTemplateAsString(meta, "category x", catName),
		buf.String(),
	))
	if err != nil {
		log.Println(err)
	}
}

func CategoryList(meta viewutil.Meta) {
	var buf strings.Builder
	err := localizedCatTemplates(meta).ExecuteTemplate(&buf, "category list", struct {
		Categories []string
	}{
		categories.List(),
	})
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(meta.W, Base(
		meta,
		localizedCatTemplateAsString(meta, "category list heading"),
		buf.String(),
	))
	if err != nil {
		log.Println(err)
	}
}
