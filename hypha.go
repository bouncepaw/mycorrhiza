package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func GetRevision(hyphae map[string]*Hypha, hyphaName string, rev string, w http.ResponseWriter) (Revision, bool) {
	for name, _ := range hyphae {
		if name == hyphaName {
			for _, r := range hyphae[name].Revisions {
				id, err := strconv.Atoi(rev)
				if err != nil {
					log.Println("No such revision", rev, "at hypha", hyphaName)
					w.WriteHeader(http.StatusNotFound)
					return Revision{}, false
				}
				if r.Id == id {
					return r, true
				}
			}
		}
	}
	return Revision{}, false
}

// `rev` is the id of revision to render. If it = 0, the last one is rendered. If the revision is not found, an error is returned.
func (h Hypha) Render(hyphae map[string]*Hypha, rev int) (ret string, err error) {
	for _, r := range h.Revisions {
		if r.Id == rev {
			return r.Render(hyphae)
		}
	}
	return "", errors.New("Revision was not found")
}
