package misc

import (
	"embed"
	"iter"
	"strings"
	"text/template"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

var (
	//go:embed *html
	fs                        embed.FS
	tmplList, tmplTitleSearch *template.Template
)

func initViews() {
	tmplList = viewutil.CopyWith(fs, "view_list.html")
	tmplTitleSearch = viewutil.CopyWith(fs, "view_title_search.html")
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

func viewList(meta viewdata.Meta, entries []listDatum) {
	viewutil.ExecutePage(meta, tmplList, listData{
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

func viewTitleSearch(meta viewdata.Meta, query string, hyphaName string, hasExactMatch bool, results []string) {
	viewutil.ExecutePage(meta, tmplTitleSearch, titleSearchData{
		BaseData:         &viewutil.BaseData{},
		Query:            query,
		Results:          results,
		MatchedHyphaName: hyphaName,
		HasExactMatch:    hasExactMatch,
	})
}
