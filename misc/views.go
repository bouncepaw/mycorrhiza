package misc

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/viewutil"
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
{{define "search no results"}}Ничего не найдено.{{end}}
{{define "x total"}}{{.}} всего.{{end}}
{{define "go to hypha"}}Перейти к гифе <a class="wikilink{{if .HasExactMatch | not}} wikilink_new{{end}}" href="/hypha/{{.MatchedHyphaName}}">{{beautifulName .MatchedHyphaName}}</a>.{{end}}
`
)

func initViews() {
	chainList = viewutil.CopyEnRuWith(fs, "view_list.html", ruTranslation)
	chainTitleSearch = viewutil.CopyEnRuWith(fs, "view_title_search.html", ruTranslation)
}

type listDatum struct {
	Name string
	Ext  string
}

type listData struct {
	*viewutil.BaseData
	Entries    []listDatum
	HyphaCount int
}

func viewList(meta viewutil.Meta, entries []listDatum) {
	viewutil.ExecutePage(meta, chainList, listData{
		BaseData:   &viewutil.BaseData{},
		Entries:    entries,
		HyphaCount: hyphae.Count(),
	})
}

type titleSearchData struct {
	*viewutil.BaseData
	Query            string
	Results          []string
	MatchedHyphaName string
	HasExactMatch    bool
}

func viewTitleSearch(meta viewutil.Meta, query string, hyphaName string, hasExactMatch bool, results []string) {
	viewutil.ExecutePage(meta, chainTitleSearch, titleSearchData{
		BaseData:         &viewutil.BaseData{},
		Query:            query,
		Results:          results,
		MatchedHyphaName: hyphaName,
		HasExactMatch:    hasExactMatch,
	})
}
