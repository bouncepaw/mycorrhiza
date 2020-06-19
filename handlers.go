package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Boilerplate code present in many handlers. Good to have it.
func HandlerBase(w http.ResponseWriter, r *http.Request) (Revision, bool) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	return GetRevision(hyphae, vars["hypha"], revno, w)
}

func HandlerGetBinary(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionGetBinary(w)
	}
}

func HandlerRaw(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionRaw(w)
	}
}

func HandlerZen(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionZen(w)
	}
}

func HandlerView(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionView(w, HyphaPage)
	}
}

func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ActionEdit(vars["hypha"], w)
}

func HandlerRewind(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerRename(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func makeTagsSlice(responseTagsString string) (ret []string) {
	// `responseTagsString` is string like "foo,, bar,kek". Whitespace around commas is insignificant. Expected output: []string{"foo", "bar", "kek"}
	for _, tag := range strings.Split(responseTagsString, ",") {
		if trimmed := strings.TrimSpace(tag); "" == trimmed {
			ret = append(ret, trimmed)
		}
	}
	return ret
}

// Return an existing hypha it exists in `hyphae` or create a new one. If it `isNew`, you'll have to insert it to `hyphae` yourself.
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

// Create a new revison for hypha `h`. All data is fetched from `r`, except for BinaryMime and BinaryPath which require additional processing. You'll have te insert the revision to `h` yourself.
func revisionFromHttpData(h *Hypha, r *http.Request) *Revision {
	idStr := strconv.Itoa(h.NewestRevisionInt() + 1)
	log.Println(idStr)
	rev := &Revision{
		Id:       h.NewestRevisionInt() + 1,
		FullName: h.FullName,
		Tags:     makeTagsSlice(r.PostFormValue("tags")),
		Comment:  r.PostFormValue("comment"),
		Author:   r.PostFormValue("author"),
		Time:     int(time.Now().Unix()),
		TextMime: r.PostFormValue("text_mime"),
		TextPath: filepath.Join(h.Path, idStr+".txt"),
		// Left: BinaryMime, BinaryPath
	}
	return rev
}

func writeTextFileFromHttpData(rev *Revision, r *http.Request) error {
	data := []byte(r.PostFormValue("text"))
	err := ioutil.WriteFile(rev.TextPath, data, 0644)
	if err != nil {
		log.Println("Failed to write", len(data), "bytes to", rev.TextPath)
	}
	return err
}

func writeBinaryFileFromHttpData(h *Hypha, oldRev Revision, newRev *Revision, r *http.Request) error {
	// 10 MB file size limit
	r.ParseMultipartForm(10 << 20)
	// Read file
	file, handler, err := r.FormFile("binary")
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		log.Println("No binary data passed for", newRev.FullName)
		newRev.BinaryMime = oldRev.BinaryMime
		newRev.BinaryPath = oldRev.BinaryPath
		log.Println("Set previous revision's binary data")
		return nil
	}
	newRev.BinaryMime = handler.Header.Get("Content-Type")
	newRev.BinaryPath = filepath.Join(h.Path, newRev.IdAsStr()+".bin")
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

func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("Attempt to update hypha", mux.Vars(r)["hypha"])
	h, isNew := getHypha(vars["hypha"])
	oldRev := h.GetNewestRevision()
	newRev := revisionFromHttpData(h, r)

	if isNew {
		h.CreateDir()
	}
	err := writeTextFileFromHttpData(newRev, r)
	if err != nil {
		log.Println(err)
		return
	}
	err = writeBinaryFileFromHttpData(h, oldRev, newRev, r)
	if err != nil {
		log.Println(err)
		return
	}

	h.Revisions[newRev.IdAsStr()] = newRev
	h.SaveJson()

	log.Println("Current hyphae storage is", hyphae)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Saved successfully"))
}
