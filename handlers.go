package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// There are handlers below. See main() for their usage.

// Boilerplate code present in many handlers. Good to have it.
func HandlerBase(w http.ResponseWriter, rq *http.Request) (Revision, bool) {
	vars := mux.Vars(rq)
	revno := RevInMap(vars)
	return GetRevision(vars["hypha"], revno)
}

func HandlerGetBinary(w http.ResponseWriter, rq *http.Request) {
	if rev, ok := HandlerBase(w, rq); ok {
		rev.ActionGetBinary(w)
	}
}

func HandlerRaw(w http.ResponseWriter, rq *http.Request) {
	if rev, ok := HandlerBase(w, rq); ok {
		rev.ActionRaw(w)
	}
}

func HandlerZen(w http.ResponseWriter, rq *http.Request) {
	if rev, ok := HandlerBase(w, rq); ok {
		rev.ActionZen(w)
	}
}

func HandlerView(w http.ResponseWriter, rq *http.Request) {
	if rev, ok := HandlerBase(w, rq); ok {
		rev.ActionView(w, HyphaPage)
	}
}

func HandlerHistory(w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerEdit(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	ActionEdit(vars["hypha"], w)
}

func HandlerRewind(w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerDelete(w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerRename(w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

// makeTagsSlice turns strings like `"foo,, bar,kek"` to slice of strings that represent tag names. Whitespace around commas is insignificant.
// Expected output for string above: []string{"foo", "bar", "kek"}
func makeTagsSlice(responseTagsString string) (ret []string) {
	for _, tag := range strings.Split(responseTagsString, ",") {
		if trimmed := strings.TrimSpace(tag); "" == trimmed {
			ret = append(ret, trimmed)
		}
	}
	return ret
}

// getHypha returns an existing hypha if it exists in `hyphae` or creates a new one. If it `isNew`, you'll have to insert it to `hyphae` yourself.
func getHypha(name string) (*Hypha, bool) {
	log.Println("Accessing hypha", name)
	if h, ok := hyphae[name]; ok {
		log.Println("Got hypha", name)
		return h, false
	}
	log.Println("Create hypha", name)
	h := &Hypha{
		FullName:   name,
		Path:       filepath.Join(rootWikiDir, name),
		Revisions:  make(map[string]*Revision),
		parentName: filepath.Dir(name),
	}
	return h, true
}

// revisionFromHttpData creates a new revison for hypha `h`. All data is fetched from `rq`, except for BinaryMime and BinaryPath which require additional processing. You'll have te insert the revision to `h` yourself.
func revisionFromHttpData(h *Hypha, rq *http.Request) *Revision {
	idStr := strconv.Itoa(h.NewestRevisionInt() + 1)
	log.Printf("Creating revision %s from http data", idStr)
	rev := &Revision{
		Id:        h.NewestRevisionInt() + 1,
		FullName:  h.FullName,
		ShortName: filepath.Base(h.FullName),
		Tags:      makeTagsSlice(rq.PostFormValue("tags")),
		Comment:   rq.PostFormValue("comment"),
		Author:    rq.PostFormValue("author"),
		Time:      int(time.Now().Unix()),
		TextMime:  rq.PostFormValue("text_mime"),
		// Fields left: BinaryMime, BinaryPath, BinaryName, TextName, TextPath
	}
	rev.desiredTextFilename() // TextName is set now
	rev.TextPath = filepath.Join(h.Path, rev.TextName)
	return rev
}

// writeTextFileFromHttpData tries to fetch text content from `rq` for revision `rev` and write it to a corresponding text file. It used in `HandlerUpdate`.
func writeTextFileFromHttpData(rev *Revision, rq *http.Request) error {
	data := []byte(rq.PostFormValue("text"))
	err := ioutil.WriteFile(rev.TextPath, data, 0644)
	if err != nil {
		log.Println("Failed to write", len(data), "bytes to", rev.TextPath)
	}
	return err
}

// writeBinaryFileFromHttpData tries to fetch binary content from `rq` for revision `newRev` and write it to a corresponding binary file. If there is no content, it is taken from `oldRev`.
func writeBinaryFileFromHttpData(h *Hypha, oldRev Revision, newRev *Revision, rq *http.Request) error {
	// 10 MB file size limit
	rq.ParseMultipartForm(10 << 20)
	// Read file
	file, handler, err := rq.FormFile("binary")
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		log.Println("No binary data passed for", newRev.FullName)
		newRev.BinaryMime = oldRev.BinaryMime
		newRev.BinaryPath = oldRev.BinaryPath
		newRev.BinaryName = oldRev.BinaryName
		log.Println("Set previous revision's binary data")
		return nil
	}
	newRev.BinaryMime = handler.Header.Get("Content-Type")
	newRev.BinaryPath = filepath.Join(h.Path, newRev.IdAsStr()+".bin")
	newRev.BinaryName = newRev.desiredBinaryFilename()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Got", len(data), "of binary data for", newRev.FullName)
	err = ioutil.WriteFile(newRev.BinaryPath, data, 0644)
	if err != nil {
		log.Println("Failed to write", len(data), "bytes to", newRev.TextPath)
		return err
	}
	log.Println("Written", len(data), "of binary data for", newRev.FullName)
	return nil
}

func HandlerUpdate(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	log.Println("Attempt to update hypha", mux.Vars(rq)["hypha"])
	h, isNew := getHypha(vars["hypha"])
	oldRev := h.GetNewestRevision()
	newRev := revisionFromHttpData(h, rq)

	if isNew {
		h.CreateDir()
	}
	err := writeTextFileFromHttpData(newRev, rq)
	if err != nil {
		log.Println(err)
		return
	}
	err = writeBinaryFileFromHttpData(h, oldRev, newRev, rq)
	if err != nil {
		log.Println(err)
		return
	}

	h.Revisions[newRev.IdAsStr()] = newRev
	h.SaveJson()

	log.Println("Current hyphae storage is", hyphae)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	d := map[string]string{"Name": h.FullName}
	w.Write([]byte(renderFromMap(d, "updateOk.html")))
}
