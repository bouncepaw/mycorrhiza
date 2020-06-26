package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

type Hypha struct {
	Exists    bool                 `json:"-"`
	FullName  string               `json:"-"`
	ViewCount int                  `json:"views"`
	Deleted   bool                 `json:"deleted"`
	Revisions map[string]*Revision `json:"revisions"`
	actual    *Revision            `json:"-"`
	Invalid   bool
	Err       error
}

func (h *Hypha) Invalidate(err error) *Hypha {
	h.Invalid = true
	h.Err = err
	return h
}

func (h *Hypha) MetaJsonPath() string {
	return filepath.Join(h.Path(), "meta.json")
}

func (h *Hypha) Path() string {
	return filepath.Join(cfg.WikiDir, h.FullName)
}

func (h *Hypha) TextPath() string {
	return h.actual.TextPath
}

func (h *Hypha) TagsJoined() string {
	if h.Exists {
		return strings.Join(h.actual.Tags, ", ")
	}
	return ""
}

func (h *Hypha) TextMime() string {
	if h.Exists {
		return h.actual.TextMime
	}
	return "text/markdown"
}

func (h *Hypha) TextContent() string {
	if h.Exists {
		contents, err := ioutil.ReadFile(h.TextPath())
		if err != nil {
			log.Println("Could not read", h.FullName)
			return "Error: could not hypha text content file. It is recommended to cancel editing. Please contact the wiki admin. If you are the admin, see the logs."
		}
		return string(contents)
	}
	return fmt.Sprintf(cfg.DescribeHyphaHerePattern, h.FullName)
}

func (s *Storage) Open(name string) *Hypha {
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
			return h.Invalidate(err)
		}

		err = json.Unmarshal(metaJsonText, &h)
		if err != nil {
			return h.Invalidate(err)
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

		return h.OnRevision("0")
	}
	return h
}

func (h *Hypha) parentName() string {
	return filepath.Dir(h.FullName)
}

func (h *Hypha) metaJsonPath() string {
	return filepath.Join(cfg.WikiDir, h.FullName, "meta.json")
}

