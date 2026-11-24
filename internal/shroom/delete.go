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
func Delete(u *usersvc.User, h hyphae.ExistingHypha) error {
	hop := history.
		Operation(history.TypeDeleteHypha).
		WithMsg(fmt.Sprintf("Delete ‘%s’", h.CanonicalName())).
		WithUser(u)

	originalText, _ := hyphae.FetchMycomarkupFile(h)
	hop.WithFilesRemoved(hyphae.FilePaths(h)...)
	if hop.Apply().HasErrors() {
		return hop.Errs[0]
	}
	backlinks.UpdateBacklinksAfterDelete(h, originalText)
	catsvc.RemoveHyphaFromAllCategories(h.CanonicalName())
	hyphae.DeleteHypha(h)
	return nil
}
