package shroom

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"
	"github.com/bouncepaw/mycorrhiza/mimetype"
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
	"github.com/bouncepaw/mycorrhiza/user"
)

func historyMessageForTextUpload(h hyphae.Hypher, userMessage string) string {
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

func writeTextToDiskForEmptyHypha(eh *hyphae.EmptyHypha, data []byte) error {
	h := hyphae.FillEmptyHyphaUpToTextualHypha(eh, filepath.Join(files.HyphaeDir(), eh.CanonicalName()+".myco"))

	return writeTextToDiskForNonEmptyHypha(h, data)
}

func writeTextToDiskForNonEmptyHypha(h *hyphae.MediaHypha, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(h.TextPartPath()), 0777); err != nil {
		return err
	}

	if err := os.WriteFile(h.TextPartPath(), data, 0666); err != nil {
		return err
	}
	return nil
}

// UploadText edits the hypha's text part and makes a history record about that.
func UploadText(h hyphae.Hypher, data []byte, userMessage string, u *user.User, lc *l18n.Localizer) (hop *history.Op, errtitle string) {
	hop = history.
		Operation(history.TypeEditText).
		WithMsg(historyMessageForTextUpload(h, userMessage))

	// Privilege check
	if !u.CanProceed("upload-text") {
		rejectEditLog(h, u, "no rights")
		return hop.WithErrAbort(errors.New(lc.Get("ui.act_norights_edit"))), lc.Get("ui.act_no_rights")
	}

	// Hypha name exploit check
	if !hyphae.IsValidName(h.CanonicalName()) {
		// We check for the name only. I suppose the filepath would be valid as well.
		err := errors.New("invalid hypha name")
		return hop.WithErrAbort(err), err.Error()
	}

	// Empty data check
	if len(bytes.TrimSpace(data)) == 0 { // if nothing but whitespace
		switch h := h.(type) {
		case *hyphae.EmptyHypha:
			// It's ok, just like cancel button.
			return hop.Abort(), ""
		case *hyphae.MediaHypha:
			switch h.Kind() {
			case hyphae.HyphaMedia:
				// Writing no description, it's ok, just like cancel button.
				return hop.Abort(), ""
			case hyphae.HyphaText:
				// What do you want passing nothing for a textual hypha?
				return hop.WithErrAbort(errors.New("No data passed")), "Empty"
			}
		}
	}

	// At this point, we have a savable user-generated Mycomarkup document. Gotta save it.

	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		err := writeTextToDiskForEmptyHypha(h, data)
		if err != nil {
			return hop.WithErrAbort(err), err.Error()
		}

		hyphae.InsertIfNew(h)
	case *hyphae.MediaHypha:
		oldText, err := FetchTextPart(h)
		if err != nil {
			return hop.WithErrAbort(err), err.Error()
		}

		// TODO: that []byte(...) part should be removed
		if bytes.Compare(data, []byte(oldText)) == 0 {
			// No changes! Just like cancel button
			return hop.Abort(), ""
		}

		err = writeTextToDiskForNonEmptyHypha(h, data)
		if err != nil {
			return hop.WithErrAbort(err), err.Error()
		}

		backlinks.UpdateBacklinksAfterEdit(h, oldText)
	}

	return hop.
		WithFiles(h.TextPartPath()).
		WithUser(u).
		Apply(), ""
}

// UploadBinary edits the hypha's media part and makes a history record about that.
func UploadBinary(h hyphae.Hypher, mime string, file multipart.File, u *user.User, lc *l18n.Localizer) (*history.Op, string) {
	var (
		hop       = history.Operation(history.TypeEditBinary).WithMsg(fmt.Sprintf("Upload attachment for ‘%s’ with type ‘%s’", h.CanonicalName(), mime))
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

	ext := mimetype.ToExtension(mime)

	var (
		fullPath       = filepath.Join(files.HyphaeDir(), h.CanonicalName()+ext)
		sourceFullPath = h.TextPartPath()
	)
	if !isValidPath(fullPath) || !hyphae.IsValidName(h.CanonicalName()) {
		err := errors.New("bad path")
		return hop.WithErrAbort(err), err.Error()
	}
	if h := h.(*hyphae.MediaHypha); hop.Type == history.TypeEditBinary {
		sourceFullPath = h.BinaryPath()
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		return hop.WithErrAbort(err), err.Error()
	}

	if err := os.WriteFile(fullPath, data, 0666); err != nil {
		return hop.WithErrAbort(err), err.Error()
	}

	switch h.(type) {
	case *hyphae.EmptyHypha:
	default:
		if sourceFullPath != fullPath && sourceFullPath != "" {
			if err := history.Rename(sourceFullPath, fullPath); err != nil {
				return hop.WithErrAbort(err), err.Error()
			}
			log.Println("Move", sourceFullPath, "to", fullPath)
		}
	}

	hyphae.InsertIfNew(h)

	switch h.(type) {
	case *hyphae.EmptyHypha:
	default:
		if h.HasTextPart() && hop.Type == history.TypeEditText && !history.FileChanged(fullPath) {
			return hop.Abort(), "No changes"
		}
	}

	// sic!
	h.(*hyphae.MediaHypha).SetBinaryPath(fullPath)
	return hop.WithFiles(fullPath).WithUser(u).Apply(), ""
}

func isValidPath(pathname string) bool {
	return strings.HasPrefix(pathname, files.HyphaeDir())
}
