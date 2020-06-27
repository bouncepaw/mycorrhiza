// This file contains methods for Hypha that calculate data about the hypha based on known information.
package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

func (h *Hypha) MetaJsonPath() string {
	return filepath.Join(h.Path(), "meta.json")
}

func (h *Hypha) Path() string {
	return filepath.Join(cfg.WikiDir, h.FullName)
}

func (h *Hypha) TextPath() string {
	return h.actual.TextPath
}

func (h *Hypha) parentName() string {
	return filepath.Dir(h.FullName)
}

// hasBinaryData returns true if the revision has any binary data associated.
// During initialisation, it is guaranteed that r.BinaryMime is set to "" if the revision has no binary data. (is it?)
func (h *Hypha) hasBinaryData() bool {
	return h.actual.BinaryMime != ""
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

func (h *Hypha) mimeTypeForActionRaw() string {
	// If text mime type is text/html, it is not good as it will be rendered.
	if h.actual.TextMime == "text/html" {
		return "text/plain"
	}
	return h.actual.TextMime
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
