package web

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	hyphae2 "github.com/bouncepaw/mycorrhiza/internal/hyphae"
	shroom2 "github.com/bouncepaw/mycorrhiza/internal/shroom"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	viewutil2 "github.com/bouncepaw/mycorrhiza/web/viewutil"
	"html/template"
	"log"
	"net/http"

	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"github.com/bouncepaw/mycorrhiza/hypview"
	"github.com/bouncepaw/mycorrhiza/mycoopts"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/util"
)

func initMutators(r *mux.Router) {
	r.PathPrefix("/edit/").HandlerFunc(handlerEdit)
	r.PathPrefix("/rename/").HandlerFunc(handlerRename).Methods("GET", "POST")
	r.PathPrefix("/delete/").HandlerFunc(handlerDelete).Methods("GET", "POST")
	r.PathPrefix("/remove-media/").HandlerFunc(handlerRemoveMedia).Methods("GET", "POST")
	r.PathPrefix("/upload-binary/").HandlerFunc(handlerUploadBinary)
	r.PathPrefix("/upload-text/").HandlerFunc(handlerUploadText)
}

/// TODO: this is no longer ridiculous, but is now ugly. Gotta make it at least bearable to look at :-/

func handlerRemoveMedia(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		h    = hyphae2.ByName(util.HyphaNameFromRq(rq, "remove-media"))
		meta = viewutil2.MetaFrom(w, rq)
	)
	if !u.CanProceed("remove-media") {
		viewutil2.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "no rights")
		return
	}
	if rq.Method == "GET" {
		hypview.RemoveMedia(viewutil2.MetaFrom(w, rq), h.CanonicalName())
		return
	}
	switch h := h.(type) {
	case *hyphae2.EmptyHypha, *hyphae2.TextualHypha:
		viewutil2.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "no media to remove")
		return
	case *hyphae2.MediaHypha:
		if err := shroom2.RemoveMedia(u, h); err != nil {
			viewutil2.HttpErr(meta, http.StatusInternalServerError, h.CanonicalName(), err.Error())
			return
		}
	}
	http.Redirect(w, rq, "/hypha/"+h.CanonicalName(), http.StatusSeeOther)
}

func handlerDelete(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		h    = hyphae2.ByName(util.HyphaNameFromRq(rq, "delete"))
		meta = viewutil2.MetaFrom(w, rq)
	)

	if !u.CanProceed("delete") {
		log.Printf("%s has no rights to delete ‘%s’\n", u.Name, h.CanonicalName())
		viewutil2.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "No rights")
		return
	}

	switch h.(type) {
	case *hyphae2.EmptyHypha:
		log.Printf("%s tries to delete empty hypha ‘%s’\n", u.Name, h.CanonicalName())
		// TODO: localize
		viewutil2.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "Cannot delete an empty hypha")
		return
	}

	if rq.Method == "GET" {
		hypview.DeleteHypha(meta, h.CanonicalName())
		return
	}

	if err := shroom2.Delete(u, h.(hyphae2.ExistingHypha)); err != nil {
		log.Println(err)
		viewutil2.HttpErr(meta, http.StatusInternalServerError, h.CanonicalName(), err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+h.CanonicalName(), http.StatusSeeOther)
}

func handlerRename(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		lc   = l18n.FromRequest(rq)
		h    = hyphae2.ByName(util.HyphaNameFromRq(rq, "rename"))
		meta = viewutil2.MetaFrom(w, rq)
	)

	switch h.(type) {
	case *hyphae2.EmptyHypha:
		log.Printf("%s tries to rename empty hypha ‘%s’", u.Name, h.CanonicalName())
		viewutil2.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "Cannot rename an empty hypha") // TODO: localize
		return
	}

	if !u.CanProceed("rename") {
		log.Printf("%s has no rights to rename ‘%s’\n", u.Name, h.CanonicalName())
		viewutil2.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "No rights")
		return
	}

	var (
		oldHypha          = h.(hyphae2.ExistingHypha)
		newName           = util.CanonicalName(rq.PostFormValue("new-name"))
		recursive         = rq.PostFormValue("recursive") == "true"
		leaveRedirections = rq.PostFormValue("redirection") == "true"
	)

	if rq.Method == "GET" {
		hypview.RenameHypha(meta, h.CanonicalName())
		return
	}

	if err := shroom2.Rename(oldHypha, newName, recursive, leaveRedirections, u); err != nil {
		log.Printf("%s tries to rename ‘%s’: %s", u.Name, oldHypha.CanonicalName(), err.Error())
		viewutil2.HttpErr(meta, http.StatusForbidden, oldHypha.CanonicalName(), lc.Get(err.Error())) // TODO: localize
		return
	}
	http.Redirect(w, rq, "/hypha/"+newName, http.StatusSeeOther)
}

// handlerEdit shows the edit form. It doesn't edit anything actually.
func handlerEdit(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		lc   = l18n.FromRequest(rq)
		meta = viewutil2.MetaFrom(w, rq)

		hyphaName = util.HyphaNameFromRq(rq, "edit")
		h         = hyphae2.ByName(hyphaName)

		isNew   bool
		content string
		err     error
	)

	if err := shroom2.CanEdit(u, h, lc); err != nil {
		viewutil2.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
		return
	}

	switch h.(type) {
	case *hyphae2.EmptyHypha:
		isNew = true
	default:
		content, err = hyphae2.FetchMycomarkupFile(h)
		if err != nil {
			log.Println(err)
			viewutil2.HttpErr(meta, http.StatusInternalServerError, hyphaName, lc.Get("ui.error_text_fetch"))
			return
		}
	}
	hypview.EditHypha(meta, hyphaName, isNew, content, "", "")
}

// handlerUploadText uploads a new text part for the hypha.
func handlerUploadText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		meta = viewutil2.MetaFrom(w, rq)

		hyphaName = util.HyphaNameFromRq(rq, "upload-text")
		h         = hyphae2.ByName(hyphaName)
		_, isNew  = h.(*hyphae2.EmptyHypha)

		textData = rq.PostFormValue("text")
		action   = rq.PostFormValue("action")
		message  = rq.PostFormValue("message")
	)

	if action == "preview" {
		ctx, _ := mycocontext.ContextFromStringInput(textData, mycoopts.MarkupOptions(hyphaName))
		preview := template.HTML(mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx)))
		hypview.EditHypha(meta, hyphaName, isNew, textData, message, preview)
		return
	}

	if err := shroom2.UploadText(h, []byte(textData), message, u); err != nil {
		viewutil2.HttpErr(meta, http.StatusForbidden, hyphaName, err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}

// handlerUploadBinary uploads a new media for the hypha.
func handlerUploadBinary(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	rq.ParseMultipartForm(10 << 20) // Set upload limit
	var (
		hyphaName          = util.HyphaNameFromRq(rq, "upload-binary")
		h                  = hyphae2.ByName(hyphaName)
		u                  = user.FromRequest(rq)
		lc                 = l18n.FromRequest(rq)
		file, handler, err = rq.FormFile("binary")
		meta               = viewutil2.MetaFrom(w, rq)
	)
	if err != nil {
		viewutil2.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
	}
	if err := shroom2.CanAttach(u, h, lc); err != nil {
		viewutil2.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
	}

	// If file is not passed:
	if err != nil {
		return
	}

	// If file is passed:
	if file != nil {
		defer file.Close()
	}

	var (
		mime = handler.Header.Get("Content-Type")
	)

	if err := shroom2.UploadBinary(h, mime, file, u); err != nil {
		viewutil2.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}
