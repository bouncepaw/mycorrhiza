package backlinks

import (
	"github.com/bouncepaw/mycomarkup/v3"
	"github.com/bouncepaw/mycomarkup/v3/links"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"github.com/bouncepaw/mycomarkup/v3/tools"
	"github.com/bouncepaw/mycorrhiza/hyphae"
)

// UpdateBacklinksAfterEdit is a creation/editing hook for backlinks index
func UpdateBacklinksAfterEdit(h hyphae.Hypha, oldText string) {
	oldLinks := extractHyphaLinksFromContent(h.CanonicalName(), oldText)
	newLinks := extractHyphaLinks(h)
	backlinkConveyor <- backlinkIndexEdit{h.CanonicalName(), oldLinks, newLinks}
}

// UpdateBacklinksAfterDelete is a deletion hook for backlinks index
func UpdateBacklinksAfterDelete(h hyphae.Hypha, oldText string) {
	oldLinks := extractHyphaLinksFromContent(h.CanonicalName(), oldText)
	backlinkConveyor <- backlinkIndexDeletion{h.CanonicalName(), oldLinks}
}

// UpdateBacklinksAfterRename is a renaming hook for backlinks index
func UpdateBacklinksAfterRename(h hyphae.Hypha, oldName string) {
	actualLinks := extractHyphaLinks(h)
	backlinkConveyor <- backlinkIndexRenaming{oldName, h.CanonicalName(), actualLinks}
}

// extractHyphaLinks extracts hypha links from a desired hypha
func extractHyphaLinks(h hyphae.Hypha) []string {
	return extractHyphaLinksFromContent(h.CanonicalName(), fetchText(h))
}

// extractHyphaLinksFromContent extracts local hypha links from the provided text.
func extractHyphaLinksFromContent(hyphaName string, contents string) []string {
	ctx, _ := mycocontext.ContextFromStringInput(hyphaName, contents)
	linkVisitor, getLinks := tools.LinkVisitor(ctx)
	// Ignore the result of BlockTree because we call it for linkVisitor.
	_ = mycomarkup.BlockTree(ctx, linkVisitor)
	foundLinks := getLinks()
	var result []string
	for _, link := range foundLinks {
		if link.OfKind(links.LinkLocalHypha) {
			result = append(result, link.TargetHypha())
		}
	}
	return result
}
