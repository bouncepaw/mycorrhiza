package shroom

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/user"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func historyMessageForTextUpload(h hyphae.Hypha, userMessage string) string {
	var verb string
	switch h.(type) {
	case *hyphae.EmptyHypha:
		verb = "Create"
	default:
		verb = "Edit"
	}

	if userMessage == "" {
		return fmt.Sprintf("%s ‘%s’", verb, h.CanonicalName())
	}
	return fmt.Sprintf("%s ‘%s’: %s", verb, h.CanonicalName(), userMessage)
}

func writeTextToDisk(h hyphae.ExistingHypha, data []byte, hop *history.Op) error {
	if err := hyphae.WriteToMycoFile(h, data); err != nil {
		return err
	}
	hop.WithFiles(h.TextFilePath())

	return nil
}

// UploadText edits the hypha's text part and makes a history record about that.
func UploadText(h hyphae.Hypha, data []byte, userMessage string, u *user.User) error {
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
	if !hyphae.IsValidName(h.CanonicalName()) {
		// We check for the name only. I suppose the filepath would be valid as well.
		hop.Abort()
		return errors.New("invalid hypha name")
	}

	// Empty data check
	if len(bytes.TrimSpace(data)) == 0 { // if nothing but whitespace
		switch h.(type) {
		case *hyphae.EmptyHypha, *hyphae.MediaHypha:
			// Writing no description, it's ok, just like cancel button.
			hop.Abort()
			return nil
		case *hyphae.TextualHypha:
			// What do you want passing nothing for a textual hypha?
			return errors.New("No data passed")
		}
	}

	// At this point, we have a savable user-generated Mycomarkup document. Gotta save it.

	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		H := hyphae.ExtendEmptyToTextual(h, filepath.Join(files.HyphaeDir(), h.CanonicalName()+".myco"))

		err := writeTextToDisk(H, data, hop)
		if err != nil {
			hop.Abort()
			return err
		}

		hyphae.Insert(H)
		backlinks.UpdateBacklinksAfterEdit(H, "")
	case *hyphae.MediaHypha:
		oldText, err := FetchTextFile(h)
		if err != nil {
			hop.Abort()
			return err
		}

		// TODO: that []byte(...) part should be removed
		if bytes.Compare(data, []byte(oldText)) == 0 {
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
	case *hyphae.TextualHypha:
		oldText, err := FetchTextFile(h)
		if err != nil {
			hop.Abort()
			return err
		}

		// TODO: that []byte(...) part should be removed
		if bytes.Compare(data, []byte(oldText)) == 0 {
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

func historyMessageForMediaUpload(h hyphae.Hypha, mime string) string {
	return fmt.Sprintf("Upload media for ‘%s’ with type ‘%s’", h.CanonicalName(), mime)
}

// writeMediaToDisk saves the given data with the given mime type for the given hypha to the disk and returns the path to the saved file and an error, if any.
func writeMediaToDisk(h hyphae.Hypha, mime string, data []byte) (string, error) {
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
func UploadBinary(h hyphae.Hypha, mime string, file multipart.File, u *user.User) error {

	// Privilege check
	if !u.CanProceed("upload-binary") {
		rejectAttachLog(h, u, "no rights")
		return errors.New("ui.act_no_rights")
	}

	// Hypha name exploit check
	if !hyphae.IsValidName(h.CanonicalName()) {
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
	case *hyphae.EmptyHypha:
		H := hyphae.ExtendEmptyToMedia(h, uploadedFilePath)
		hyphae.Insert(H)
	case *hyphae.TextualHypha:
		hyphae.Insert(hyphae.ExtendTextualToMedia(h, uploadedFilePath))
	case *hyphae.MediaHypha: // If this is not the first media the hypha gets
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
