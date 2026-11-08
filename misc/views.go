package misc

import (
	"embed"
	"iter"
	"log/slog"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
)

var (
	//go:embed *html
	fs                        embed.FS
	pageList, pageTitleSearch *newtmpl.Page
)

func initViews() {
	pageAbout = newtmpl.NewPage(fs, "view_about.html")
	pageList = newtmpl.NewPage(fs, "view_list.html")
	pageTitleSearch = newtmpl.NewPage(fs, "view_title_search.html")
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
	newtmpl.BaseData
	Entries    []*listDatum
	HyphaCount int
}

func viewList(meta viewdata.Meta, entries []listDatum) {
	if err := pageList.RenderDataTo(meta, &listData{
		Entries:    nestList(entries),
		HyphaCount: hyphae.Count(),
	}); err != nil {
		slog.Error("Failed to execute page", "page", "list", "err", err)
	}
}

type titleSearchData struct {
	newtmpl.BaseData
	Query            string
	Results          []string
	MatchedHyphaName string
	HasExactMatch    bool
}

func viewTitleSearch(meta viewdata.Meta, query string, hyphaName string, hasExactMatch bool, results []string) {
	if err := pageTitleSearch.RenderDataTo(meta, &titleSearchData{
		Query:            query,
		Results:          results,
		MatchedHyphaName: hyphaName,
		HasExactMatch:    hasExactMatch,
	}); err != nil {
		slog.Error("Failed to execute page", "page", "title-search", "err", err)
	}
}
