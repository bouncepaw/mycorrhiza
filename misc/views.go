package misc

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"text/template"
)

var (
	//go:embed *html
	fs                          embed.FS
	chainList, chainTitleSearch viewutil.Chain
	ruTranslation               = `
{{define "list of hyphae"}}Список гиф{{end}}
{{define "search:"}}Поиск:{{end}}
{{define "search results for"}}Результаты поиска для «{{.}}»{{end}}
{{define "search desc"}}Название каждой из существующих гиф сопоставлено с запросом. Подходящие гифы приведены ниже.{{end}}
{{define "search no results"}}Ничего не найдено{{end}}
`
)

func initViews() {
	m := template.Must
	chainList = viewutil.
		En(viewutil.CopyEnWith(fs, "view_list.html")).
		Ru(m(viewutil.CopyRuWith(fs, "view_list.html").Parse(ruTranslation)))
	chainTitleSearch = viewutil.
		En(viewutil.CopyEnWith(fs, "view_title_search.html")).
		Ru(m(viewutil.CopyRuWith(fs, "view_title_search.html").Parse(ruTranslation)))
}

type listDatum struct {
	Name string
	Ext  string
}

type listData struct {
	*viewutil.BaseData
	Entries []listDatum
}

func viewList(meta viewutil.Meta, entries []listDatum) {
	viewutil.ExecutePage(meta, chainList, listData{
		BaseData: &viewutil.BaseData{},
		Entries:  entries,
	})
}

type titleSearchData struct {
	*viewutil.BaseData
	Query   string
	Results []string
}

func viewTitleSearch(meta viewutil.Meta, query string, results []string) {
	viewutil.ExecutePage(meta, chainTitleSearch, titleSearchData{
		BaseData: &viewutil.BaseData{},
		Query:    query,
		Results:  results,
	})
}
