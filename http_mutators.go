package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
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

func handlerUnattachAsk(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "unattach-ask")
		hd, isOld = HyphaStorage[hyphaName]
		hasAmnt   = hd != nil && hd.BinaryPath != ""
	)
	if !hasAmnt {
		HttpErr(w, http.StatusBadRequest, hyphaName, "Cannot unattach", "No attachment attached yet, therefore you cannot unattach")
		log.Println("Rejected (no amnt):", rq.URL)
		return
	} else if ok := user.CanProceed(rq, "unattach-confirm"); !ok {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a trusted editor to unattach attachments")
		log.Println("Rejected (no rights):", rq.URL)
		return
	}
	util.HTTP200Page(w, base("Unattach "+hyphaName+"?", templates.UnattachAskHTML(rq, hyphaName, isOld), user.FromRequest(rq)))
}

func handlerUnattachConfirm(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "unattach-confirm")
		hyphaData, isOld = HyphaStorage[hyphaName]
		hasAmnt          = hyphaData != nil && hyphaData.BinaryPath != ""
		u                = user.FromRequest(rq)
	)
	if !u.CanProceed("unattach-confirm") {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a trusted editor to unattach attachments")
		log.Println("Rejected (no rights):", rq.URL)
		return
	}
	if !hasAmnt {
		HttpErr(w, http.StatusBadRequest, hyphaName, "Cannot unattach", "No attachment attached yet, therefore you cannot unattach")
		log.Println("Rejected (no amnt):", rq.URL)
		return
	} else if !isOld {
		// The precondition is to have the hypha in the first place.
		HttpErr(w, http.StatusPreconditionFailed, hyphaName,
			"Error: no such hypha",
			"Could not unattach this hypha because it does not exist")
		return
	}
	if hop := hyphaData.UnattachHypha(hyphaName, u); len(hop.Errs) != 0 {
		HttpErr(w, http.StatusInternalServerError, hyphaName,
			"Error: could not unattach hypha",
			fmt.Sprintf("Could not unattach this hypha due to internal errors. Server errors: <code>%v</code>", hop.Errs))
		return
	}
	http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
}

func handlerRenameAsk(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "rename-ask")
		_, isOld  = HyphaStorage[hyphaName]
		u         = user.FromRequest(rq)
	)
	if !u.CanProceed("rename-confirm") {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a trusted editor to rename pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	util.HTTP200Page(w, base("Rename "+hyphaName+"?", templates.RenameAskHTML(rq, hyphaName, isOld), u))
}

func handlerRenameConfirm(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "rename-confirm")
		_, isOld         = HyphaStorage[hyphaName]
		newName          = CanonicalName(rq.PostFormValue("new-name"))
		_, newNameIsUsed = HyphaStorage[newName]
		recursive        = rq.PostFormValue("recursive") == "true"
		u                = user.FromRequest(rq)
	)
	switch {
	case !u.CanProceed("rename-confirm"):
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a trusted editor to rename pages.")
		log.Println("Rejected", rq.URL)
	case newNameIsUsed:
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error: hypha exists",
			fmt.Sprintf("Hypha named <a href='/page/%s'>%s</a> already exists.", hyphaName, hyphaName))
	case newName == "":
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error: no name",
			"No new name is given.")
	case !isOld:
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error: no such hypha",
			"Cannot rename a hypha that does not exist yet.")
	case !HyphaPattern.MatchString(newName):
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error: invalid name",
			"Invalid new name. Names cannot contain characters <code>^?!:#@&gt;&lt;*|\"\\'&amp;%</code>")
	default:
		if hop := RenameHypha(hyphaName, newName, recursive, u); len(hop.Errs) != 0 {
			HttpErr(w, http.StatusInternalServerError, hyphaName,
				"Error: could not rename hypha",
				fmt.Sprintf("Could not rename this hypha due to an internal error. Server errors: <code>%v</code>", hop.Errs))
		} else {
			http.Redirect(w, rq, "/page/"+newName, http.StatusSeeOther)
		}
	}
}

