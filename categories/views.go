package categories

import (
	"embed"
	viewutil2 "github.com/bouncepaw/mycorrhiza/web/viewutil"
	"log"
	"sort"
	"strings"
)

const ruTranslation = `
{{define "empty cat"}}Эта категория пуста.{{end}}
{{define "cat"}}Категория{{end}}
{{define "hypha name"}}Название гифы{{end}}
{{define "categories"}}Категории{{end}}
{{define "placeholder"}}Название категории...{{end}}
{{define "remove from category title"}}Убрать гифу из этой категории{{end}}
{{define "add to category title"}}Добавить гифу в эту категорию{{end}}
{{define "category list"}}Список категорий{{end}}
{{define "no categories"}}В этой вики нет категорий.{{end}}
{{define "category x"}}Категория {{. | beautifulName}}{{end}}

{{define "edit category x"}}Редактирование категории {{beautifulName .}}{{end}}
{{define "edit category heading"}}Редактирование категории <a href="/category/{{.}}">{{beautifulName .}}</a>{{end}}
{{define "add"}}Добавить{{end}}
{{define "remove hyphae"}}Убрать гифы из этой категории{{end}}
{{define "remove"}}Убрать{{end}}
{{define "edit"}}Редактировать{{end}}
`

var (
	//go:embed *.html
	fs                                                         embed.FS
	viewListChain, viewPageChain, viewCardChain, viewEditChain viewutil2.Chain
)

func prepareViews() {
	viewCardChain = viewutil2.CopyEnRuWith(fs, "view_card.html", ruTranslation)
	viewListChain = viewutil2.CopyEnRuWith(fs, "view_list.html", ruTranslation)
	viewPageChain = viewutil2.CopyEnRuWith(fs, "view_page.html", ruTranslation)
	viewEditChain = viewutil2.CopyEnRuWith(fs, "view_edit.html", ruTranslation)
}

type cardData struct {
	HyphaName               string
	Categories              []string
	GivenPermissionToModify bool
}

// CategoryCard is the little sidebar that is shown nearby the hypha view.
func CategoryCard(meta viewutil2.Meta, hyphaName string) string {
	var buf strings.Builder
	err := viewCardChain.Get(meta).ExecuteTemplate(&buf, "category card", cardData{
		HyphaName:               hyphaName,
		Categories:              CategoriesWithHypha(hyphaName),
		GivenPermissionToModify: meta.U.CanProceed("add-to-category"),
	})
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

type catData struct {
	*viewutil2.BaseData
	CatName                 string
	Hyphae                  []string
	GivenPermissionToModify bool
}

func categoryEdit(meta viewutil2.Meta, catName string) {
	viewutil2.ExecutePage(meta, viewEditChain, catData{
		BaseData: &viewutil2.BaseData{
			Addr: "/edit-category/" + catName,
		},
		CatName:                 catName,
		Hyphae:                  hyphaeInCategory(catName),
		GivenPermissionToModify: meta.U.CanProceed("add-to-category"),
	})
}

func categoryPage(meta viewutil2.Meta, catName string) {
	viewutil2.ExecutePage(meta, viewPageChain, catData{
		BaseData: &viewutil2.BaseData{
			Addr: "/category/" + catName,
		},
		CatName:                 catName,
		Hyphae:                  hyphaeInCategory(catName),
		GivenPermissionToModify: meta.U.CanProceed("add-to-category"),
	})
}

type listData struct {
	*viewutil2.BaseData
	Categories []string
}

func categoryList(meta viewutil2.Meta) {
	cats := listOfCategories()
	sort.Strings(cats)
	viewutil2.ExecutePage(meta, viewListChain, listData{
		BaseData: &viewutil2.BaseData{
			Addr: "/category",
		},
		Categories: cats,
	})
}
