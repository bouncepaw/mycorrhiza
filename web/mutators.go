package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bouncepaw/mycomarkup/v3"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
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
		u  = user.FromRequest(rq)
		lc = l18n.FromRequest(rq)
		h  = hyphae.ByName(util.HyphaNameFromRq(rq, "delete"))
	)
	if !u.CanProceed("remove-media") {
		httpErr(w, lc, http.StatusForbidden, h.CanonicalName(), "no rights")
		return
	}
	if rq.Method == "GET" {
		util.HTTP200Page(
			w,
			views.BaseHTML(
				fmt.Sprintf(lc.Get("ui.ask_remove_media"), util.BeautifulName(h.CanonicalName())),
				views.RemoveMediaAskHTML(rq, h.CanonicalName()),
				lc,
				u))
		return
	}
	switch h := h.(type) {
	case *hyphae.EmptyHypha, *hyphae.TextualHypha:
		httpErr(w, lc, http.StatusForbidden, h.CanonicalName(), "no media to remove")
		return
	case *hyphae.MediaHypha:
		if err := shroom.RemoveMedia(u, h); err != nil {
			httpErr(w, lc, http.StatusInternalServerError, h.CanonicalName(), err.Error())
			return
		}
	}
}

func handlerDelete(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u  = user.FromRequest(rq)
		lc = l18n.FromRequest(rq)
		h  = hyphae.ByName(util.HyphaNameFromRq(rq, "delete"))
	)

	if !u.CanProceed("delete") {
		log.Printf("%s has no rights to delete ‘%s’\n", u.Name, h.CanonicalName())
		httpErr(w, lc, http.StatusForbidden, h.CanonicalName(), "No rights")
		return
	}

	switch h.(type) {
	case *hyphae.EmptyHypha:
		log.Printf("%s tries to delete empty hypha ‘%s’\n", u.Name, h.CanonicalName())
		// TODO: localize
		httpErr(w, lc, http.StatusForbidden, h.CanonicalName(), "Cannot delete an empty hypha")
		return
	}

	if rq.Method == "GET" {
		util.HTTP200Page(
			w,
			views.BaseHTML(
				fmt.Sprintf(lc.Get("ui.ask_delete"), util.BeautifulName(h.CanonicalName())),
				views.DeleteAskHTML(rq, h.CanonicalName()),
				lc,
				u))
		return
	}

	if err := shroom.Delete(u, h.(hyphae.ExistingHypha)); err != nil {
		log.Println(err)
		httpErr(w, lc, http.StatusInternalServerError, h.CanonicalName(), err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+h.CanonicalName(), http.StatusSeeOther)
}

func handlerRename(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u  = user.FromRequest(rq)
		lc = l18n.FromRequest(rq)
		h  = hyphae.ByName(util.HyphaNameFromRq(rq, "rename"))
	)

	switch h.(type) {
	case *hyphae.EmptyHypha:
		log.Printf("%s tries to rename empty hypha ‘%s’", u.Name, h.CanonicalName())
		httpErr(w, lc, http.StatusForbidden, h.CanonicalName(), "Cannot rename an empty hypha") // TODO: localize
		return
	}

	if !u.CanProceed("rename") {
		log.Printf("%s has no rights to rename ‘%s’\n", u.Name, h.CanonicalName())
		httpErr(w, lc, http.StatusForbidden, h.CanonicalName(), "No rights")
		return
	}

	var (
		oldHypha  = h.(hyphae.ExistingHypha)
		newName   = util.CanonicalName(rq.PostFormValue("new-name"))
		recursive = rq.PostFormValue("recursive") == "true"
	)

	if rq.Method == "GET" {
		util.HTTP200Page(
			w,
			views.BaseHTML(
				fmt.Sprintf(lc.Get("ui.ask_rename"), util.BeautifulName(oldHypha.CanonicalName())),
				views.RenameAskHTML(rq, oldHypha.CanonicalName()),
				lc,
				u))
		return
	}

	if err := shroom.Rename(oldHypha, newName, recursive, u); err != nil {
		log.Printf("%s tries to rename ‘%s’: %s", u.Name, oldHypha.CanonicalName(), err.Error())
		httpErr(w, lc, http.StatusForbidden, oldHypha.CanonicalName(), lc.Get(err.Error())) // TODO: localize
		return
	}
	http.Redirect(w, rq, "/hypha/"+newName, http.StatusSeeOther)
}

// handlerEdit shows the edit form. It doesn't edit anything actually.
func handlerEdit(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		hyphaName    = util.HyphaNameFromRq(rq, "edit")
		h            = hyphae.ByName(hyphaName)
		warning      string
		textAreaFill string
		err          error
		u            = user.FromRequest(rq)
		lc           = l18n.FromRequest(rq)
	)
	if err := shroom.CanEdit(u, h, lc); err != nil {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName,
			err.Error())
		return
	}
	switch h.(type) {
	case *hyphae.EmptyHypha:
		warning = fmt.Sprintf(`<p class="warning warning_new-hypha">%s</p>`, lc.Get("edit.new_hypha"))
	default:
		textAreaFill, err = shroom.FetchTextFile(h)
		if err != nil {
			log.Println(err)
			httpErr(w, lc, http.StatusInternalServerError, hyphaName,
				lc.Get("ui.error_text_fetch"))
			return
		}
	}
	util.HTTP200Page(
		w,
		views.BaseHTML(
			fmt.Sprintf(lc.Get("edit.title"), util.BeautifulName(hyphaName)),
			views.EditHTML(rq, hyphaName, textAreaFill, warning),
			lc,
			u))
}

// handlerUploadText uploads a new text part for the hypha.
func handlerUploadText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		hyphaName = util.HyphaNameFromRq(rq, "upload-text")
		h         = hyphae.ByName(hyphaName)
		textData  = rq.PostFormValue("text")
		action    = rq.PostFormValue("action")
		message   = rq.PostFormValue("message")
		u         = user.FromRequest(rq)
		lc        = l18n.FromRequest(rq)
	)

	if action != "Preview" {
		if err := shroom.UploadText(h, []byte(textData), message, u); err != nil {
			httpErr(w, lc, http.StatusForbidden, hyphaName, err.Error())
			return
		}
	}

	if action == "Preview" {
		ctx, _ := mycocontext.ContextFromStringInput(hyphaName, textData)

		util.HTTP200Page(
			w,
			views.BaseHTML(
				fmt.Sprintf(lc.Get("edit.preview_title"), util.BeautifulName(hyphaName)),
				views.PreviewHTML(
					rq,
					hyphaName,
					textData,
					message,
					"",
					mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx))),
				lc,
				u))
	} else {
		http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
	}
}

// handlerUploadBinary uploads a new media for the hypha.
func handlerUploadBinary(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	rq.ParseMultipartForm(10 << 20) // Set upload limit
	var (
		hyphaName          = util.HyphaNameFromRq(rq, "upload-binary")
		h                  = hyphae.ByName(hyphaName)
		u                  = user.FromRequest(rq)
		lc                 = l18n.FromRequest(rq)
		file, handler, err = rq.FormFile("binary")
	)
	if err != nil {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName,
			err.Error())
	}
	if err := shroom.CanAttach(u, h, lc); err != nil {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName,
			err.Error())
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

	if err := shroom.UploadBinary(h, mime, file, u); err != nil {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName, err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}
