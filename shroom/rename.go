package shroom

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func canRenameThisToThat(oh *hyphae.Hypha, nh *hyphae.Hypha, u *user.User) (err error, errtitle string) {
	if nh.Exists {
		rejectRenameLog(oh, u, fmt.Sprintf("name ‘%s’ taken already", nh.Name))
		return errors.New(fmt.Sprintf("Hypha named <a href='/hypha/%[1]s'>%[1]s</a> already exists, cannot rename", nh.Name)), "Name taken"
	}

	if nh.Name == "" {
		rejectRenameLog(oh, u, "no new name given")
		return errors.New("No new name is given"), "No name given"
	}

	if !hyphae.HyphaPattern.MatchString(nh.Name) {
		rejectRenameLog(oh, u, fmt.Sprintf("new name ‘%s’ invalid", nh.Name))
		return errors.New("Invalid new name. Names cannot contain characters <code>^?!:#@&gt;&lt;*|\"\\'&amp;%</code>"), "Invalid name"
	}

	return nil, ""
}

// RenameHypha renames hypha from old name `hyphaName` to `newName` and makes a history record about that. If `recursive` is `true`, its subhyphae will be renamed the same way.
func RenameHypha(h *hyphae.Hypha, newHypha *hyphae.Hypha, recursive bool, u *user.User) (hop *history.HistoryOp, errtitle string) {
	newHypha.Lock()
	defer newHypha.Unlock()
	hop = history.Operation(history.TypeRenameHypha)

	if err, errtitle := CanRename(u, h); errtitle != "" {
		hop.WithError(err)
		return hop, errtitle
	}
	if err, errtitle := canRenameThisToThat(h, newHypha, u); errtitle != "" {
		hop.WithError(err)
		return hop, errtitle
	}

	var (
		re          = regexp.MustCompile(`(?i)` + h.Name)
		replaceName = func(str string) string {
			return re.ReplaceAllString(util.CanonicalName(str), newHypha.Name)
		}
		hyphaeToRename = findHyphaeToRename(h, recursive)
		renameMap, err = renamingPairs(hyphaeToRename, replaceName)
		renameMsg      = "Rename ‘%s’ to ‘%s’"
	)
	if err != nil {
		hop.Errs = append(hop.Errs, err)
		return hop, hop.FirstErrorText()
	}
	if recursive && len(hyphaeToRename) > 0 {
		renameMsg += " recursively"
	}
	hop.WithFilesRenamed(renameMap).
		WithMsg(fmt.Sprintf(renameMsg, h.Name, newHypha.Name)).
		WithUser(u).
		Apply()
	if len(hop.Errs) == 0 {
		for _, h := range hyphaeToRename {
			h.RenameTo(replaceName(h.Name))
			h.Lock()
			h.TextPath = replaceName(h.TextPath)
			h.BinaryPath = replaceName(h.BinaryPath)
			h.Unlock()
		}
	}
	return hop, ""
}

func findHyphaeToRename(superhypha *hyphae.Hypha, recursive bool) []*hyphae.Hypha {
	hyphae := []*hyphae.Hypha{superhypha}
	if recursive {
		hyphae = append(hyphae, superhypha.Subhyphae()...)
	}
	return hyphae
}

func renamingPairs(hyphaeToRename []*hyphae.Hypha, replaceName func(string) string) (map[string]string, error) {
	renameMap := make(map[string]string)
	newNames := make([]string, len(hyphaeToRename))
	for _, h := range hyphaeToRename {
		h.RLock()
		newNames = append(newNames, replaceName(h.Name))
		if h.TextPath != "" {
			renameMap[h.TextPath] = replaceName(h.TextPath)
		}
		if h.BinaryPath != "" {
			renameMap[h.BinaryPath] = replaceName(h.BinaryPath)
		}
		h.RUnlock()
	}
	if firstFailure, ok := hyphae.AreFreeNames(newNames...); !ok {
		return nil, errors.New("Hypha " + firstFailure + " already exists")
	}
	return renameMap, nil
}