// OnRevision tries to change to a revision specified by `id`.
func (h *Hypha) OnRevision(id string) *Hypha {
	if h.Invalid || !h.Exists {
		return h
	}
	if len(h.Revisions) == 0 {
		return h.Invalidate(errors.New("This hypha has no revisions"))
	}
	if id == "0" {
		id = h.NewestId()
	}
	// Revision must be there, so no error checking
	if rev, _ := h.Revisions[id]; true {
		h.actual = rev
	}
	return h
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
	if h.Exists {
		log.Println(h.FullName, h.actual.Id, s)
	} else {
		log.Println("nonexistent", h.FullName, s)
	}
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

// ActionBinary is used with `?action=binary`.
// It writes contents of binary content file.
func (h *Hypha) ActionBinary(w http.ResponseWriter) {
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

// ActionView is used with `?action=view` or no action at all.
// It renders the page, the layout and everything else.
func (h *Hypha) ActionView(w http.ResponseWriter, renderExists, renderNotExists func(string, string) string) {
	var html string
	var err error
	if h.Exists {
		html, err = h.asHtml()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if h.Exists {
		w.Write([]byte(renderExists(h.FullName, html)))
	} else {
		w.Write([]byte(renderNotExists(h.FullName, "")))
	}
	h.PlainLog("Rendering hypha view")
}

// CreateDirIfNeeded creates directory where the hypha must reside if needed.
// It is not needed if the dir already exists.
func (h *Hypha) CreateDirIfNeeded() *Hypha {
	if h.Invalid {
		return h
	}
	// os.MkdirAll created dir if it is not there. Basically, checks it for us.
	err := os.MkdirAll(filepath.Join(cfg.WikiDir, h.FullName), os.ModePerm)
	if err != nil {
		h.Invalidate(err)
	}
	return h
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

// revisionFromHttpData creates a new revison for hypha `h`. All data is fetched from `rq`, except for BinaryMime and BinaryPath which require additional processing. The revision is inserted for you. You'll have to pop it out if there is an error.
func (h *Hypha) AddRevisionFromHttpData(rq *http.Request) *Hypha {
	if h.Invalid {
		return h
	}
	id := 1
	if h.Exists {
		id = h.actual.Id + 1
	}
	log.Printf("Creating revision %d from http data", id)
	rev := &Revision{
		Id:        id,
		FullName:  h.FullName,
		ShortName: filepath.Base(h.FullName),
		Tags:      makeTagsSlice(rq.PostFormValue("tags")),
		Comment:   rq.PostFormValue("comment"),
		Author:    rq.PostFormValue("author"),
		Time:      int(time.Now().Unix()),
		TextMime:  rq.PostFormValue("text_mime"),
		// Fields left: BinaryMime, BinaryPath, BinaryName, TextName, TextPath
	}
	rev.generateTextFilename() // TextName is set now
	rev.TextPath = filepath.Join(h.Path(), rev.TextName)
	return h.AddRevision(rev)
}

func (h *Hypha) AddRevision(rev *Revision) *Hypha {
	if h.Invalid {
		return h
	}
	h.Revisions[strconv.Itoa(rev.Id)] = rev
	h.actual = rev
	return h
}

// WriteTextFileFromHttpData tries to fetch text content from `rq` for revision `rev` and write it to a corresponding text file. It used in `HandlerUpdate`.
func (h *Hypha) WriteTextFileFromHttpData(rq *http.Request) *Hypha {
	if h.Invalid {
		return h
	}
	data := []byte(rq.PostFormValue("text"))
	err := ioutil.WriteFile(h.TextPath(), data, 0644)
	if err != nil {
		log.Println("Failed to write", len(data), "bytes to", h.TextPath())
		h.Invalidate(err)
	}
	return h
}

// WriteBinaryFileFromHttpData tries to fetch binary content from `rq` for revision `newRev` and write it to a corresponding binary file. If there is no content, it is taken from a previous revision, if there is any.
func (h *Hypha) WriteBinaryFileFromHttpData(rq *http.Request) *Hypha {
	if h.Invalid {
		return h
	}
	// 10 MB file size limit
	rq.ParseMultipartForm(10 << 20)
	// Read file
	file, handler, err := rq.FormFile("binary")
	if file != nil {
		defer file.Close()
	}
	// If file is not passed:
	if err != nil {
		// Let's hope there are no other errors 🙏
		// TODO: actually check if there any other errors
		log.Println("No binary data passed for", h.FullName)
		// It is expected there is at least one revision
		if len(h.Revisions) > 1 {
			prevRev := h.Revisions[strconv.Itoa(h.actual.Id-1)]
			h.actual.BinaryMime = prevRev.BinaryMime
			h.actual.BinaryPath = prevRev.BinaryPath
			h.actual.BinaryName = prevRev.BinaryName
			log.Println("Set previous revision's binary data")
		}
		return h
	}
	// If file is passed:
	h.actual.BinaryMime = handler.Header.Get("Content-Type")
	h.actual.generateBinaryFilename()
	h.actual.BinaryPath = filepath.Join(h.Path(), h.actual.BinaryName)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return h.Invalidate(err)
	}
	log.Println("Got", len(data), "of binary data for", h.FullName)
	err = ioutil.WriteFile(h.actual.BinaryPath, data, 0644)
	if err != nil {
		return h.Invalidate(err)
	}
	log.Println("Written", len(data), "of binary data for", h.FullName)
	return h
}

// SaveJson dumps the hypha's metadata to `meta.json` file.
func (h *Hypha) SaveJson() *Hypha {
	if h.Invalid {
		return h
	}
	data, err := json.MarshalIndent(h, "", "\t")
	if err != nil {
		return h.Invalidate(err)
	}
	err = ioutil.WriteFile(h.MetaJsonPath(), data, 0644)
	if err != nil {
		return h.Invalidate(err)
	}
	log.Println("Saved JSON data of", h.FullName)
	return h
}

// Store adds `h` to the `Hs` if it is not already there
func (h *Hypha) Store() *Hypha {
	if !h.Invalid {
		Hs.paths[h.FullName] = h.Path()
	}
	return h
}
