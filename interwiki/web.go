package interwikisvc

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/interwiki/domain"
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"

	"github.com/gorilla/mux"
)

var (
	pageInterwiki *newtmpl.Page
	pageNameTaken *newtmpl.Page
)

func InitHandlers(rtr *mux.Router) {
	pageInterwiki = newtmpl.NewPage("interwiki.html")
	pageNameTaken = newtmpl.NewPage("interwiki-name-taken.html")
	rtr.HandleFunc("/interwiki", handlerInterwiki)
	rtr.HandleFunc("/interwiki/add-entry", handlerAddEntry).Methods(http.MethodPost)
	rtr.HandleFunc("/interwiki/modify-entry/{target}", handlerModifyEntry).Methods(http.MethodPost)
}

// canEdit reports whether the user may change the interwiki map. Like the admin panel, it needs authorization to be on.
func canEdit(u *usersvc.User) bool {
	return cfg.UseAuth && u.CanProceed(userdomain.ActionEditInterwiki)
}

// prepareEdit checks that the user in rq may change the interwiki map, and reads and checks the entry in rq.
// If either check fails, it writes an error page and returns false.
func prepareEdit(w http.ResponseWriter, rq *http.Request) (Wiki, bool) {
	meta := viewdata.MetaFrom(w, rq)
	if !canEdit(meta.U) {
		slog.Info("No rights to change the interwiki map", "username", meta.U.Name)
		newtmpl.Error(meta, http.StatusForbidden,
			"You must be an admin to change the interwiki map.", "/interwiki", "Go back to the interwiki map.")
		return Wiki{}, false
	}

	wiki := readInterwikiEntryFromRequest(rq)
	if err := wiki.canonize(); err != nil {
		slog.Info("Invalid interwiki entry", "name", wiki.Name, "err", err)
		newtmpl.Error(meta, http.StatusBadRequest,
			"Invalid interwiki entry: "+err.Error()+".", "/interwiki", "Go back to the interwiki map.")
		return Wiki{}, false
	}
	return wiki, true
}

func readInterwikiEntryFromRequest(rq *http.Request) Wiki {
	return Wiki{WikiDTO: interwikidomain.WikiDTO{
		Name:           rq.PostFormValue("name"),
		Aliases:        strings.Split(rq.PostFormValue("aliases"), ","),
		URL:            rq.PostFormValue("url"),
		LinkHrefFormat: rq.PostFormValue("link-href-format"),
		ImgSrcFormat:   rq.PostFormValue("img-src-format"),
		Engine:         interwikidomain.WikiEngine(rq.PostFormValue("engine")),
	}}
}

func handlerModifyEntry(w http.ResponseWriter, rq *http.Request) {
	newData, ok := prepareEdit(w, rq)
	if !ok {
		return
	}
	var (
		oldData *Wiki
		name    = mux.Vars(rq)["target"]
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
	wiki, ok := prepareEdit(w, rq)
	if !ok {
		return
	}
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
		CanEdit: canEdit(meta.U),
		Error:   "",
	}); err != nil {
		slog.Error("Failed to execute page", "page", "interwiki", "err", err)
	}
}
