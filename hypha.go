package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// `Hypha` represents a hypha. It is the thing MycorrhizaWiki generally serves.
// Each hypha has 1 or more revisions.
type Hypha struct {
	FullName      string               `json:"-"`
	Path          string               `json:"-"`
	ViewCount     int                  `json:"views"`
	Deleted       bool                 `json:"deleted"`
	Revisions     map[string]*Revision `json:"revisions"`
	ChildrenNames []string             `json:"-"`
	parentName    string
}

// AsHtml returns HTML representation of the hypha.
// No layout or navigation are present here. Just the hypha.
func (h *Hypha) AsHtml(id string, w http.ResponseWriter) (string, error) {
	if "0" == id {
		id = h.NewestRevision()
	}
	if rev, ok := h.Revisions[id]; ok {
		return rev.AsHtml(w)
	}
	return "", fmt.Errorf("Hypha %v has no such revision: %v", h.FullName, id)
}

// GetNewestRevision returns the most recent Revision.
func (h *Hypha) GetNewestRevision() Revision {
	return *h.Revisions[h.NewestRevision()]
}

// NewestRevision returns the most recent revision's id as a string.
func (h *Hypha) NewestRevision() string {
	return strconv.Itoa(h.NewestRevisionInt())
}

// NewestRevision returns the most recent revision's id as an integer.
func (h *Hypha) NewestRevisionInt() (ret int) {
	for k, _ := range h.Revisions {
		id, _ := strconv.Atoi(k)
		if id > ret {
			ret = id
		}
	}
	return ret
}

// MetaJsonPath returns rooted path to the hypha's `meta.json` file.
// It is not promised that the file exists.
func (h *Hypha) MetaJsonPath() string {
	return filepath.Join(h.Path, "meta.json")
}

// CreateDir creates directory where the hypha must reside.
// It is meant to be used with new hyphae.
func (h *Hypha) CreateDir() error {
	return os.MkdirAll(h.Path, os.ModePerm)
}

// SaveJson dumps the hypha's metadata to `meta.json` file.
func (h *Hypha) SaveJson() {
	data, err := json.MarshalIndent(h, "", "\t")
	if err != nil {
		log.Println("Failed to create JSON of hypha.", err)
		return
	}
	err = ioutil.WriteFile(h.MetaJsonPath(), data, 0644)
	if err != nil {
		log.Println("Failed to save JSON of hypha.", err)
		return
	}
	log.Println("Saved JSON data of", h.FullName)
}

// ActionEdit is called with `?acton=edit`.
// It represents the hypha editor.
func ActionEdit(hyphaName string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var initContents, initTextMime, initTags string
	if h, ok := hyphae[hyphaName]; ok {
		newestRev := h.GetNewestRevision()
		contents, err := ioutil.ReadFile(newestRev.TextPath)
		if err != nil {
			log.Println("Could not read", newestRev.TextPath)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(GenericErrorMsg))
			return
		}
		initContents = string(contents)
		initTextMime = newestRev.TextMime
		initTags = strings.Join(newestRev.Tags, ",")
	} else {
		initContents = "Describe " + hyphaName + "here."
		initTextMime = "text/markdown"
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EditHyphaPage(hyphaName, initTextMime, initContents, initTags)))
}
