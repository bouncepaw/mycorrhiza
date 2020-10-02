package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/templates"
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
	util.HTTP200Page(w, base("Rename "+hyphaName+"?", templates.RenameAskHTML(hyphaName, isOld)))
}

func handlerRenameConfirm(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "rename-confirm")
		hyphaData, isOld = HyphaStorage[hyphaName]
		newName          = CanonicalName(rq.PostFormValue("new-name"))
		_, newNameIsUsed = HyphaStorage[newName]
	)
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
		if hop := hyphaData.RenameHypha(hyphaName, newName); len(hop.Errs) == 0 {
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
	util.HTTP200Page(w, base("Delete "+hyphaName+"?", templates.DeleteAskHTML(hyphaName, isOld)))
}

// handlerDeleteConfirm deletes a hypha for sure
func handlerDeleteConfirm(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "delete-confirm")
		hyphaData, isOld = HyphaStorage[hyphaName]
	)
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
	if isOld {
		textAreaFill, err = FetchTextPart(hyphaData)
		if err != nil {
			log.Println(err)
			HttpErr(w, http.StatusInternalServerError, hyphaName, "Error",
				"Could not fetch text data")
			return
		}
	} else {
		warning = `<p>You are creating a new hypha.</p>`
	}
	util.HTTP200Page(w, base("Edit"+hyphaName, templates.EditHTML(hyphaName, textAreaFill, warning)))
}

// handlerUploadText uploads a new text part for the hypha.
func handlerUploadText(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName        = HyphaNameFromRq(rq, "upload-text")
		hyphaData, isOld = HyphaStorage[hyphaName]
		textData         = rq.PostFormValue("text")
		textDataBytes    = []byte(textData)
		fullPath         = filepath.Join(WikiDir, hyphaName+"&.gmi")
	)
	if textData == "" {
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error",
			"No text data passed")
		return
	}
	// For some reason, only 0777 works. Why?
	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		log.Println(err)
	}
	if err := ioutil.WriteFile(fullPath, textDataBytes, 0644); err != nil {
		log.Println(err)
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error",
			fmt.Sprintf("Failed to write %d bytes to %s",
				len(textDataBytes), fullPath))
		return
	}
	if !isOld {
		HyphaStorage[hyphaName] = &HyphaData{
			textType: TextGemini,
			textPath: fullPath,
		}
	} else {
		hyphaData.textType = TextGemini
		hyphaData.textPath = fullPath
	}
	http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
	history.Operation(history.TypeEditText).
		WithFiles(fullPath).
		WithMsg(fmt.Sprintf("Edit ‘%s’", hyphaName)).
		WithSignature("anon").
		Apply()
}

// handlerUploadBinary uploads a new binary part for the hypha.
func handlerUploadBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "upload-binary")
	rq.ParseMultipartForm(10 << 20)
	// Read file
	file, handler, err := rq.FormFile("binary")
	if file != nil {
		defer file.Close()
	}
	// If file is not passed:
	if err != nil {
		HttpErr(w, http.StatusBadRequest, hyphaName, "Error",
			"No binary data passed")
		return
	}
	// If file is passed:
	var (
		hyphaData, isOld = HyphaStorage[hyphaName]
		mimeType         = MimeToBinaryType(handler.Header.Get("Content-Type"))
		ext              = mimeType.Extension()
		fullPath         = filepath.Join(WikiDir, hyphaName+"&"+ext)
	)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error",
			"Could not read passed data")
		return
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		log.Println(err)
	}
	if !isOld {
		HyphaStorage[hyphaName] = &HyphaData{
			binaryPath: fullPath,
			binaryType: mimeType,
		}
	} else {
		if hyphaData.binaryPath != fullPath {
			if err := history.Rename(hyphaData.binaryPath, fullPath); err != nil {
				log.Println(err)
			} else {
				log.Println("Moved", hyphaData.binaryPath, "to", fullPath)
			}
		}
		hyphaData.binaryPath = fullPath
		hyphaData.binaryType = mimeType
	}
	if err = ioutil.WriteFile(fullPath, data, 0644); err != nil {
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error",
			"Could not save passed data")
		return
	}
	log.Println("Written", len(data), "of binary data for", hyphaName, "to path", fullPath)
	http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
	history.Operation(history.TypeEditText).
		WithFiles(fullPath, hyphaData.binaryPath).
		WithMsg(fmt.Sprintf("Upload binary part for ‘%s’ with type ‘%s’", hyphaName, mimeType.Mime())).
		WithSignature("anon").
		Apply()
}
