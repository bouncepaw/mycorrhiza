package shroom

import (
	"fmt"
	"log/slog"
	"path"
	"path/filepath"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/categories"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Rename renames the old hypha to the new name and makes a history record about that. Call if and only if the user has the permission to rename.
func Rename(oldHypha hyphae.ExistingHypha, newName string, recursive bool, leaveRedirections bool, u *usersvc.User) error {
	// * bouncepaw hates this function and related renaming functions
	if newName == "" {
		slog.Info("Reject rename",
			"hyphaName", oldHypha.CanonicalName(),
			"username", u.Name,
			"reason", "no new name given")
		return fmt.Errorf("ui.rename_noname_tip")
	}

	if !hyphae.IsValidName(newName) {
		slog.Info("Reject rename",
			"hyphaName", oldHypha.CanonicalName(),
			"username", u.Name,
			"reason", "new name is invalid",
			"newName", newName)
		return fmt.Errorf("ui.rename_badname_tip") // FIXME: There is a bug related to this.
	}

	switch targetHypha := hyphae.ByName(newName); targetHypha.(type) {
	case hyphae.ExistingHypha:
		if targetHypha.CanonicalName() == oldHypha.CanonicalName() {
			return nil
		}
		slog.Info("Reject rename",
			"hyphaName", oldHypha.CanonicalName(),
			"username", u.Name,
			"reason", "new name is taken already",
			"newName", newName)
		return fmt.Errorf("ui.rename_taken_tip") // FIXME: There is a bug related to this.
	}

	var (
		r              = renaming{hyphaeDir: files.HyphaeDir(), from: oldHypha.CanonicalName(), to: newName}
		hyphaeToRename = findHyphaeToRename(oldHypha, recursive)
		newNames       = r.newNames(hyphaeToRename)
	)

	if firstTaken, ok := hyphae.AreFreeNames(newNames...); !ok {
		return fmt.Errorf("Hypha %s already exists", firstTaken)
	}

	movedFiles, err := r.movedFiles(hyphaeToRename)
	if err != nil {
		return err
	}

	hop := history.Operation(history.TypeRenameHypha).WithUser(u)

	if len(hyphaeToRename) > 0 {
		hop.WithMsg(fmt.Sprintf(
			"Rename ‘%s’ to ‘%s’ recursively",
			oldHypha.CanonicalName(),
			newName))
	} else {
		hop.WithMsg(fmt.Sprintf(
			"Rename ‘%s’ to ‘%s’",
			oldHypha.CanonicalName(),
			newName))
	}

	if hop.WithFilesRenamed(movedFiles).HasErrors() {
		hop.Abort()
		return hop.Errs[0]
	}

	for i, h := range hyphaeToRename {
		var (
			hyphaOldName = h.CanonicalName()
			hyphaNewName = newNames[i]
		)
		hyphae.RenameHyphaTo(h, hyphaNewName, movedFiles)
		backlinks.UpdateBacklinksAfterRename(h, hyphaOldName)
		catsvc.RenameHyphaInAllCategories(hyphaOldName, hyphaNewName)
		if leaveRedirections {
			if err := leaveRedirection(hyphaOldName, hyphaNewName, hop); err != nil {
				hop.WithErrAbort(err)
				return err
			}
		}
	}

	hop.Apply()

	return nil
}

const redirectionTemplate = `=> %[1]s | 👁️➡️ %[2]s
<= %[1]s | full
`

func leaveRedirection(oldName, newName string, hop *history.Op) error {
	var (
		text       = fmt.Sprintf(redirectionTemplate, newName, util.BeautifulName(newName))
		emptyHypha = hyphae.ByName(oldName)
	)
	switch emptyHypha := emptyHypha.(type) {
	case *hyphae.EmptyHypha:
		h := hyphae.ExtendEmptyToTextual(emptyHypha, filepath.Join(files.HyphaeDir(), oldName+".myco"))
		hyphae.Insert(h)
		catsvc.AddHyphaToCategory(oldName, cfg.RedirectionCategory)
		defer backlinks.UpdateBacklinksAfterEdit(h, "")
		return writeTextToDisk(h, []byte(text), hop)
	default:
		return fmt.Errorf("invalid state for hypha %s renamed to %s", oldName, newName)
	}
}

func findHyphaeToRename(superhypha hyphae.ExistingHypha, recursive bool) []hyphae.ExistingHypha {
	hyphaList := []hyphae.ExistingHypha{superhypha}
	if recursive {
		hyphaList = append(hyphaList, hyphae.Subhyphae(superhypha)...)
	}
	return hyphaList
}

// renaming is the renaming of hypha from to to.
// It applies to the subhyphae
// of from as well: from/sub becomes to/sub. Names are substituted literally.
type renaming struct {
	hyphaeDir string
	from, to  string
}

// name is r.from or one of its subhyphae.
func (r renaming) newName(name string) string {
	return r.to + strings.TrimPrefix(name, r.from)
}

func (r renaming) newNames(hs []hyphae.ExistingHypha) []string {
	names := make([]string, len(hs))
	for i, h := range hs {
		names[i] = r.newName(h.CanonicalName())
	}
	return names
}

func (r renaming) newFilePath(filePath string) (string, error) {
	rel := util.CanonicalName(strings.TrimPrefix(filePath, r.hyphaeDir+"/"))
	if !strings.HasPrefix(rel, r.from+".") && !strings.HasPrefix(rel, r.from+"/") {
		return "", fmt.Errorf("file %s does not belong to hypha %s", filePath, r.from)
	}

	newRel := r.newName(rel)
	if r.to == "" || !files.IsInsideHyphaeDir(newRel) {
		return "", fmt.Errorf("refusing to move file %s to %s", filePath, newRel)
	}
	return path.Join(r.hyphaeDir, newRel), nil
}

func (r renaming) movedFiles(hs []hyphae.ExistingHypha) (map[string]string, error) {
	moved := make(map[string]string)
	for _, filePath := range hyphae.FilePaths(hs...) {
		newPath, err := r.newFilePath(filePath)
		if err != nil {
			return nil, err
		}
		moved[filePath] = newPath
	}
	return moved, nil
}
