package categories

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/cfg"
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
	fs                                          embed.FS
	m                                           = template.Must
	baseEn, baseRu                              *template.Template
	viewListChain, viewPageChain, viewCardChain viewutil.Chain
)

func prepareViews() {

	baseEn = m(viewutil.BaseEn.Clone())
	baseRu = m(viewutil.BaseRu.Clone())

	viewCardChain = viewutil.
		En(
			m(m(baseEn.Clone()).ParseFS(fs, "view_card.html"))).
		Ru(
			m(m(m(baseRu.Clone()).ParseFS(fs, "view_card.html")).Parse(categoriesRu)))
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

type cardData struct {
	HyphaName               string
	Categories              []string
	GivenPermissionToModify bool
}

func CategoryCard(meta viewutil.Meta, hyphaName string) string {
	var buf strings.Builder
	err := viewCardChain.Get(meta).ExecuteTemplate(&buf, "category card", cardData{
		hyphaName,
		categoriesWithHypha(hyphaName),
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
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		CatName:                 catName,
		Hyphae:                  hyphaeInCategory(catName),
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
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		Categories: listOfCategories(),
	}); err != nil {
		log.Println(err)
	}
}
