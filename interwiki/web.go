package interwiki

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

var (
	//go:embed *html
	fs            embed.FS
	ruTranslation = `
{{define "interwiki map"}}Интеркарта{{end}}
{{define "name"}}Название:{{end}}
{{define "aliases"}}Псевдонимы:{{end}}
{{define "aliases (,)"}}Псевдонимы (разделённые запятыми):{{end}}
{{define "engine"}}Движок:{{end}}
	{{define "engine/mycorrhiza"}}Микориза{{end}}
	{{define "engine/betula"}}Бетула{{end}}
	{{define "engine/agora"}}Агора{{end}}
	{{define "engine/generic"}}Любой сайт{{end}}
{{define "link href format"}}Строка форматирования атрибута href ссылки:{{end}}
{{define "img src format"}}Строка форматирования атрибута src изображения:{{end}}
{{define "unset map"}}Интеркарта не задана.{{end}}
{{define "documentation."}}Документация.{{end}}
{{define "edit separately."}}Изменяйте записи по отдельности.{{end}}
{{define "add interwiki entry"}}Добавить запись в интеркарту{{end}}
`
	chainInterwiki viewutil.Chain
	chainNameTaken viewutil.Chain
)

func InitHandlers(rtr *mux.Router) {
	chainInterwiki = viewutil.CopyEnRuWith(fs, "view_interwiki.html", ruTranslation)
	chainNameTaken = viewutil.CopyEnRuWith(fs, "view_name_taken.html", ruTranslation)
	rtr.HandleFunc("/interwiki", handlerInterwiki)
	rtr.HandleFunc("/interwiki/add-entry", handlerAddEntry).Methods(http.MethodPost)
	rtr.HandleFunc("/interwiki/modify-entry/{target}", handlerModifyEntry).Methods(http.MethodPost)
}

func readInterwikiEntryFromRequest(rq *http.Request) Wiki {
	wiki := Wiki{
		Name:           rq.PostFormValue("name"),
		Aliases:        strings.Split(rq.PostFormValue("aliases"), ","),
		URL:            rq.PostFormValue("url"),
		LinkHrefFormat: rq.PostFormValue("link-href-format"),
		ImgSrcFormat:   rq.PostFormValue("img-src-format"),
		Engine:         WikiEngine(rq.PostFormValue("engine")),
	}
	wiki.canonize()
	return wiki
}

func handlerModifyEntry(w http.ResponseWriter, rq *http.Request) {
	var (
		oldData *Wiki
		ok      bool
		name    = mux.Vars(rq)["target"]
		newData = readInterwikiEntryFromRequest(rq)
	)

	if oldData, ok = entriesByName[name]; !ok {
		log.Printf("Could not modify interwiki entry ‘%s’ because it does not exist", name)
		viewutil.HandlerNotFound(w, rq)
		return
	}

	if err := replaceEntry(oldData, &newData); err != nil {
		log.Printf("Could not modify interwiki entry ‘%s’ because one of the proposed aliases/name is taken\n", name)
		viewNameTaken(viewutil.MetaFrom(w, rq), oldData, err.Error(), "modify-entry/"+name)
		return
	}

	saveInterwikiJson()
	log.Printf("Modified interwiki entry ‘%s’\n", name)
	http.Redirect(w, rq, "/interwiki", http.StatusSeeOther)
}

func handlerAddEntry(w http.ResponseWriter, rq *http.Request) {
	wiki := readInterwikiEntryFromRequest(rq)
	if err := addEntry(&wiki); err != nil {
		viewNameTaken(viewutil.MetaFrom(w, rq), &wiki, err.Error(), "add-entry")
		return
	}
	saveInterwikiJson()
	http.Redirect(w, rq, "/interwiki", http.StatusSeeOther)
}

type nameTakenData struct {
	*viewutil.BaseData
	*Wiki
	TakenName string
	Action    string
}

func viewNameTaken(meta viewutil.Meta, wiki *Wiki, takenName, action string) {
	viewutil.ExecutePage(meta, chainNameTaken, nameTakenData{
		BaseData:  &viewutil.BaseData{},
		Wiki:      wiki,
		TakenName: takenName,
		Action:    action,
	})
}

func handlerInterwiki(w http.ResponseWriter, rq *http.Request) {
	viewInterwiki(viewutil.MetaFrom(w, rq))
}

type interwikiData struct {
	*viewutil.BaseData
	Entries []*Wiki
	CanEdit bool
	Error   string
}

func viewInterwiki(meta viewutil.Meta) {
	viewutil.ExecutePage(meta, chainInterwiki, interwikiData{
		BaseData: &viewutil.BaseData{},
		Entries:  listOfEntries,
		CanEdit:  meta.U.Group == "admin",
		Error:    "",
	})
}
