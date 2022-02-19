package shroom

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"
	"regexp"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func canRenameThisToThat(oh hyphae.Hypha, nh hyphae.Hypha, u *user.User, lc *l18n.Localizer) (errtitle string, err error) {
	switch nh.(type) {
	case *hyphae.EmptyHypha:
	default:
		rejectRenameLog(oh, u, fmt.Sprintf("name ‘%s’ taken already", nh.CanonicalName()))
		return lc.Get("ui.rename_taken"), fmt.Errorf(lc.Get("ui.rename_taken_tip", &l18n.Replacements{"name": "<a href='/hypha/%[1]s'>%[1]s</a>"}), nh.CanonicalName())
	}

	if nh.CanonicalName() == "" {
		rejectRenameLog(oh, u, "no new name given")
		return lc.Get("ui.rename_noname"), errors.New(lc.Get("ui.rename_noname_tip"))
	}

	if !hyphae.IsValidName(nh.CanonicalName()) {
		rejectRenameLog(oh, u, fmt.Sprintf("new name ‘%s’ invalid", nh.CanonicalName()))
		return lc.Get("ui.rename_badname"), errors.New(lc.Get("ui.rename_badname_tip", &l18n.Replacements{"chars": "<code>^?!:#@&gt;&lt;*|\"\\'&amp;%</code>"}))
	}

	return "", nil
}

// RenameHypha renames hypha from old name `hyphaName` to `newName` and makes a history record about that. If `recursive` is `true`, its subhyphae will be renamed the same way.
func RenameHypha(h hyphae.Hypha, newHypha hyphae.Hypha, recursive bool, u *user.User, lc *l18n.Localizer) (hop *history.Op, errtitle string) {
	newHypha.Lock()
	defer newHypha.Unlock()
	hop = history.Operation(history.TypeRenameHypha)

	if errtitle, err := CanRename(u, h, lc); errtitle != "" {
		hop.WithErrAbort(err)
		return hop, errtitle
	}
	if errtitle, err := canRenameThisToThat(h, newHypha, u, lc); errtitle != "" {
		hop.WithErrAbort(err)
		return hop, errtitle
	}

	var (
		re          = regexp.MustCompile(`(?i)` + h.CanonicalName())
		replaceName = func(str string) string {
			return re.ReplaceAllString(util.CanonicalName(str), newHypha.CanonicalName())
		}
		hyphaeToRename = findHyphaeToRename(h.(hyphae.ExistingHypha), recursive)
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
		WithMsg(fmt.Sprintf(renameMsg, h.CanonicalName(), newHypha.CanonicalName())).
		WithUser(u).
		Apply()
	if len(hop.Errs) == 0 {
		for _, H := range hyphaeToRename {
			h := H.(hyphae.ExistingHypha) // ontology think
			oldName := h.CanonicalName()
			hyphae.RenameHyphaTo(h, replaceName(h.CanonicalName()), replaceName)
			backlinks.UpdateBacklinksAfterRename(h, oldName)
		}
	}
	return hop, ""
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
