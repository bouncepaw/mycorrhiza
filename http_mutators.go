package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/history"
)

func init() {
	http.HandleFunc("/upload-binary/", handlerUploadBinary)
	http.HandleFunc("/upload-text/", handlerUploadText)
	http.HandleFunc("/edit/", handlerEdit)
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
	form := fmt.Sprintf(`
		<main>
			<h1>Edit %[1]s</h1>
			%[3]s
			<form method="post" class="upload-text-form"
			      action="/upload-text/%[1]s">
				<textarea name="text">%[2]s</textarea>
				<br/>
				<input type="submit"/>
				<a href="/page/%[1]s">Cancel</a>
			</form>
		</main>
`, hyphaName, textAreaFill, warning)

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base(
		"Edit "+hyphaName, form)))
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
	history.Operation(history.TypeEditText).
		WithFiles(fullPath).
		WithMsg(fmt.Sprintf("Edit ‘%s’", hyphaName)).
		WithSignature("anon").
		Apply()
	http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
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
	if err = ioutil.WriteFile(fullPath, data, 0644); err != nil {
		HttpErr(w, http.StatusInternalServerError, hyphaName, "Error",
			"Could not save passed data")
		return
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
	log.Println("Written", len(data), "of binary data for", hyphaName, "to path", fullPath)
	history.Operation(history.TypeEditText).
		WithFiles(fullPath).
		WithMsg(fmt.Sprintf("Upload binary part for ‘%s’ with type ‘%s’", hyphaName, mimeType.Mime())).
		WithSignature("anon").
		Apply()
	http.Redirect(w, rq, "/page/"+hyphaName, http.StatusSeeOther)
}
