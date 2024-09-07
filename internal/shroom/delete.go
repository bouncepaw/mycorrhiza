package shroom

import (
	"fmt"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/categories"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
)

// Delete deletes the hypha and makes a history record about that.
func Delete(u *user.User, h hyphae.ExistingHypha) error {
	hop := history.
		Operation(history.TypeDeleteHypha).
		WithMsg(fmt.Sprintf("Delete ‘%s’", h.CanonicalName())).
		WithUser(u)

	originalText, _ := hyphae.FetchMycomarkupFile(h)
	switch h := h.(type) {
	case *hyphae.MediaHypha:
		if h.HasTextFile() {
			hop.WithFilesRemoved(h.MediaFilePath(), h.TextFilePath())
		} else {
			hop.WithFilesRemoved(h.MediaFilePath())
		}
	case *hyphae.TextualHypha:
		hop.WithFilesRemoved(h.TextFilePath())
	}
	if hop.Apply().HasErrors() {
		return hop.Errs[0]
	}
	backlinks.UpdateBacklinksAfterDelete(h, originalText)
	categories.RemoveHyphaFromAllCategories(h.CanonicalName())
	hyphae.DeleteHypha(h)
	return nil
}
