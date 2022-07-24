package categories

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"sort"
	"strings"
)

const ruTranslation = `
{{define "empty cat"}}Эта категория пуста.{{end}}
{{define "add hypha"}}Добавить в категорию{{end}}
{{define "cat"}}Категория{{end}}
{{define "hypha name"}}Имя гифы{{end}}
{{define "categories"}}Категории{{end}}
{{define "placeholder"}}Имя категории...{{end}}
{{define "remove from category title"}}Убрать гифу из этой категории{{end}}
{{define "add to category title"}}Добавить гифу в эту категорию{{end}}
{{define "category list"}}Список категорий{{end}}
{{define "no categories"}}В этой вики нет категорий.{{end}}
{{define "category x"}}Категория {{. | beautifulName}}{{end}}
`

var (
	//go:embed *.html
	fs                                                         embed.FS
	viewListChain, viewPageChain, viewCardChain, viewEditChain viewutil.Chain
)

func prepareViews() {
	viewCardChain = viewutil.CopyEnRuWith(fs, "view_card.html", ruTranslation)
	viewListChain = viewutil.CopyEnRuWith(fs, "view_list.html", ruTranslation)
	viewPageChain = viewutil.CopyEnRuWith(fs, "view_page.html", ruTranslation)
	viewEditChain = viewutil.CopyEnRuWith(fs, "view_edit.html", ruTranslation)
}

type cardData struct {
	HyphaName               string
	Categories              []string
	GivenPermissionToModify bool
}

// CategoryCard is the little sidebar that is shown nearby the hypha view.
func CategoryCard(meta viewutil.Meta, hyphaName string) string {
	var buf strings.Builder
	err := viewCardChain.Get(meta).ExecuteTemplate(&buf, "category card", cardData{
		HyphaName:               hyphaName,
		Categories:              categoriesWithHypha(hyphaName),
		GivenPermissionToModify: meta.U.CanProceed("add-to-category"),
	})
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

type catData struct {
	*viewutil.BaseData
	CatName                 string
	Hyphae                  []string
	GivenPermissionToModify bool
}

func categoryEdit(meta viewutil.Meta, catName string) {
	viewutil.ExecutePage(meta, viewEditChain, catData{
		BaseData: &viewutil.BaseData{
			Addr: "/edit-category/" + catName,
		},
		CatName:                 catName,
		Hyphae:                  hyphaeInCategory(catName),
		GivenPermissionToModify: meta.U.CanProceed("add-to-category"),
	})
}

func categoryPage(meta viewutil.Meta, catName string) {
	viewutil.ExecutePage(meta, viewPageChain, catData{
		BaseData: &viewutil.BaseData{
			Addr: "/category/" + catName,
		},
		CatName:                 catName,
		Hyphae:                  hyphaeInCategory(catName),
		GivenPermissionToModify: meta.U.CanProceed("add-to-category"),
	})
}

type listData struct {
	*viewutil.BaseData
	Categories []string
}

func categoryList(meta viewutil.Meta) {
	cats := listOfCategories()
	sort.Strings(cats)
	viewutil.ExecutePage(meta, viewListChain, listData{
		BaseData: &viewutil.BaseData{
			Addr: "/category",
		},
		Categories: cats,
	})
}
