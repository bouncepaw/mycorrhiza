package interwikisvc

import (
	"embed"
	"log/slog"
	"net/http"
	"strings"
	"text/template"

	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/interwiki/domain"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"

	"github.com/gorilla/mux"
)

var (
	//go:embed *html
	fs            embed.FS
	tmplInterwiki *template.Template
	tmplNameTaken *template.Template
)

func InitHandlers(rtr *mux.Router) {
	tmplInterwiki = viewutil.CopyWith(fs, "view_interwiki.html")
	tmplNameTaken = viewutil.CopyWith(fs, "view_name_taken.html")
	rtr.HandleFunc("/interwiki", handlerInterwiki)
	rtr.HandleFunc("/interwiki/add-entry", handlerAddEntry).Methods(http.MethodPost)
	rtr.HandleFunc("/interwiki/modify-entry/{target}", handlerModifyEntry).Methods(http.MethodPost)
}

func readInterwikiEntryFromRequest(rq *http.Request) Wiki {
	wiki := Wiki{WikiDTO: interwikidomain.WikiDTO{
		Name:           rq.PostFormValue("name"),
		Aliases:        strings.Split(rq.PostFormValue("aliases"), ","),
		URL:            rq.PostFormValue("url"),
		LinkHrefFormat: rq.PostFormValue("link-href-format"),
		ImgSrcFormat:   rq.PostFormValue("img-src-format"),
		Engine:         interwikidomain.WikiEngine(rq.PostFormValue("engine")),
	}}
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
		slog.Info("Could not modify entry",
			"name", name,
			"reason", "does not exist")
		viewutil.HandlerNotFound(w, rq)
		return
	}

	if err := replaceEntry(oldData, &newData); err != nil {
		slog.Info("Could not modify entry",
			"name", name,
			"reason", "one of the proposed aliases or the name is taken",
			"err", err)
		viewNameTaken(viewdata.MetaFrom(w, rq), oldData, err.Error(), "modify-entry/"+name)
		return
	}

	saveInterwikiJson()
	slog.Info("Modified entry", "name", name)
	http.Redirect(w, rq, "/interwiki", http.StatusSeeOther)
}

func handlerAddEntry(w http.ResponseWriter, rq *http.Request) {
	wiki := readInterwikiEntryFromRequest(rq)
	if err := addEntry(&wiki); err != nil {
		viewNameTaken(viewdata.MetaFrom(w, rq), &wiki, err.Error(), "add-entry")
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

func viewNameTaken(meta viewdata.Meta, wiki *Wiki, takenName, action string) {
	viewutil.ExecutePage(meta, tmplNameTaken, nameTakenData{
		BaseData:  &viewutil.BaseData{},
		Wiki:      wiki,
		TakenName: takenName,
		Action:    action,
	})
}

func handlerInterwiki(w http.ResponseWriter, rq *http.Request) {
	viewInterwiki(viewdata.MetaFrom(w, rq))
}

type interwikiData struct {
	*viewutil.BaseData
	Entries []*Wiki
	CanEdit bool
	Error   string
}

func viewInterwiki(meta viewdata.Meta) {
	viewutil.ExecutePage(meta, tmplInterwiki, interwikiData{
		BaseData: &viewutil.BaseData{},
		Entries:  listOfEntries,
		CanEdit:  meta.U.Group == userdomain.GroupAdmin,
		Error:    "",
	})
}
