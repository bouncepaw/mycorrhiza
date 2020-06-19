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

type Hypha struct {
	FullName      string               `json:"-"`
	Path          string               `json:"-"`
	ViewCount     int                  `json:"views"`
	Deleted       bool                 `json:"deleted"`
	Revisions     map[string]*Revision `json:"revisions"`
	ChildrenNames []string             `json:"-"`
	parentName    string
}

func (h *Hypha) AddChild(childName string) {
	h.ChildrenNames = append(h.ChildrenNames, childName)
}

// Used with action=zen|view
func (h *Hypha) AsHtml(hyphae map[string]*Hypha, rev string) (string, error) {
	if "0" == rev {
		rev = h.NewestRevision()
	}
	r, ok := h.Revisions[rev]
	if !ok {
		return "", fmt.Errorf("Hypha %v has no such revision: %v", h.FullName, rev)
	}
	html, err := r.AsHtml(hyphae)
	return html, err
}

func (h *Hypha) Name() string {
	return h.FullName
}

func (h *Hypha) GetNewestRevision() Revision {
	return *h.Revisions[h.NewestRevision()]
}

func (h *Hypha) NewestRevision() string {
	return strconv.Itoa(h.NewestRevisionInt())
}

func (h *Hypha) NewestRevisionInt() int {
	var largest int
	for k, _ := range h.Revisions {
		rev, _ := strconv.Atoi(k)
		if rev > largest {
			largest = rev
		}
	}
	return largest
}

func (h *Hypha) MetaJsonPath() string {
	return filepath.Join(h.Path, "meta.json")
}

func (h *Hypha) CreateDir() error {
	return os.MkdirAll(h.Path, 0644)
}

func (h *Hypha) ParentName() string {
	return h.parentName
}

func (h *Hypha) SaveJson() {
	data, err := json.Marshal(h)
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

func ActionEdit(hyphaName string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var initContents, initTextMime, initTags string
	hypha, ok := hyphae[hyphaName]
	if !ok {
		initContents = "Describe " + hyphaName + "here."
		initTextMime = "text/markdown"
	} else {
		newestRev := hypha.Revisions[hypha.NewestRevision()]
		contents, err := ioutil.ReadFile(newestRev.TextPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("<b>Sorry, something went wrong</b>"))
			return
		}
		initContents = string(contents)
		initTextMime = newestRev.TextMime
		initTags = strings.Join(newestRev.Tags, ",")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EditHyphaPage(hyphaName, initTextMime, initContents, initTags)))
}