// handlerDeleteAsk shows a delete dialog.
func handlerDeleteAsk(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "delete-ask")
		_, isOld  = HyphaStorage[hyphaName]
		u         = user.FromRequest(rq)
	)
	if !u.CanProceed("delete-ask") {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a moderator to delete pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	util.HTTP200Page(w, base("Delete "+hyphaName+"?", templates.DeleteAskHTML(rq, hyphaName, isOld), u))
}

// handlerDeleteConfirm deletes a hypha for sure
func handlerDeleteConfirm(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "delete-confirm")
		hyphaData, isOld = HyphaStorage[hyphaName]
		u                = user.FromRequest(rq)
	)
	if !u.CanProceed("delete-confirm") {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a moderator to delete pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	if !isOld {
		// The precondition is to have the hypha in the first place.
		HttpErr(w, http.StatusPreconditionFailed, hyphaName,
			"Error: no such hypha",
			"Could not delete this hypha because it does not exist.")
		return
	}
	if hop := hyphaData.DeleteHypha(hyphaName, u); len(hop.Errs) != 0 {
		HttpErr(w, http.StatusInternalServerError, hyphaName,
			"Error: could not delete hypha",
			fmt.Sprintf("Could not delete this hypha due to internal errors. Server errors: <code>%v</code>", hop.Errs))
		return
	}
	http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
}

// handlerEdit shows the edit form. It doesn't edit anything actually.
func handlerEdit(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "edit")
		hyphaData, isOld = HyphaStorage[hyphaName]
		warning          string
		textAreaFill     string
		err              error
		u                = user.FromRequest(rq)
	)
	if !u.CanProceed("edit") {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be an editor to edit pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	if isOld {
		textAreaFill, err = FetchTextPart(hyphaData)
		if err != nil {
			log.Println(err)
			HttpErr(w, http.StatusInternalServerError, hyphaName, "Error", "Could not fetch text data")
			return
		}
	} else {
		warning = `<p>You are creating a new hypha.</p>`
	}
	util.HTTP200Page(w, base("Edit "+hyphaName, templates.EditHTML(rq, hyphaName, textAreaFill, warning), u))
}

// handlerUploadText uploads a new text part for the hypha.
func handlerUploadText(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "upload-text")
		textData  = rq.PostFormValue("text")
		action    = rq.PostFormValue("action")
		u         = user.FromRequest(rq)
	)
	if !u.CanProceed("upload-text") {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be an editor to edit pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	if textData == "" {
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error", "No text data passed")
		return
	}
	if action == "Preview" {
		util.HTTP200Page(w, base("Preview "+hyphaName, templates.PreviewHTML(rq, hyphaName, textData, "", markup.Doc(hyphaName, textData).AsHTML()), u))
	} else if hop := UploadText(hyphaName, textData, u); len(hop.Errs) != 0 {
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error", hop.Errs[0].Error())
	} else {
		http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
	}
}

// handlerUploadBinary uploads a new binary part for the hypha.
func handlerUploadBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "upload-binary")
		u         = user.FromRequest(rq)
	)
	if !u.CanProceed("upload-binary") {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be an editor to upload attachments.")
		log.Println("Rejected", rq.URL)
		return
	}

	rq.ParseMultipartForm(10 << 20) // Set upload limit
	file, handler, err := rq.FormFile("binary")
	if file != nil {
		defer file.Close()
	}

	// If file is not passed:
	if err != nil {
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error", "No binary data passed")
		return
	}

	// If file is passed:
	var (
		mime = handler.Header.Get("Content-Type")
		hop  = UploadBinary(hyphaName, mime, file, u)
	)

	if len(hop.Errs) != 0 {
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error", hop.Errs[0].Error())
		return
	}
	http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
}
