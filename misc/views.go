package misc

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	viewutil2 "github.com/bouncepaw/mycorrhiza/web/viewutil"
)

var (
	//go:embed *html
	fs                          embed.FS
	chainList, chainTitleSearch viewutil2.Chain
	ruTranslation               = `
{{define "list of hyphae"}}Список гиф{{end}}
{{define "search:"}}Поиск: {{.}}{{end}}
{{define "search results for"}}Результаты поиска для «{{.}}»{{end}}
{{define "search no results"}}Ничего не найдено.{{end}}
{{define "x total"}}{{.}} всего.{{end}}
{{define "go to hypha"}}Перейти к гифе <a class="wikilink{{if .HasExactMatch | not}} wikilink_new{{end}}" href="/hypha/{{.MatchedHyphaName}}">{{beautifulName .MatchedHyphaName}}</a>.{{end}}
`
)

func initViews() {
	chainList = viewutil2.CopyEnRuWith(fs, "view_list.html", ruTranslation)
	chainTitleSearch = viewutil2.CopyEnRuWith(fs, "view_title_search.html", ruTranslation)
}

type listDatum struct {
	Name string
	Ext  string
}

type listData struct {
	*viewutil2.BaseData
	Entries    []listDatum
	HyphaCount int
}

func viewList(meta viewutil2.Meta, entries []listDatum) {
	viewutil2.ExecutePage(meta, chainList, listData{
		BaseData:   &viewutil2.BaseData{},
		Entries:    entries,
		HyphaCount: hyphae.Count(),
	})
}

type titleSearchData struct {
	*viewutil2.BaseData
	Query            string
	Results          []string
	MatchedHyphaName string
	HasExactMatch    bool
}

func viewTitleSearch(meta viewutil2.Meta, query string, hyphaName string, hasExactMatch bool, results []string) {
	viewutil2.ExecutePage(meta, chainTitleSearch, titleSearchData{
		BaseData:         &viewutil2.BaseData{},
		Query:            query,
		Results:          results,
		MatchedHyphaName: hyphaName,
		HasExactMatch:    hasExactMatch,
	})
}
