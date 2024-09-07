package web

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/hypview"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/shroom"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/mycoopts"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"

	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"github.com/gorilla/mux"
)

func initMutators(r *mux.Router) {
	r.PathPrefix("/edit/").HandlerFunc(handlerEdit)
	r.PathPrefix("/rename/").HandlerFunc(handlerRename).Methods("GET", "POST")
	r.PathPrefix("/delete/").HandlerFunc(handlerDelete).Methods("GET", "POST")
	r.PathPrefix("/remove-media/").HandlerFunc(handlerRemoveMedia).Methods("POST")
	r.PathPrefix("/upload-binary/").HandlerFunc(handlerUploadBinary)
	r.PathPrefix("/upload-text/").HandlerFunc(handlerUploadText)
}

/// TODO: this is no longer ridiculous, but is now ugly. Gotta make it at least bearable to look at :-/

func handlerRemoveMedia(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		h    = hyphae.ByName(util.HyphaNameFromRq(rq, "remove-media"))
		meta = viewutil.MetaFrom(w, rq)
	)
	if !u.CanProceed("remove-media") {
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "no rights")
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
		slog.Info("No rights to delete hypha",
			"username", u.Name, "hyphaName", h.CanonicalName())
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "No rights")
		return
	}

	switch h.(type) {
	case *hyphae.EmptyHypha:
		slog.Info("Trying to delete empty hyphae",
			"username", u.Name, "hyphaName", h.CanonicalName())
		// TODO: localize
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "Cannot delete an empty hypha")
		return
	}

	if rq.Method == "GET" {
		_ = pageHyphaDelete.RenderTo(
			viewutil.MetaFrom(w, rq),
			map[string]any{
				"HyphaName": h.CanonicalName(),
			})
		return
	}

	if err := shroom.Delete(u, h.(hyphae.ExistingHypha)); err != nil {
		slog.Error("Failed to delete hypha", "err", err)
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
		slog.Info("Trying to rename empty hypha",
			"username", u.Name, "hyphaName", h.CanonicalName())
		viewutil.HttpErr(meta, http.StatusForbidden, h.CanonicalName(), "Cannot rename an empty hypha") // TODO: localize
		return
	}

	if !u.CanProceed("rename") {
		slog.Info("No rights to rename hypha",
			"username", u.Name, "hyphaName", h.CanonicalName())
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
		slog.Error("Failed to rename hypha",
			"err", err, "username", u.Name, "hyphaName", oldHypha.CanonicalName())
		viewutil.HttpErr(meta, http.StatusForbidden, oldHypha.CanonicalName(), lc.Get(err.Error())) // TODO: localize
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
		meta = viewutil.MetaFrom(w, rq)

		hyphaName = util.HyphaNameFromRq(rq, "edit")
		h         = hyphae.ByName(hyphaName)

		isNew   bool
		content string
		err     error
	)

	if err := shroom.CanEdit(u, h, lc); err != nil {
		viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
		return
	}

	switch h.(type) {
	case *hyphae.EmptyHypha:
		isNew = true
	default:
		content, err = hyphae.FetchMycomarkupFile(h)
		if err != nil {
			slog.Error("Failed to fetch Mycomarkup file", "err", err)
			viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, lc.Get("ui.error_text_fetch"))
			return
		}
	}
	_ = pageHyphaEdit.RenderTo(
		viewutil.MetaFrom(w, rq),
		map[string]any{
			"HyphaName": hyphaName,
			"Content":   content,
			"IsNew":     isNew,
			"Message":   "",
			"Preview":   "",
		})
}

// handlerUploadText uploads a new text part for the hypha.
func handlerUploadText(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u    = user.FromRequest(rq)
		meta = viewutil.MetaFrom(w, rq)

		hyphaName = util.HyphaNameFromRq(rq, "upload-text")
		h         = hyphae.ByName(hyphaName)
		_, isNew  = h.(*hyphae.EmptyHypha)

		textData = rq.PostFormValue("text")
		action   = rq.PostFormValue("action")
		message  = rq.PostFormValue("message")
	)

	if action == "preview" {
		ctx, _ := mycocontext.ContextFromStringInput(textData, mycoopts.MarkupOptions(hyphaName))
		preview := template.HTML(mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx)))

		_ = pageHyphaEdit.RenderTo(
			viewutil.MetaFrom(w, rq),
			map[string]any{
				"HyphaName": hyphaName,
				"Content":   textData,
				"IsNew":     isNew,
				"Message":   message,
				"Preview":   preview,
			})
		return
	}

	if err := shroom.UploadText(h, []byte(textData), message, u); err != nil {
		viewutil.HttpErr(meta, http.StatusForbidden, hyphaName, err.Error())
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
