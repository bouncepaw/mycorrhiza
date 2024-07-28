package categories

import (
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

const ruTranslation = `
{{define "empty cat"}}{{end}}
{{define "cat"}}{{end}}
{{define "hypha name"}}{{end}}
{{define "categories"}}Категории{{end}}
{{define "placeholder"}}Название категории...{{end}}
{{define "remove from category title"}}Убрать гифу из этой категории{{end}}
{{define "add to category title"}}{{end}}
{{define "category list"}}{{end}}
{{define "no categories"}}{{end}}
{{define "category x"}}{{end}}

{{define "edit category x"}}{{end}}
{{define "edit category heading"}}{{end}}
{{define "add"}}{{end}}
{{define "remove hyphae"}}{{end}}
{{define "remove"}}{{end}}
{{define "edit"}}{{end}}
`

type catData struct {
	*viewutil.BaseData
	CatName                 string
	Hyphae                  []string
	GivenPermissionToModify bool
}

func categoryPage(meta viewutil.Meta, catName string) {

}

type listData struct {
	*viewutil.BaseData
	Categories []string
}

func categoryList(meta viewutil.Meta) {

}
