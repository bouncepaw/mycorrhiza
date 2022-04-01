package categories

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"strings"
	"text/template" // TODO: Fight
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
	//go:embed *.html
	fs                           embed.FS
	m                            = template.Must
	baseEn, baseRu               *template.Template
	viewListChain, viewPageChain viewutil.Chain
	categoryTemplatesEn          *template.Template
	categoryTemplatesRu          *template.Template
)

func prepareViews() {
	categoryTemplatesEn = template.Must(template.
		New("category").
		Funcs(
			template.FuncMap{
				"beautifulName": util.BeautifulName,
			}).
		ParseFS(fs, "categories.html"))
	categoryTemplatesRu = template.Must(template.Must(categoryTemplatesEn.Clone()).Parse(categoriesRu))

	baseEn = m(viewutil.BaseEn.Clone())
	baseRu = m(viewutil.BaseRu.Clone())

	viewListChain = viewutil.
		En(
			m(m(baseEn.Clone()).ParseFS(fs, "view_list.html"))).
		Ru(
			m(m(m(baseRu.Clone()).ParseFS(fs, "view_list.html")).Parse(categoriesRu)))
	viewPageChain = viewutil.
		En(
			m(m(baseEn.Clone()).ParseFS(fs, "view_page.html"))).
		Ru(
			m(m(m(baseRu.Clone()).ParseFS(fs, "view_page.html")).Parse(categoriesRu)))
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

func CategoryCard(meta viewutil.Meta, hyphaName string) string {
	var buf strings.Builder
	err := localizedCatTemplates(meta).ExecuteTemplate(&buf, "category card", struct {
		HyphaName               string
		Categories              []string
		GivenPermissionToModify bool
	}{
		hyphaName,
		WithHypha(hyphaName),
		meta.U.CanProceed("add-to-category"),
	})
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

type pageData struct {
	viewutil.BaseData
	CatName                 string
	Hyphae                  []string
	GivenPermissionToModify bool
}

func categoryPage(meta viewutil.Meta, catName string) {
	if err := viewPageChain.Get(meta).ExecuteTemplate(meta.W, "page", pageData{
		BaseData: viewutil.BaseData{
			Meta:          meta,
			Title:         localizedCatTemplateAsString(meta, "category x", catName),
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		CatName:                 catName,
		Hyphae:                  Contents(catName),
		GivenPermissionToModify: meta.U.CanProceed("add-to-category"),
	}); err != nil {
		log.Println(err)
	}
}

type listData struct {
	viewutil.BaseData
	Categories []string
}

func categoryList(meta viewutil.Meta) {
	if err := viewListChain.Get(meta).ExecuteTemplate(meta.W, "page", listData{
		BaseData: viewutil.BaseData{
			Meta:          meta,
			Title:         localizedCatTemplateAsString(meta, "category list heading"),
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		Categories: List(),
	}); err != nil {
		log.Println(err)
	}
}
