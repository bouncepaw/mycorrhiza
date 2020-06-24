package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/gomarkdown/markdown"
)

type Hypha struct {
	Exists    bool                 `json:"-"`
	FullName  string               `json:"-"`
	ViewCount int                  `json:"views"`
	Deleted   bool                 `json:"deleted"`
	Revisions map[string]*Revision `json:"revisions"`
	actual    *Revision            `json:"-"`
}

func (s *Storage) Open(name string) (*Hypha, error) {
	h := &Hypha{
		Exists:   true,
		FullName: name,
	}
	path, ok := s.paths[name]
	// This hypha does not exist yet
	if !ok {
		log.Println("Hypha", name, "does not exist")
		h.Exists = false
		h.Revisions = make(map[string]*Revision)
	} else {
		metaJsonText, err := ioutil.ReadFile(filepath.Join(path, "meta.json"))
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		err = json.Unmarshal(metaJsonText, &h)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		// fill in rooted paths to content files and full names
		for idStr, rev := range h.Revisions {
			rev.FullName = filepath.Join(h.parentName(), rev.ShortName)
			rev.Id, _ = strconv.Atoi(idStr)
			if rev.BinaryName != "" {
				rev.BinaryPath = filepath.Join(path, rev.BinaryName)
			}
			rev.TextPath = filepath.Join(path, rev.TextName)
		}

		err = h.OnRevision("0")
		return h, err
	}
	return h, nil
}

func (h *Hypha) parentName() string {
	return filepath.Dir(h.FullName)
}

func (h *Hypha) metaJsonPath() string {
	return filepath.Join(cfg.WikiDir, h.FullName, "meta.json")
}

// OnRevision tries to change to a revision specified by `id`.
func (h *Hypha) OnRevision(id string) error {
	if len(h.Revisions) == 0 {
		return errors.New("This hypha has no revisions")
	}
	if id == "0" {
		id = h.NewestId()
	}
	// Revision must be there, so no error checking
	if rev, _ := h.Revisions[id]; true {
		h.actual = rev
	}
	return nil
}

// NewestId finds the largest id among all revisions.
func (h *Hypha) NewestId() string {
	var largest int
	for k, _ := range h.Revisions {
		id, _ := strconv.Atoi(k)
		if id > largest {
			largest = id
		}
	}
	return strconv.Itoa(largest)
}

func (h *Hypha) PlainLog(s string) {
	log.Println(h.FullName, h.actual.Id, s)
}

func (h *Hypha) mimeTypeForActionRaw() string {
	// If text mime type is text/html, it is not good as it will be rendered.
	if h.actual.TextMime == "text/html" {
		return "text/plain"
	}
	return h.actual.TextMime
}

// hasBinaryData returns true if the revision has any binary data associated.
// During initialisation, it is guaranteed that r.BinaryMime is set to "" if the revision has no binary data. (is it?)
func (h *Hypha) hasBinaryData() bool {
	return h.actual.BinaryMime != ""
}

func (h *Hypha) asHtml() (string, error) {
	rev := h.actual
	ret := `<article class="page">
	<h1 class="page__title">` + rev.FullName + `</h1>
`
	// What about using <figure>?
	if h.hasBinaryData() {
		ret += fmt.Sprintf(`<img src="/%s?action=getBinary&rev=%d" cla
ss="page__amnt"/>`, rev.FullName, rev.Id)
	}

	contents, err := ioutil.ReadFile(rev.TextPath)
	if err != nil {
		log.Println("Failed to render", rev.FullName, ":", err)
		return "", err
	}

	// TODO: support more markups.
	// TODO: support mycorrhiza extensions like transclusion.
	switch rev.TextMime {
	case "text/markdown":
		html := markdown.ToHTML(contents, nil, nil)
		ret += string(html)
	default:
		ret += fmt.Sprintf(`<pre>%s</pre>`, contents)
	}

	ret += `
</article>`

	return ret, nil
}

// ActionRaw is used with `?action=raw`.
// It writes text content of the revision without any parsing or rendering.
func (h *Hypha) ActionRaw(w http.ResponseWriter) {
	fileContents, err := ioutil.ReadFile(h.actual.TextPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", h.mimeTypeForActionRaw())
	w.WriteHeader(http.StatusOK)
	w.Write(fileContents)
	h.PlainLog("Serving raw text")
}

// ActionGetBinary is used with `?action=getBinary`.
// It writes contents of binary content file.
func (h *Hypha) ActionGetBinary(w http.ResponseWriter) {
	fileContents, err := ioutil.ReadFile(h.actual.BinaryPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", h.actual.BinaryMime)
	w.WriteHeader(http.StatusOK)
	w.Write(fileContents)
	h.PlainLog("Serving raw text")
}

// ActionZen is used with `?action=zen`.
// It renders the hypha but without any layout or styles. Pure. Zen.
func (h *Hypha) ActionZen(w http.ResponseWriter) {
	html, err := h.asHtml()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
	h.PlainLog("Rendering zen")
}
