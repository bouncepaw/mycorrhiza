package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func init() {
	// Those that do not actually mutate anything:
	http.HandleFunc("/edit/", handlerEdit)
	http.HandleFunc("/delete-ask/", handlerDeleteAsk)
	http.HandleFunc("/rename-ask/", handlerRenameAsk)
	http.HandleFunc("/unattach-ask/", handlerUnattachAsk)
	// And those that do mutate something:
	http.HandleFunc("/upload-binary/", handlerUploadBinary)
	http.HandleFunc("/upload-text/", handlerUploadText)
	http.HandleFunc("/delete-confirm/", handlerDeleteConfirm)
	http.HandleFunc("/rename-confirm/", handlerRenameConfirm)
	http.HandleFunc("/unattach-confirm/", handlerUnattachConfirm)
}

func factoryHandlerAsker(
	actionPath string,
	asker func(*user.User, *hyphae.Hypha) (error, string),
	succTitleTemplate string,
	succPageTemplate func(*http.Request, string, bool) string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		log.Println(rq.URL)
		var (
			hyphaName = HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
		)
		if err, errtitle := asker(u, h); err != nil {
			HttpErr(
				w,
				http.StatusInternalServerError,
				hyphaName,
				errtitle,
				err.Error())
			return
		}
		util.HTTP200Page(
			w,
			base(
				fmt.Sprintf(succTitleTemplate, hyphaName),
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
		log.Println(rq.URL)
		var (
			hyphaName = HyphaNameFromRq(rq, actionPath)
			h         = hyphae.ByName(hyphaName)
			u         = user.FromRequest(rq)
		)
		if hop, errtitle := confirmer(h, u, rq); hop.HasErrors() {
			HttpErr(w, http.StatusInternalServerError, hyphaName,
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

var handlerRenameConfirm = factoryHandlerConfirmer(
	"rename-confirm",
	func(oldHypha *hyphae.Hypha, u *user.User, rq *http.Request) (*history.HistoryOp, string) {
		var (
			newName   = util.CanonicalName(rq.PostFormValue("new-name"))
			recursive = rq.PostFormValue("recursive") == "true"
			newHypha  = hyphae.ByName(newName)
		)
		return shroom.RenameHypha(oldHypha, newHypha, recursive, u)
	},
)

// handlerEdit shows the edit form. It doesn't edit anything actually.
func handlerEdit(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName    = HyphaNameFromRq(rq, "edit")
		h            = hyphae.ByName(hyphaName)
		warning      string
		textAreaFill string
		err          error
		u            = user.FromRequest(rq)
	)
	if err, errtitle := shroom.CanEdit(u, h); err != nil {
		HttpErr(w, http.StatusInternalServerError, hyphaName,
			errtitle,
			err.Error())
		return
	}
	if h.Exists {
		textAreaFill, err = shroom.FetchTextPart(h)
		if err != nil {
			log.Println(err)
			HttpErr(w, http.StatusInternalServerError, hyphaName,
				"Error",
				"Could not fetch text data")
			return
		}
	} else {
		warning = `<p class="warning warning_new-hypha">You are creating a new hypha.</p>`
	}
	util.HTTP200Page(
		w,
		base(
			"Edit "+hyphaName,
			templates.EditHTML(rq, hyphaName, textAreaFill, warning),
			u))
}

// handlerUploadText uploads a new text part for the hypha.
func handlerUploadText(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "upload-text")
		h         = hyphae.ByName(hyphaName)
		textData  = rq.PostFormValue("text")
		action    = rq.PostFormValue("action")
		u         = user.FromRequest(rq)
		hop       *history.HistoryOp
		errtitle  string
	)

	if action != "Preview" {
		hop, errtitle = shroom.UploadText(h, []byte(textData), u)
		if hop.HasErrors() {
			HttpErr(w, http.StatusForbidden, hyphaName,
				errtitle,
				hop.FirstErrorText())
			return
		}
	}

	if action == "Preview" {
		util.HTTP200Page(
			w,
			base(
				"Preview "+hyphaName,
				templates.PreviewHTML(
					rq,
					hyphaName,
					textData,
					"",
					markup.Doc(hyphaName, textData).AsHTML()),
				u))
	} else {
		http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
	}
}

// handlerUploadBinary uploads a new binary part for the hypha.
func handlerUploadBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	rq.ParseMultipartForm(10 << 20) // Set upload limit
	var (
		hyphaName          = HyphaNameFromRq(rq, "upload-binary")
		h                  = hyphae.ByName(hyphaName)
		u                  = user.FromRequest(rq)
		file, handler, err = rq.FormFile("binary")
	)
	if err != nil {
		HttpErr(w, http.StatusInternalServerError, hyphaName,
			"Error",
			err.Error())
	}
	if err, errtitle := shroom.CanAttach(u, h); err != nil {
		HttpErr(w, http.StatusInternalServerError, hyphaName,
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
		HttpErr(w, http.StatusInternalServerError, hyphaName, errtitle, hop.FirstErrorText())
		return
	}
	http.Redirect(w, rq, "/hypha/"+hyphaName, http.StatusSeeOther)
}
