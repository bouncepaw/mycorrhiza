package main

import (
	"fmt"
)

type Hypha struct {
	// Hypha is physically located here. Most fields below are stored in <path>/mm.ini (mm for metametadata). Its revisions are physically located in <path>/<n>/ subfolders. <n> ∈ [0;∞] with 0 being latest revision, 1 the first.
	Path string
	// Every hypha was created at some point
	CreationTime int `json:"creationTime"`
	// Hypha has name but it can be changed
	Name string `json:"name"`
	// Hypha can be deleted. If it is deleted, it is not indexed by most of the software but still can be recovered at some point.
	Deleted bool `json:"deleted"`
	// Fields below are not part of m.ini and are created when traversing the file tree.
	// Hypha can be a child of any other hypha except its children. The parent hypha is stored in <path>/..
	ParentName string
	// Hypha can have any number of children which are stored as subfolders in <path>.
	ChildrenNames []string
	Revisions     []Revision
}

func (h Hypha) String() string {
	var revbuf string
	for _, r := range h.Revisions {
		revbuf += r.String() + "\n"
	}
	return fmt.Sprintf("Hypha %v {\n\t"+
		"path %v\n\t"+
		"created at %v\n\t"+
		"child of %v\n\t"+
		"parent of %v\n\t"+
		"Having these revisions:\n%v"+
		"}\n", h.Name, h.Path, h.CreationTime, h.ParentName, h.ChildrenNames,
		revbuf)
}

type Revision struct {
	// Revision is hypha's state at some point in time. Future revisions are not really supported. Most data here is stored in m.ini.
	Id int
	// Name used at this revision
	Name string `json:"name"`
	// Present in every hypha. Stored in t.txt.
	TextPath string
	// In at least one markup. Supported ones are "myco", "html", "md", "plain"
	Markup string `json:"markup"`
	// Some hyphæ have binary contents such as images. Their presence change hypha's behavior in a lot of ways (see methods' implementations). If stored, it is stored in b (filename "b")
	BinaryPath string
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
	// Rest of fields are ignored
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
