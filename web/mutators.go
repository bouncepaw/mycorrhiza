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
	asker func(*user.User, *hyphae.Hypha) (error, string),
	succTitleTemplate string,
	succPageTemplate func(*http.Request, string, bool) string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		util.PrepareRq(rq)
		var (
			hyphaName = util.HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
		)
		if err, errtitle := asker(u, h); err != nil {
			httpErr(
				w,
				http.StatusInternalServerError,
				hyphaName,
				errtitle,
				err.Error())
			return
		}
		util.HTTP200Page(
			w,
			views.BaseHTML(
				fmt.Sprintf(succTitleTemplate, util.BeautifulName(hyphaName)),
				succPageTemplate(rq, hyphaName, h.Exists),
				u))
	}
}

var handlerUnattachAsk = factoryHandlerAsker(
	"unattach-ask",
	shroom.CanUnattach,
	"Unattach %s?",
	views.UnattachAskHTML,
)

var handlerDeleteAsk = factoryHandlerAsker(
	"delete-ask",
	shroom.CanDelete,
	"Delete %s?",
	views.DeleteAskHTML,
)

var handlerRenameAsk = factoryHandlerAsker(
	"rename-ask",
	shroom.CanRename,
	"Rename %s?",
	views.RenameAskHTML,
)

func factoryHandlerConfirmer(
	actionPath string,
	confirmer func(*hyphae.Hypha, *user.User, *http.Request) (*history.HistoryOp, string),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		util.PrepareRq(rq)
		var (
			hyphaName = util.HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
		)
		if hop, errtitle := confirmer(h, u, rq); hop.HasErrors() {
			httpErr(w, http.StatusInternalServerError, hyphaName,
				errtitle,
				hop.FirstErrorText())
			return
		}
		http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
	}
}

var handlerUnattachConfirm = factoryHandlerConfirmer(
	"unattach-confirm",
	func(h *hyphae.Hypha, u *user.User, _ *http.Request) (*history.HistoryOp, string) {
		return shroom.UnattachHypha(u, h)
	},
)

var handlerDeleteConfirm = factoryHandlerConfirmer(
	"delete-confirm",
	func(h *hyphae.Hypha, u *user.User, _ *http.Request) (*history.HistoryOp, string) {
		return shroom.DeleteHypha(u, h)
	},
)

// handlerRenameConfirm should redirect to the new hypha, thus it's out of factory
func handlerRenameConfirm(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		u         = user.FromRequest(rq)
		hyphaName = util.HyphaNameFromRq(rq, "rename-confirm")
		oldHypha  = hyphae.ByName(hyphaName)
		newName   = util.CanonicalName(rq.PostFormValue("new-name"))
		newHypha  = hyphae.ByName(newName)
		recursive = rq.PostFormValue("recursive") == "true"
	)
	hop, errtitle := shroom.RenameHypha(oldHypha, newHypha, recursive, u)
	if hop.HasErrors() {
		httpErr(w, http.StatusInternalServerError, hyphaName,
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
	)
	if err, errtitle := shroom.CanEdit(u, h); err != nil {
		httpErr(w, http.StatusInternalServerError, hyphaName,
			errtitle,
			err.Error())
		return
	}
	if h.Exists {
		textAreaFill, err = shroom.FetchTextPart(h)
		if err != nil {
			log.Println(err)
			httpErr(w, http.StatusInternalServerError, hyphaName,
				"Error",
				"Could not fetch text data")
			return
		}
	} else {
		warning = `<p class="warning warning_new-hypha">You are creating a new hypha.</p>`
	}
	util.HTTP200Page(
		w,
		views.BaseHTML(
			fmt.Sprintf("Edit %s", util.BeautifulName(hyphaName)),
			views.EditHTML(rq, hyphaName, textAreaFill, warning),
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
		hop       *history.HistoryOp
		errtitle  string
	)

	if action != "Preview" {
		hop, errtitle = shroom.UploadText(h, []byte(textData), message, u)
		if hop.HasErrors() {
			httpErr(w, http.StatusForbidden, hyphaName,
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
				fmt.Sprintf("Preview of %s", util.BeautifulName(hyphaName)),
				views.PreviewHTML(
					rq,
					hyphaName,
					textData,
					message,
					"",
					mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx))),
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
		file, handler, err = rq.FormFile("binary")
	)
	if err != nil {
		httpErr(w, http.StatusInternalServerError, hyphaName,
			"Error",
			err.Error())
	}
	if err, errtitle := shroom.CanAttach(u, h); err != nil {
		httpErr(w, http.StatusInternalServerError, hyphaName,
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
		hop, errtitle = shroom.UploadBinary(h, mime, file, u)
	)

	if hop.HasErrors() {
		httpErr(w, http.StatusInternalServerError, hyphaName, errtitle, hop.FirstErrorText())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}
