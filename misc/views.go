package misc

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"path/filepath"
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
	viewutil.BaseData
	Entries []listDatum
}

func viewList(meta viewutil.Meta) {
	// TODO: make this more effective, there are too many loops and vars
	var (
		hyphaNames  = make(chan string)
		sortedHypha = hyphae.PathographicSort(hyphaNames)
		data        []listDatum
	)
	for hypha := range hyphae.YieldExistingHyphae() {
		hyphaNames <- hypha.CanonicalName()
	}
	close(hyphaNames)
	for hyphaName := range sortedHypha {
		switch h := hyphae.ByName(hyphaName).(type) {
		case *hyphae.TextualHypha:
			data = append(data, listDatum{h.CanonicalName(), ""})
		case *hyphae.MediaHypha:
			data = append(data, listDatum{h.CanonicalName(), filepath.Ext(h.MediaFilePath())[1:]})
		}
	}

	if err := chainList.Get(meta).ExecuteTemplate(meta.W, "page", listData{
		BaseData: viewutil.BaseData{
			Meta:          meta,
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		Entries: data,
	}); err != nil {
		log.Println(err)
	}
}

type titleSearchData struct {
	viewutil.BaseData
	Query   string
	Results []string
}

func viewTitleSearch(meta viewutil.Meta, query string, results []string) {
	if err := chainTitleSearch.Get(meta).ExecuteTemplate(meta.W, "page", titleSearchData{
		BaseData: viewutil.BaseData{
			Meta:          meta,
			HeaderLinks:   cfg.HeaderLinks,
			CommonScripts: cfg.CommonScripts,
		},
		Query:   query,
		Results: results,
	}); err != nil {
		log.Println(err)
	}
}
