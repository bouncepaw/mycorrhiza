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
	// Those that do not actually mutate anything:
	r.PathPrefix("/edit/").HandlerFunc(handlerEdit)
	r.PathPrefix("/rename/").HandlerFunc(handlerRename).Methods("GET", "POST")
	r.PathPrefix("/delete-ask/").HandlerFunc(handlerDeleteAsk)
	r.PathPrefix("/unattach-ask/").HandlerFunc(handlerUnattachAsk)
	// And those that do mutate something:
	r.PathPrefix("/upload-binary/").HandlerFunc(handlerUploadBinary)
	r.PathPrefix("/upload-text/").HandlerFunc(handlerUploadText)
	r.PathPrefix("/delete-confirm/").HandlerFunc(handlerDeleteConfirm)
	r.PathPrefix("/unattach-confirm/").HandlerFunc(handlerUnattachConfirm)
}

/// TODO: this is ridiculous, refactor heavily:

func factoryHandlerAsker(
	actionPath string,
	asker func(*user.User, hyphae.Hypha, *l18n.Localizer) error,
	succTitleKey string,
	succPageTemplate func(*http.Request, string) string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		util.PrepareRq(rq)
		var (
			hyphaName = util.HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
			lc        = l18n.FromRequest(rq)
		)
		if err := asker(u, h, lc); err != nil {
			httpErr(
				w,
				lc,
				http.StatusInternalServerError,
				hyphaName,
				err.Error())
			return
		}
		util.HTTP200Page(
			w,
			views.BaseHTML(
				fmt.Sprintf(lc.Get(succTitleKey), util.BeautifulName(hyphaName)),
				succPageTemplate(rq, hyphaName),
				lc,
				u))
	}
}

var handlerUnattachAsk = factoryHandlerAsker(
	"unattach-ask",
	shroom.CanUnattach,
	"ui.ask_unattach",
	views.UnattachAskHTML,
)

var handlerDeleteAsk = factoryHandlerAsker(
	"delete-ask",
	shroom.CanDelete,
	"ui.ask_delete",
	views.DeleteAskHTML,
)

func factoryHandlerConfirmer(
	actionPath string,
	confirmer func(hyphae.Hypha, *user.User, *http.Request) error,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		util.PrepareRq(rq)
		var (
			hyphaName = util.HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
			lc        = l18n.FromRequest(rq)
		)
		if err := confirmer(h, u, rq); err != nil {
			httpErr(w, lc, http.StatusInternalServerError, hyphaName,
				err.Error())
			return
		}
		http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
	}
}

var handlerUnattachConfirm = factoryHandlerConfirmer(
	"unattach-confirm",
	func(h hyphae.Hypha, u *user.User, rq *http.Request) error {
		return shroom.UnattachHypha(u, h, l18n.FromRequest(rq))
	},
)

var handlerDeleteConfirm = factoryHandlerConfirmer(
	"delete-confirm",
	func(h hyphae.Hypha, u *user.User, rq *http.Request) error {
		return shroom.DeleteHypha(u, h, l18n.FromRequest(rq))
	},
)

func handlerRename(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u  = user.FromRequest(rq)
		lc = l18n.FromRequest(rq)
		h  = hyphae.ByName(util.HyphaNameFromRq(rq, "rename-confirm"))
	)

	switch h.(type) {
	case *hyphae.EmptyHypha:
		log.Printf("%s tries to rename empty hypha ‘%s’", u.Name, h.CanonicalName())
		httpErr(w, lc, http.StatusForbidden, h.CanonicalName(), "Cannot rename an empty hypha") // TODO: localize
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
		if err := shroom.UploadText(h, []byte(textData), message, u, lc); err != nil {
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

// handlerUploadBinary uploads a new attachment for the hypha.
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

	if err := shroom.UploadBinary(h, mime, file, u, lc); err != nil {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName, err.Error())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}
