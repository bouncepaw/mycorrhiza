package interwikisvc

import (
	"embed"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/interwiki/domain"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"

	"github.com/gorilla/mux"
)

var (
	//go:embed *html
	fs            embed.FS
	pageInterwiki *newtmpl.Page
	pageNameTaken *newtmpl.Page
)

func InitHandlers(rtr *mux.Router) {
	pageInterwiki = newtmpl.NewPage(fs, "view_interwiki.html")
	pageNameTaken = newtmpl.NewPage(fs, "view_name_taken.html")
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
		Engine:         interwikidomain.WikiEngine(rq.PostFormValue("engine")),
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
		slog.Info("Could not modify entry",
			"name", name,
			"reason", "does not exist")
		newtmpl.HandlerNotFound(w, rq)
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
	newtmpl.BaseData
	*Wiki
	TakenName string
	Action    string
}

func viewNameTaken(meta viewdata.Meta, wiki *Wiki, takenName, action string) {
	if err := pageNameTaken.RenderDataTo(meta, &nameTakenData{
		Wiki:      wiki,
		TakenName: takenName,
		Action:    action,
	}); err != nil {
		slog.Error("Failed to execute page", "page", "interwiki-name-taken", "err", err)
	}
}

func handlerInterwiki(w http.ResponseWriter, rq *http.Request) {
	viewInterwiki(viewdata.MetaFrom(w, rq))
}

type interwikiData struct {
	newtmpl.BaseData
	Entries []*Wiki
	CanEdit bool
	Error   string
}

func viewInterwiki(meta viewdata.Meta) {
	if err := pageInterwiki.RenderDataTo(meta, &interwikiData{
		Entries: listOfEntries,
		CanEdit: meta.U.Group == userdomain.GroupAdmin,
		Error:   "",
	}); err != nil {
		slog.Error("Failed to execute page", "page", "interwiki", "err", err)
	}
}
