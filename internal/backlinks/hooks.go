package backlinks

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/tools"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/mycoopts"
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
	ctx, _ := mycocontext.ContextFromStringInput(contents, mycoopts.MarkupOptions(hyphaName))
	linkVisitor, getLinks := tools.LinkVisitor(ctx)
	// Ignore the result of BlockTree because we call it for linkVisitor.
	_ = mycomarkup.BlockTree(ctx, linkVisitor)
	foundLinks := getLinks()
	var result []string
	for _, link := range foundLinks {
		switch link := link.(type) {
		case *links.LocalLink:
			result = append(result, link.Target(ctx))
		}
	}
	return result
}
