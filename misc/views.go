package misc

import (
	"embed"
	"iter"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

var (
	//go:embed *html
	fs                          embed.FS
	chainList, chainTitleSearch viewutil.Chain
	ruTranslation               = `
{{define "hypha count"}}{{.HyphaCount}} {{if eq .RussianHyphaCountForm "one"}}гифа{{else if eq .RussianHyphaCountForm "few"}}гифы{{else}}гиф{{end}}{{end}}
{{define "search:"}}Поиск: {{.}}{{end}}
{{define "search results for"}}Результаты поиска для «{{.}}»{{end}}
{{define "search no results"}}Ничего не найдено.{{end}}
{{define "go to hypha"}}Перейти к гифе <a class="wikilink{{if .HasExactMatch | not}} wikilink_new{{end}}" href="/hypha/{{.MatchedHyphaName}}">{{beautifulName .MatchedHyphaName}}</a>.{{end}}
`
)

func initViews() {
	chainList = viewutil.CopyEnRuWith(fs, "view_list.html", ruTranslation)
	chainTitleSearch = viewutil.CopyEnRuWith(fs, "view_title_search.html", ruTranslation)
}

type listDatum struct {
	Name     string
	Ext      string
	Label    string
	Exists   bool
	Children []*listDatum
}

func (d listDatum) Supernames() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		parts := strings.Split(d.Name, "/")
		for i, part := range parts {
			if !yield(strings.Join(parts[:i+1], "/"), part) {
				return
			}
		}
	}
}

func nestList(entries []listDatum) []*listDatum {
	var (
		root  = &listDatum{}
		nodes = map[string]*listDatum{}
	)
	for _, entry := range entries {
		node := root
		for name, part := range entry.Supernames() {
			child := nodes[name]
			if child == nil {
				child = &listDatum{Name: name, Label: part}
				nodes[name] = child
				node.Children = append(node.Children, child)
			}
			node = child
		}
		node.Ext = entry.Ext
		node.Exists = true
	}
	return root.Children
}

type listData struct {
	*viewutil.BaseData
	Entries    []*listDatum
	HyphaCount int
}

func (d listData) RussianHyphaCountForm() string {
	lastDigit, lastTwoDigits := d.HyphaCount%10, d.HyphaCount%100
	switch {
	case lastDigit == 1 && lastTwoDigits != 11:
		return "one"
	case lastDigit >= 2 && lastDigit <= 4 && (lastTwoDigits < 12 || lastTwoDigits > 14):
		return "few"
	default:
		return "many"
	}
}

func viewList(meta viewutil.Meta, entries []listDatum) {
	viewutil.ExecutePage(meta, chainList, listData{
		BaseData:   &viewutil.BaseData{},
		Entries:    nestList(entries),
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
