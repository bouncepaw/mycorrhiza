package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Hypha struct {
	FullName      string
	Path          string
	ViewCount     int                  `json:"views"`
	Deleted       bool                 `json:"deleted"`
	Revisions     map[string]*Revision `json:"revisions"`
	ChildrenNames []string
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

func (h *Hypha) NewestRevision() string {
	var largest int
	for k, _ := range h.Revisions {
		rev, _ := strconv.Atoi(k)
		if rev > largest {
			largest = rev
		}
	}
	return strconv.Itoa(largest)
}

func (h *Hypha) ParentName() string {
	return h.parentName
}

func ActionEdit(hyphaName string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var initContents, initTextMime, initBinaryMime, initTags string
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
		initBinaryMime = newestRev.BinaryMime
		initTags = strings.Join(newestRev.Tags, ",")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(EditHyphaPage(hyphaName, initTextMime, initBinaryMime, initContents, initTags)))
}
