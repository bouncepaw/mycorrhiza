package web

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v5"
	"github.com/bouncepaw/mycorrhiza/hypview"
	"github.com/bouncepaw/mycorrhiza/mycoopts"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"net/http"

	"github.com/bouncepaw/mycomarkup/v5/mycocontext"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
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
		lc   = l18n.FromRequest(rq)
		h    = hyphae.ByName(util.HyphaNameFromRq(rq, "remove-media"))
		meta = viewutil.MetaFrom(w, rq)
	)
	if !u.CanProceed("remove-media") {
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "no rights")
		return
	}
	if rq.Method == "GET" {
		util.HTTP200Page(
			w,
			viewutil.Base(
				meta,
				fmt.Sprintf(lc.Get("ui.ask_remove_media"), util.BeautifulName(h.CanonicalName())),
				hypview.RemoveMediaAsk(rq, h.CanonicalName()),
				map[string]string{},
			))
		return
	}
	switch h := h.(type) {
	case *hyphae.EmptyHypha, *hyphae.TextualHypha:
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "no media to remove")
		return
	case *hyphae.MediaHypha:
		if err := shroom.RemoveMedia(u, h); err != nil {
			viewutil.HttpErr(meta, http.StatusInternalServerError, h.CanonicalName(), err.Error())
			return
		}
	}
	http.Redirect(w, rq, "/hypha/"+h.CanonicalName(), http.StatusSeeOther)
}

func handlerDelete(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		h    = hyphae.ByName(util.HyphaNameFromRq(rq, "delete"))
		meta = viewutil.MetaFrom(w, rq)
	)

	if !u.CanProceed("delete") {
		log.Printf("%s has no rights to delete ‘%s’\n", u.Name, h.CanonicalName())
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "No rights")
		return
	}

	switch h.(type) {
	case *hyphae.EmptyHypha:
		log.Printf("%s tries to delete empty hypha ‘%s’\n", u.Name, h.CanonicalName())
		// TODO: localize
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "Cannot delete an empty hypha")
		return
	}

	if rq.Method == "GET" {
		hypview.DeleteHypha(meta, h.CanonicalName())
		return
	}

	if err := shroom.Delete(u, h.(hyphae.ExistingHypha)); err != nil {
		log.Println(err)
		viewutil.HttpErr(meta, http.StatusInternalServerError, h.CanonicalName(), err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+h.CanonicalName(), http.StatusSeeOther)
}

func handlerRename(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		lc   = l18n.FromRequest(rq)
		h    = hyphae.ByName(util.HyphaNameFromRq(rq, "rename"))
		meta = viewutil.MetaFrom(w, rq)
	)

	switch h.(type) {
	case *hyphae.EmptyHypha:
		log.Printf("%s tries to rename empty hypha ‘%s’", u.Name, h.CanonicalName())
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "Cannot rename an empty hypha") // TODO: localize
		return
	}

	if !u.CanProceed("rename") {
		log.Printf("%s has no rights to rename ‘%s’\n", u.Name, h.CanonicalName())
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "No rights")
		return
	}

	var (
		oldHypha          = h.(hyphae.ExistingHypha)
		newName           = util.CanonicalName(rq.PostFormValue("new-name"))
		recursive         = rq.PostFormValue("recursive") == "true"
		leaveRedirections = rq.PostFormValue("redirection") == "true"
	)

	if rq.Method == "GET" {
		hypview.RenameHypha(meta, h.CanonicalName())
		return
	}

	if err := shroom.Rename(oldHypha, newName, recursive, leaveRedirections, u); err != nil {
		log.Printf("%s tries to rename ‘%s’: %s", u.Name, oldHypha.CanonicalName(), err.Error())
		viewutil.HttpErr(meta, http.StatusForbidden, oldHypha.CanonicalName(), lc.Get(err.Error())) // TODO: localize
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
		meta         = viewutil.MetaFrom(w, rq)
	)
	if err := shroom.CanEdit(u, h, lc); err != nil {
		viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
		return
	}
	switch h.(type) {
	case *hyphae.EmptyHypha:
		warning = fmt.Sprintf(`<p class="warning warning_new-hypha">%s</p>`, lc.Get("edit.new_hypha"))
	default:
		textAreaFill, err = hyphae.FetchMycomarkupFile(h)
		if err != nil {
			log.Println(err)
			viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, lc.Get("ui.error_text_fetch"))
			return
		}
	}
	util.HTTP200Page(
		w,
		viewutil.Base(
			meta,
			fmt.Sprintf(lc.Get("edit.title"), util.BeautifulName(hyphaName)),
			hypview.Editor(rq, hyphaName, textAreaFill, warning),
			map[string]string{}))
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
		meta      = viewutil.MetaFrom(w, rq)
	)

	if action != "Preview" {
		if err := shroom.UploadText(h, []byte(textData), message, u); err != nil {
			viewutil.HttpErr(meta, http.StatusForbidden, hyphaName, err.Error())
			return
		}
	}

	if action == "Preview" {
		ctx, _ := mycocontext.ContextFromStringInput(textData, mycoopts.MarkupOptions(hyphaName))

		util.HTTP200Page(
			w,
			viewutil.Base(
				meta,
				fmt.Sprintf(lc.Get("edit.preview_title"), util.BeautifulName(hyphaName)),
				hypview.Preview(
					rq,
					hyphaName,
					textData,
					message,
					"",
					mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx))),
				map[string]string{},
			))
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
		meta               = viewutil.MetaFrom(w, rq)
	)
	if err != nil {
		viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
	}
	if err := shroom.CanAttach(u, h, lc); err != nil {
		viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
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
		viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}
