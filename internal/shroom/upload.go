package shroom

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/files"
	hyphae2 "github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/mimetype"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func historyMessageForTextUpload(h hyphae2.Hypha, userMessage string) string {
	var verb string
	switch h.(type) {
	case *hyphae2.EmptyHypha:
		verb = "Create"
	default:
		verb = "Edit"
	}

	if userMessage == "" {
		return fmt.Sprintf("%s ‘%s’", verb, h.CanonicalName())
	}
	return fmt.Sprintf("%s ‘%s’: %s", verb, h.CanonicalName(), userMessage)
}

func writeTextToDisk(h hyphae2.ExistingHypha, data []byte, hop *history.Op) error {
	if err := hyphae2.WriteToMycoFile(h, data); err != nil {
		return err
	}
	hop.WithFiles(h.TextFilePath())

	return nil
}

// UploadText edits the hypha's text part and makes a history record about that.
func UploadText(h hyphae2.Hypha, data []byte, userMessage string, u *user.User) error {
	hop := history.
		Operation(history.TypeEditText).
		WithMsg(historyMessageForTextUpload(h, userMessage)).
		WithUser(u)

	// Privilege check
	if !u.CanProceed("upload-text") {
		rejectEditLog(h, u, "no rights")
		hop.Abort()
		return errors.New("ui.act_no_rights")
	}

	// Hypha name exploit check
	if !hyphae2.IsValidName(h.CanonicalName()) {
		// We check for the name only. I suppose the filepath would be valid as well.
		hop.Abort()
		return errors.New("invalid hypha name")
	}

	oldText, err := hyphae2.FetchMycomarkupFile(h)
	if err != nil {
		hop.Abort()
		return err
	}

	// Empty data check
	if len(bytes.TrimSpace(data)) == 0 && len(oldText) == 0 { // if nothing but whitespace
		hop.Abort()
		return nil
	}

	// At this point, we have a savable user-generated Mycomarkup document. Gotta save it.

	switch h := h.(type) {
	case *hyphae2.EmptyHypha:
		H := hyphae2.ExtendEmptyToTextual(h, filepath.Join(files.HyphaeDir(), h.CanonicalName()+".myco"))

		err := writeTextToDisk(H, data, hop)
		if err != nil {
			hop.Abort()
			return err
		}

		hyphae2.Insert(H)
		backlinks.UpdateBacklinksAfterEdit(H, "")
	case *hyphae2.MediaHypha:
		// TODO: that []byte(...) part should be removed
		if bytes.Equal(data, []byte(oldText)) {
			// No changes! Just like cancel button
			hop.Abort()
			return nil
		}

		err = writeTextToDisk(h, data, hop)
		if err != nil {
			hop.Abort()
			return err
		}

		backlinks.UpdateBacklinksAfterEdit(h, oldText)
	case *hyphae2.TextualHypha:
		oldText, err := hyphae2.FetchMycomarkupFile(h)
		if err != nil {
			hop.Abort()
			return err
		}

		// TODO: that []byte(...) part should be removed
		if bytes.Equal(data, []byte(oldText)) {
			// No changes! Just like cancel button
			hop.Abort()
			return nil
		}

		err = writeTextToDisk(h, data, hop)
		if err != nil {
			hop.Abort()
			return err
		}

		backlinks.UpdateBacklinksAfterEdit(h, oldText)
	}

	hop.Apply()
	return nil
}

func historyMessageForMediaUpload(h hyphae2.Hypha, mime string) string {
	return fmt.Sprintf("Upload media for ‘%s’ with type ‘%s’", h.CanonicalName(), mime)
}

// writeMediaToDisk saves the given data with the given mime type for the given hypha to the disk and returns the path to the saved file and an error, if any.
func writeMediaToDisk(h hyphae2.Hypha, mime string, data []byte) (string, error) {
	var (
		ext = mimetype.ToExtension(mime)
		// That's where the file will go
		uploadedFilePath = filepath.Join(files.HyphaeDir(), h.CanonicalName()+ext)
	)

	if err := os.MkdirAll(filepath.Dir(uploadedFilePath), 0777); err != nil {
		return uploadedFilePath, err
	}

	if err := os.WriteFile(uploadedFilePath, data, 0666); err != nil {
		return uploadedFilePath, err
	}
	return uploadedFilePath, nil
}

// UploadBinary edits the hypha's media part and makes a history record about that.
func UploadBinary(h hyphae2.Hypha, mime string, file multipart.File, u *user.User) error {

	// Privilege check
	if !u.CanProceed("upload-binary") {
		rejectUploadMediaLog(h, u, "no rights")
		return errors.New("ui.act_no_rights")
	}

	// Hypha name exploit check
	if !hyphae2.IsValidName(h.CanonicalName()) {
		// We check for the name only. I suppose the filepath would be valid as well.
		return errors.New("invalid hypha name")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Empty data check
	if len(data) == 0 {
		return errors.New("No data passed")
	}

	// At this point, we have a savable media document. Gotta save it.

	uploadedFilePath, err := writeMediaToDisk(h, mime, data)
	if err != nil {
		return err
	}

	switch h := h.(type) {
	case *hyphae2.EmptyHypha:
		H := hyphae2.ExtendEmptyToMedia(h, uploadedFilePath)
		hyphae2.Insert(H)
	case *hyphae2.TextualHypha:
		hyphae2.Insert(hyphae2.ExtendTextualToMedia(h, uploadedFilePath))
	case *hyphae2.MediaHypha: // If this is not the first media the hypha gets
		prevFilePath := h.MediaFilePath()
		if prevFilePath != uploadedFilePath {
			if err := history.Rename(prevFilePath, uploadedFilePath); err != nil {
				return err
			}
			log.Printf("Move ‘%s’ to ‘%s’\n", prevFilePath, uploadedFilePath)
			h.SetMediaFilePath(uploadedFilePath)
		}
	}

	history.
		Operation(history.TypeEditBinary).
		WithMsg(historyMessageForMediaUpload(h, mime)).
		WithUser(u).
		WithFiles(uploadedFilePath).
		Apply()
	return nil
}
