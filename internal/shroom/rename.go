package shroom

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/categories"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"
	hyphae2 "github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Rename renames the old hypha to the new name and makes a history record about that. Call if and only if the user has the permission to rename.
func Rename(oldHypha hyphae2.ExistingHypha, newName string, recursive bool, leaveRedirections bool, u *user.User) error {
	// * bouncepaw hates this function and related renaming functions
	if newName == "" {
		rejectRenameLog(oldHypha, u, "no new name given")
		return errors.New("ui.rename_noname_tip")
	}

	if !hyphae2.IsValidName(newName) {
		rejectRenameLog(oldHypha, u, fmt.Sprintf("new name ‘%s’ invalid", newName))
		return errors.New("ui.rename_badname_tip") // FIXME: There is a bug related to this.
	}

	switch targetHypha := hyphae2.ByName(newName); targetHypha.(type) {
	case hyphae2.ExistingHypha:
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
		hyphae2.RenameHyphaTo(h, newName, replaceName)
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
		emptyHypha = hyphae2.ByName(oldName)
	)
	switch emptyHypha := emptyHypha.(type) {
	case *hyphae2.EmptyHypha:
		h := hyphae2.ExtendEmptyToTextual(emptyHypha, filepath.Join(files.HyphaeDir(), oldName+".myco"))
		hyphae2.Insert(h)
		categories.AddHyphaToCategory(oldName, cfg.RedirectionCategory)
		defer backlinks.UpdateBacklinksAfterEdit(h, "")
		return writeTextToDisk(h, []byte(text), hop)
	default:
		return errors.New("invalid state for hypha " + oldName + " renamed to " + newName)
	}
}

func findHyphaeToRename(superhypha hyphae2.ExistingHypha, recursive bool) []hyphae2.ExistingHypha {
	hyphaList := []hyphae2.ExistingHypha{superhypha}
	if recursive {
		hyphaList = append(hyphaList, hyphae2.Subhyphae(superhypha)...)
	}
	return hyphaList
}

func renamingPairs(hyphaeToRename []hyphae2.ExistingHypha, replaceName func(string) string) (map[string]string, error) {
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
		case *hyphae2.MediaHypha:
			renameMap[h.MediaFilePath()] = replaceName(h.MediaFilePath())
		}
		h.Unlock()
	}
	if firstFailure, ok := hyphae2.AreFreeNames(newNames...); !ok {
		return nil, errors.New("Hypha " + firstFailure + " already exists")
	}
	return renameMap, nil
}
