package shroom

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/backlinks"
	"github.com/bouncepaw/mycorrhiza/categories"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Rename renames the old hypha to the new name and makes a history record about that. Call if and only if the user has the permission to rename.
func Rename(oldHypha hyphae.ExistingHypha, newName string, recursive bool, leaveRedirections bool, u *user.User) error {
	// * bouncepaw hates this function and related renaming functions
	if newName == "" {
		rejectRenameLog(oldHypha, u, "no new name given")
		return errors.New("ui.rename_noname_tip")
	}

	if !hyphae.IsValidName(newName) {
		rejectRenameLog(oldHypha, u, fmt.Sprintf("new name ‘%s’ invalid", newName))
		return errors.New("ui.rename_badname_tip") // FIXME: There is a bug related to this.
	}

	switch targetHypha := hyphae.ByName(newName); targetHypha.(type) {
	case hyphae.ExistingHypha:
		if targetHypha.CanonicalName() == oldHypha.CanonicalName() {
			return nil
		}
		rejectRenameLog(oldHypha, u, fmt.Sprintf("name ‘%s’ taken already", newName))
		return errors.New("ui.rename_taken_tip") // FIXME: There is a bug related to this.
	}

	var (
		re          = regexp.MustCompile(`(?i)` + oldHypha.CanonicalName())
		replaceName = func(str string) string {
			namepart := strings.TrimPrefix(str, files.HyphaeDir())
			// Can we drop that util.CanonicalName?:
			replaced := re.ReplaceAllString(util.CanonicalName(namepart), newName)
			return path.Join(files.HyphaeDir(), replaced)
		}
		hyphaeToRename = findHyphaeToRename(oldHypha, recursive)
		renameMap, err = renamingPairs(hyphaeToRename, replaceName)
	)

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

	hop.WithFilesRenamed(renameMap)

	if len(hop.Errs) != 0 {
		return hop.Errs[0]
	}

	for _, h := range hyphaeToRename {
		var (
			oldName = h.CanonicalName()
			newName = re.ReplaceAllString(oldName, newName)
		)
		hyphae.RenameHyphaTo(h, newName, replaceName)
		backlinks.UpdateBacklinksAfterRename(h, oldName)
		categories.RenameHyphaInAllCategories(oldName, newName)
		if leaveRedirections {
			if err := leaveRedirection(oldName, newName, hop); err != nil {
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
		categories.AddHyphaToCategory(oldName, cfg.RedirectionCategory)
		defer backlinks.UpdateBacklinksAfterEdit(h, "")
		return writeTextToDisk(h, []byte(text), hop)
	default:
		return errors.New("invalid state for hypha " + oldName + " renamed to " + newName)
	}
}

func findHyphaeToRename(superhypha hyphae.ExistingHypha, recursive bool) []hyphae.ExistingHypha {
	hyphaList := []hyphae.ExistingHypha{superhypha}
	if recursive {
		hyphaList = append(hyphaList, hyphae.Subhyphae(superhypha)...)
	}
	return hyphaList
}

func renamingPairs(hyphaeToRename []hyphae.ExistingHypha, replaceName func(string) string) (map[string]string, error) {
	var (
		renameMap = make(map[string]string)
		newNames  = make([]string, len(hyphaeToRename))
	)
	for _, h := range hyphaeToRename {
		h.Lock()
		newNames = append(newNames, replaceName(h.CanonicalName()))
		if h.HasTextFile() {
			renameMap[h.TextFilePath()] = replaceName(h.TextFilePath())
		}
		switch h := h.(type) {
		case *hyphae.MediaHypha:
			renameMap[h.MediaFilePath()] = replaceName(h.MediaFilePath())
		}
		h.Unlock()
	}
	if firstFailure, ok := hyphae.AreFreeNames(newNames...); !ok {
		return nil, errors.New("Hypha " + firstFailure + " already exists")
	}
	return renameMap, nil
}
