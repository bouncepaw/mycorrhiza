package main

import (
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
)

type Revision struct {
	// Revision is hypha's state at some point in time. Future revisions are not really supported. Most data here is stored in m.ini.
	Id int
	// Name used at this revision
	Name string `json:"name"`
	// Name of hypha
	FullName string
	// Present in every hypha. Stored in t.txt.
	TextPath string
	// In at least one markup. Supported ones are "myco", "html", "md", "plain"
	Markup string `json:"markup"`
	// Some hyphæ have binary contents such as images. Their presence change hypha's behavior in a lot of ways (see methods' implementations). If stored, it is stored in b (filename "b")
	BinaryPath    string
	BinaryRequest string
	// To tell what is meaning of binary content, mimeType for them is stored. If the hypha has no binary content, this field must be "application/x-hypha"
	MimeType string `json:"mimeType"`
	// Every revision was created at some point. This field stores the creation time of the latest revision
	RevisionTime int `json:"createdAt"`
	// Every hypha has any number of tags
	Tags []string `json:"tags"`
	// Current revision is authored by someone
	RevisionAuthor string `json:"author"`
	// and has a comment in plain text
	RevisionComment string `json:"comment"`
}

func (h Revision) String() string {
	return fmt.Sprintf(`Revision %v created at %v {
	name: %v
	textPath: %v
	markup: %v
	binaryPath: %v
	mimeType: %v
	tags: %v
	revisionAuthor: %v
	revisionComment: %v
}`, h.Id, h.RevisionTime, h.Name, h.TextPath, h.Markup, h.BinaryPath, h.MimeType, h.Tags, h.RevisionAuthor, h.RevisionComment)
}

// This method is meant to be called only by Hypha#Render.
func (r Revision) Render(hyphae map[string]*Hypha) (ret string, err error) {
	ret += `<article class="page">
`
	// If it is a binary hypha (we support only images for now):
	// TODO: support things other than images.
	if r.MimeType != "application/x-hypha" {
		ret += fmt.Sprintf(`<img src="/%s" class="page__image"/>`, r.BinaryRequest)
	}

	contents, err := ioutil.ReadFile(r.TextPath)
	if err != nil {
		return "", err
	}

	// TODO: support more markups.
	// TODO: support mycorrhiza extensions like transclusion.
	switch r.Markup {
	case "plain":
		ret += fmt.Sprintf(`<pre>%s</pre>`, contents)
	case "md":
		html := markdown.ToHTML(contents, nil, nil)
		ret += string(html)
	default:
		return "", errors.New("Unsupported markup: " + r.Markup)
	}

	ret += `
</article>`
	return ret, nil
}
