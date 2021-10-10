package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bouncepaw/mycomarkup/v2"
	"github.com/bouncepaw/mycomarkup/v2/mycocontext"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/history"
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
	r.PathPrefix("/delete-ask/").HandlerFunc(handlerDeleteAsk)
	r.PathPrefix("/rename-ask/").HandlerFunc(handlerRenameAsk)
	r.PathPrefix("/unattach-ask/").HandlerFunc(handlerUnattachAsk)
	// And those that do mutate something:
	r.PathPrefix("/upload-binary/").HandlerFunc(handlerUploadBinary)
	r.PathPrefix("/upload-text/").HandlerFunc(handlerUploadText)
	r.PathPrefix("/delete-confirm/").HandlerFunc(handlerDeleteConfirm)
	r.PathPrefix("/rename-confirm/").HandlerFunc(handlerRenameConfirm)
	r.PathPrefix("/unattach-confirm/").HandlerFunc(handlerUnattachConfirm)
}

func factoryHandlerAsker(
	actionPath string,
	asker func(*user.User, *hyphae.Hypha, *l18n.Localizer) (string, error),
	succTitleKey string,
	succPageTemplate func(*http.Request, string, bool) string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		util.PrepareRq(rq)
		var (
			hyphaName = util.HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
			lc        = l18n.FromRequest(rq)
		)
		if errtitle, err := asker(u, h, lc); err != nil {
			httpErr(
				w,
				lc,
				http.StatusInternalServerError,
				hyphaName,
				errtitle,
				err.Error())
			return
		}
		util.HTTP200Page(
			w,
			views.BaseHTML(
				fmt.Sprintf(lc.Get(succTitleKey), util.BeautifulName(hyphaName)),
				succPageTemplate(rq, hyphaName, h.Exists),
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

var handlerRenameAsk = factoryHandlerAsker(
	"rename-ask",
	shroom.CanRename,
	"ui.ask_rename",
	views.RenameAskHTML,
)

func factoryHandlerConfirmer(
	actionPath string,
	confirmer func(*hyphae.Hypha, *user.User, *http.Request) (*history.Op, string),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		util.PrepareRq(rq)
		var (
			hyphaName = util.HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
			lc        = l18n.FromRequest(rq)
		)
		if hop, errtitle := confirmer(h, u, rq); hop.HasErrors() {
			httpErr(w, lc, http.StatusInternalServerError, hyphaName,
				errtitle,
				hop.FirstErrorText())
			return
		}
		http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
	}
}

var handlerUnattachConfirm = factoryHandlerConfirmer(
	"unattach-confirm",
	func(h *hyphae.Hypha, u *user.User, rq *http.Request) (*history.Op, string) {
		return shroom.UnattachHypha(u, h, l18n.FromRequest(rq))
	},
)

var handlerDeleteConfirm = factoryHandlerConfirmer(
	"delete-confirm",
	func(h *hyphae.Hypha, u *user.User, rq *http.Request) (*history.Op, string) {
		return shroom.DeleteHypha(u, h, l18n.FromRequest(rq))
	},
)

// handlerRenameConfirm should redirect to the new hypha, thus it's out of factory
func handlerRenameConfirm(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u         = user.FromRequest(rq)
		lc        = l18n.FromRequest(rq)
		hyphaName = util.HyphaNameFromRq(rq, "rename-confirm")
		oldHypha  = hyphae.ByName(hyphaName)
		newName   = util.CanonicalName(rq.PostFormValue("new-name"))
		newHypha  = hyphae.ByName(newName)
		recursive = rq.PostFormValue("recursive") == "true"
	)
	hop, errtitle := shroom.RenameHypha(oldHypha, newHypha, recursive, u, lc)
	if hop.HasErrors() {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName,
			errtitle,
			hop.FirstErrorText())
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
	if errtitle, err := shroom.CanEdit(u, h, lc); err != nil {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName,
			errtitle,
			err.Error())
		return
	}
	if h.Exists {
		textAreaFill, err = shroom.FetchTextPart(h)
		if err != nil {
			log.Println(err)
			httpErr(w, lc, http.StatusInternalServerError, hyphaName,
				lc.Get("ui.error"),
				lc.Get("ui.error_text_fetch"))
			return
		}
	} else {
		warning = fmt.Sprintf(`<p class="warning warning_new-hypha">%s</p>`, lc.Get("edit.new_hypha"))
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
		hop       *history.Op
		errtitle  string
	)

	if action != "Preview" {
		hop, errtitle = shroom.UploadText(h, []byte(textData), message, u, lc)
		if hop.HasErrors() {
			httpErr(w, lc, http.StatusForbidden, hyphaName,
				errtitle,
				hop.FirstErrorText())
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
			lc.Get("ui.error"),
			err.Error())
	}
	if errtitle, err := shroom.CanAttach(u, h, lc); err != nil {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName,
			errtitle,
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
		mime          = handler.Header.Get("Content-Type")
		hop, errtitle = shroom.UploadBinary(h, mime, file, u, lc)
	)

	if hop.HasErrors() {
		httpErr(w, lc, http.StatusInternalServerError, hyphaName, errtitle, hop.FirstErrorText())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}
