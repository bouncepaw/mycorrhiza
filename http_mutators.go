package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	http.HandleFunc("/upload-binary/", handlerUploadBinary)
	http.HandleFunc("/upload-text/", handlerUploadText)
	http.HandleFunc("/edit/", handlerEdit)
	http.HandleFunc("/delete-ask/", handlerDeleteAsk)
	http.HandleFunc("/delete-confirm/", handlerDeleteConfirm)
	http.HandleFunc("/rename-ask/", handlerRenameAsk)
	http.HandleFunc("/rename-confirm/", handlerRenameConfirm)
}

func handlerRenameAsk(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "rename-ask")
		_, isOld  = HyphaStorage[hyphaName]
	)
	if ok := user.CanProceed(rq, "rename-confirm"); !ok {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a trusted editor to rename pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	util.HTTP200Page(w, base("Rename "+hyphaName+"?", templates.RenameAskHTML(rq, hyphaName, isOld)))
}

func handlerRenameConfirm(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "rename-confirm")
		hyphaData, isOld = HyphaStorage[hyphaName]
		newName          = CanonicalName(rq.PostFormValue("new-name"))
		_, newNameIsUsed = HyphaStorage[newName]
		recursive        bool
	)
	if ok := user.CanProceed(rq, "rename-confirm"); !ok {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a trusted editor to rename pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	if rq.PostFormValue("recursive") == "true" {
		recursive = true
	}
	switch {
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
		if hop := hyphaData.RenameHypha(hyphaName, newName, recursive); len(hop.Errs) == 0 {
			http.Redirect(w, rq, "/page/"+newName, http.StatusSeeOther)
		} else {
			HttpErr(w, http.StatusInternalServerError, hyphaName,
				"Error: could not rename hypha",
				fmt.Sprintf("Could not rename this hypha due to an internal error. Server errors: <code>%v</code>", hop.Errs))
		}
	}
}

// handlerDeleteAsk shows a delete dialog.
func handlerDeleteAsk(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = HyphaNameFromRq(rq, "delete-ask")
		_, isOld  = HyphaStorage[hyphaName]
	)
	if ok := user.CanProceed(rq, "delete-ask"); !ok {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a moderator to delete pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	util.HTTP200Page(w, base("Delete "+hyphaName+"?", templates.DeleteAskHTML(rq, hyphaName, isOld)))
}

// handlerDeleteConfirm deletes a hypha for sure
func handlerDeleteConfirm(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "delete-confirm")
		hyphaData, isOld = HyphaStorage[hyphaName]
	)
	if ok := user.CanProceed(rq, "delete-confirm"); !ok {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be a moderator to delete pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	if isOld {
		// If deleted successfully
		if hop := hyphaData.DeleteHypha(hyphaName); len(hop.Errs) == 0 {
			http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
		} else {
			HttpErr(w, http.StatusInternalServerError, hyphaName,
				"Error: could not delete hypha",
				fmt.Sprintf("Could not delete this hypha due to an internal error. Server errors: <code>%v</code>", hop.Errs))
		}
	} else {
		// The precondition is to have the hypha in the first place.
		HttpErr(w, http.StatusPreconditionFailed, hyphaName,
			"Error: no such hypha",
			"Could not delete this hypha because it does not exist.")
	}
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
	)
	if ok := user.CanProceed(rq, "edit"); !ok {
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
	util.HTTP200Page(w, base("Edit "+hyphaName, templates.EditHTML(rq, hyphaName, textAreaFill, warning)))
}

// handlerUploadText uploads a new text part for the hypha.
func handlerUploadText(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "upload-text")
		hyphaData, isOld = HyphaStorage[hyphaName]
		textData         = rq.PostFormValue("text")
	)
	if ok := user.CanProceed(rq, "upload-text"); !ok {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be an editor to edit pages.")
		log.Println("Rejected", rq.URL)
		return
	}
	if !isOld {
		hyphaData = &HyphaData{}
	}
	if textData == "" {
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error", "No text data passed")
		return
	}
	if hop := hyphaData.UploadText(hyphaName, textData, isOld); len(hop.Errs) != 0 {
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error", hop.Errs[0].Error())
	} else {
		http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
	}
}

// handlerUploadBinary uploads a new binary part for the hypha.
func handlerUploadBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "upload-binary")
	if ok := user.CanProceed(rq, "upload-binary"); !ok {
		HttpErr(w, http.StatusForbidden, hyphaName, "Not enough rights", "You must be an editor to upload attachments.")
		log.Println("Rejected", rq.URL)
		return
	}
	rq.ParseMultipartForm(10 << 20)

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
		hyphaData, isOld = HyphaStorage[hyphaName]
		mime             = handler.Header.Get("Content-Type")
	)
	if !isOld {
		hyphaData = &HyphaData{}
	}
	hop := hyphaData.UploadBinary(hyphaName, mime, file, isOld)

	if len(hop.Errs) != 0 {
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error", hop.Errs[0].Error())
	} else {
		http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
	}
}
