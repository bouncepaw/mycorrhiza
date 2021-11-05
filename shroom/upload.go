package shroom

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/user"
)

// UploadText edits a hypha' text part and makes a history record about that.
func UploadText(h *hyphae.Hypha, data []byte, message string, u *user.User, lc *l18n.Localizer) (hop *history.Op, errtitle string) {
	hop = history.Operation(history.TypeEditText)
	var action string
	if h.Exists {
		action = "Edit"
	} else {
		action = "Create"
	}

	if message == "" {
		hop.WithMsg(fmt.Sprintf("%s ‘%s’", action, h.Name))
	} else {
		hop.WithMsg(fmt.Sprintf("%s ‘%s’: %s", action, h.Name, message))
	}

	if errtitle, err := CanEdit(u, h, lc); err != nil {
		return hop.WithErrAbort(err), errtitle
	}
	if len(bytes.TrimSpace(data)) == 0 && h.BinaryPath == "" {
		return hop.WithErrAbort(errors.New("No data passed")), "Empty"
	}

	return uploadHelp(h, hop, ".myco", data, u)
}

// UploadBinary edits a hypha' attachment and makes a history record about that.
func UploadBinary(h *hyphae.Hypha, mime string, file multipart.File, u *user.User, lc *l18n.Localizer) (*history.Op, string) {
	var (
		hop       = history.Operation(history.TypeEditBinary).WithMsg(fmt.Sprintf("Upload attachment for ‘%s’ with type ‘%s’", h.Name, mime))
		data, err = io.ReadAll(file)
	)

	if err != nil {
		return hop.WithErrAbort(err), err.Error()
	}
	if errtitle, err := CanAttach(u, h, lc); err != nil {
		return hop.WithErrAbort(err), errtitle
	}
	if len(data) == 0 {
		return hop.WithErrAbort(errors.New("No data passed")), "Empty"
	}

	return uploadHelp(h, hop, mimetype.ToExtension(mime), data, u)
}

// uploadHelp is a helper function for UploadText and UploadBinary
func uploadHelp(h *hyphae.Hypha, hop *history.Op, ext string, data []byte, u *user.User) (*history.Op, string) {
	var (
		fullPath         = filepath.Join(files.HyphaeDir(), h.Name+ext)
		originalFullPath = &h.TextPath
		originalText     = "" // for backlink update
	)
	if !isValidPath(fullPath) || !hyphae.IsValidName(h.Name) {
		err := errors.New("bad path")
		return hop.WithErrAbort(err), err.Error()
	}
	if hop.Type == history.TypeEditBinary {
		originalFullPath = &h.BinaryPath
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		return hop.WithErrAbort(err), err.Error()
	}

	if hop.Type == history.TypeEditText {
		originalText, _ = FetchTextPart(h)
	}

	if err := os.WriteFile(fullPath, data, 0666); err != nil {
		return hop.WithErrAbort(err), err.Error()
	}

	if h.Exists && *originalFullPath != fullPath && *originalFullPath != "" {
		if err := history.Rename(*originalFullPath, fullPath); err != nil {
			return hop.WithErrAbort(err), err.Error()
		}
		log.Println("Move", *originalFullPath, "to", fullPath)
	}

	h.InsertIfNew()
	if h.Exists && h.TextPath != "" && hop.Type == history.TypeEditText && !history.FileChanged(fullPath) {
		return hop.Abort(), "No changes"
	}
	*originalFullPath = fullPath
	if hop.Type == history.TypeEditText {
		hyphae.BacklinksOnEdit(h, originalText)
	}
	return hop.WithFiles(fullPath).WithUser(u).Apply(), ""
}

func isValidPath(pathname string) bool {
	return strings.HasPrefix(pathname, files.HyphaeDir())
}
