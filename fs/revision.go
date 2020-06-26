package fs

import (
	"mime"
	"strconv"
)

type Revision struct {
	Id         int      `json:"-"`
	FullName   string   `json:"-"`
	Tags       []string `json:"tags"`
	ShortName  string   `json:"name"`
	Comment    string   `json:"comment"`
	Author     string   `json:"author"`
	Time       int      `json:"time"`
	TextMime   string   `json:"text_mime"`
	BinaryMime string   `json:"binary_mime"`
	TextPath   string   `json:"-"`
	BinaryPath string   `json:"-"`
	TextName   string   `json:"text_name"`
	BinaryName string   `json:"binary_name"`
}

// TODO: https://github.com/bouncepaw/mycorrhiza/issues/4
// Some filenames are wrong?
func (rev *Revision) generateTextFilename() {
	ts, err := mime.ExtensionsByType(rev.TextMime)
	if err != nil || ts == nil {
		rev.TextName = strconv.Itoa(rev.Id) + ".txt"
	} else {
		rev.TextName = strconv.Itoa(rev.Id) + ts[0]
	}
}

func (rev *Revision) generateBinaryFilename() {
	ts, err := mime.ExtensionsByType(rev.BinaryMime)
	if err != nil || ts == nil {
		rev.BinaryName = strconv.Itoa(rev.Id) + ".bin"
	} else {
		rev.BinaryName = strconv.Itoa(rev.Id) + ts[0]
	}
}
