package shroom

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/backlinks"
	"github.com/bouncepaw/mycorrhiza/categories"
	"regexp"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Rename renames the old hypha to the new name and makes a history record about that. Call if and only if the user has the permission to rename.
func Rename(oldHypha hyphae.ExistingHypha, newName string, recursive bool, u *user.User) error {
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
			// Can we drop that util.CanonicalName?
			return re.ReplaceAllString(util.CanonicalName(str), newName)
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

	hop.WithFilesRenamed(renameMap).Apply()

	if len(hop.Errs) != 0 {
		return hop.Errs[0]
	}

	for _, h := range hyphaeToRename {
		var (
			oldName = h.CanonicalName()
			newName = replaceName(oldName)
		)
		hyphae.RenameHyphaTo(h, newName, replaceName)
		backlinks.UpdateBacklinksAfterRename(h, oldName)
		categories.RenameHyphaInAllCategories(oldName, newName)
	}

	return nil
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
