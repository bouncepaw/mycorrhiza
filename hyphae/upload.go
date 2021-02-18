package hyphae

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func rejectEditLog(h *Hypha, u *user.User, errmsg string) {
	log.Printf("Reject edit ‘%s’ by @%s: %s\n", h.Name, u.Name, errmsg)
}

func rejectAttachLog(h *Hypha, u *user.User, errmsg string) {
	log.Printf("Reject attach ‘%s’ by @%s: %s\n", h.Name, u.Name, errmsg)
}

func (h *Hypha) CanEdit(u *user.User) (err error, errtitle string) {
	if !u.CanProceed("edit") {
		rejectEditLog(h, u, "no rights")
		return errors.New("You must be an editor to edit pages."), "Not enough rights"
	}
	return nil, ""
}

func (h *Hypha) CanUploadThat(data []byte, u *user.User) (err error, errtitle string) {
	if len(data) == 0 {
		return errors.New("No text data passed"), "Empty"
	}
	return nil, ""
}

func (h *Hypha) UploadText(textData []byte, u *user.User) (hop *history.HistoryOp, errtitle string) {
	hop = history.Operation(history.TypeEditText)
	if h.Exists {
		hop.WithMsg(fmt.Sprintf("Edit ‘%s’", h.Name))
	} else {
		hop.WithMsg(fmt.Sprintf("Create ‘%s’", h.Name))
	}

	if err, errtitle := h.CanEdit(u); err != nil {
		return hop.WithError(err), errtitle
	}
	if err, errtitle := h.CanUploadThat(textData, u); err != nil {
		return hop.WithError(err), errtitle
	}

	return h.uploadHelp(hop, ".myco", textData, u)
}

func (h *Hypha) CanAttach(err error, u *user.User) (error, string) {
	if !u.CanProceed("upload-binary") {
		rejectAttachLog(h, u, "no rights")
		return errors.New("You must be an editor to upload attachments."), "Not enough rights"
	}

	if err != nil {
		rejectAttachLog(h, u, err.Error())
		return errors.New("No binary data passed"), err.Error()
	}
	return nil, ""
}

func (h *Hypha) UploadBinary(mime string, file multipart.File, u *user.User) (*history.HistoryOp, string) {
	var (
		hop       = history.Operation(history.TypeEditBinary).WithMsg(fmt.Sprintf("Upload binary part for ‘%s’ with type ‘%s’", h.Name, mime))
		data, err = ioutil.ReadAll(file)
	)

	if err != nil {
		return hop.WithError(err), err.Error()
	}
	if err, errtitle := h.CanEdit(u); err != nil {
		return hop.WithError(err), errtitle
	}
	if err, errtitle := h.CanUploadThat(data, u); err != nil {
		return hop.WithError(err), errtitle
	}

	return h.uploadHelp(hop, mimetype.ToExtension(mime), data, u)
}

// uploadHelp is a helper function for UploadText and UploadBinary
func (h *Hypha) uploadHelp(hop *history.HistoryOp, ext string, data []byte, u *user.User) (*history.HistoryOp, string) {
	var (
		fullPath         = filepath.Join(util.WikiDir, h.Name+ext)
		originalFullPath = &h.TextPath
	)
	if hop.Type == history.TypeEditBinary {
		originalFullPath = &h.BinaryPath
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		return hop.WithError(err), err.Error()
	}

	if err := ioutil.WriteFile(fullPath, data, 0644); err != nil {
		return hop.WithError(err), err.Error()
	}

	if h.Exists && *originalFullPath != fullPath && *originalFullPath != "" {
		if err := history.Rename(*originalFullPath, fullPath); err != nil {
			return hop.WithError(err), err.Error()
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
