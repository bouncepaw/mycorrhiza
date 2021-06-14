package shroom

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/user"
)

func UploadText(h *hyphae.Hypha, data []byte, message string, u *user.User) (hop *history.HistoryOp, errtitle string) {
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


	if err, errtitle := CanEdit(u, h); err != nil {
		return hop.WithErrAbort(err), errtitle
	}
	if len(data) == 0 {
		return hop.WithErrAbort(errors.New("No data passed")), "Empty"
	}

	return uploadHelp(h, hop, ".myco", data, u)
}

func UploadBinary(h *hyphae.Hypha, mime string, file multipart.File, u *user.User) (*history.HistoryOp, string) {
	var (
		hop       = history.Operation(history.TypeEditBinary).WithMsg(fmt.Sprintf("Upload binary part for ‘%s’ with type ‘%s’", h.Name, mime))
		data, err = ioutil.ReadAll(file)
	)

	if err != nil {
		return hop.WithErrAbort(err), err.Error()
	}
	if err, errtitle := CanAttach(u, h); err != nil {
		return hop.WithErrAbort(err), errtitle
	}
	if len(data) == 0 {
		return hop.WithErrAbort(errors.New("No data passed")), "Empty"
	}

	return uploadHelp(h, hop, mimetype.ToExtension(mime), data, u)
}

// uploadHelp is a helper function for UploadText and UploadBinary
func uploadHelp(h *hyphae.Hypha, hop *history.HistoryOp, ext string, data []byte, u *user.User) (*history.HistoryOp, string) {
	var (
		fullPath         = filepath.Join(cfg.WikiDir, h.Name+ext)
		originalFullPath = &h.TextPath
	)
	if hop.Type == history.TypeEditBinary {
		originalFullPath = &h.BinaryPath
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		return hop.WithErrAbort(err), err.Error()
	}

	if err := ioutil.WriteFile(fullPath, data, 0644); err != nil {
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
	return hop.WithFiles(fullPath).WithUser(u).Apply(), ""
}
