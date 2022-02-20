package shroom

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
)

// Delete deletes the hypha and makes a history record about that.
func Delete(u *user.User, h hyphae.ExistingHypha) error {
	hop := history.
		Operation(history.TypeDeleteHypha).
		WithMsg(fmt.Sprintf("Delete ‘%s’", h.CanonicalName())).
		WithUser(u)

	originalText, _ := FetchTextFile(h)
	switch h := h.(type) {
	case *hyphae.MediaHypha:
		hop.WithFilesRemoved(h.MediaFilePath(), h.TextFilePath())
	case *hyphae.TextualHypha:
		hop.WithFilesRemoved(h.TextFilePath())
	}
	if hop.Apply().HasErrors() {
		return hop.Errs[0]
	}
	backlinks.UpdateBacklinksAfterDelete(h, originalText)
	hyphae.DeleteHypha(h)
	return nil
}
