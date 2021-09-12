package hyphae

import (
	"github.com/bouncepaw/mycomarkup/v2/tools"
	"os"

	"github.com/bouncepaw/mycorrhiza/util"

	"github.com/bouncepaw/mycomarkup/v2"
	"github.com/bouncepaw/mycomarkup/v2/links"
	"github.com/bouncepaw/mycomarkup/v2/mycocontext"
)

// Using set here seems like the most appropriate solution
type linkSet map[string]struct{}

func toLinkSet(xs []string) linkSet {
	result := make(linkSet)
	for _, x := range xs {
		result[x] = struct{}{}
	}
	return result
}

func fetchText(h *Hypha) string {
	if h.TextPath == "" {
		return ""
	}
	text, err := os.ReadFile(h.TextPath)
	if err == nil {
		return string(text)
	}
	return ""
}

// BacklinkIndexOperation is an operation for the backlink index. This operation is executed async-safe.
type BacklinkIndexOperation interface {
	Apply()
}

type BacklinkIndexEdit struct {
	Name     string
	OldLinks []string
	NewLinks []string
}

func (op BacklinkIndexEdit) Apply() {
	oldLinks := toLinkSet(op.OldLinks)
	newLinks := toLinkSet(op.NewLinks)
	for link := range oldLinks {
		if _, exists := newLinks[link]; !exists {
			delete(backlinkIndex[link], op.Name)
		}
	}
	for link := range newLinks {
		if _, exists := oldLinks[link]; !exists {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][op.Name] = struct{}{}
		}
	}
}

type BacklinkIndexDeletion struct {
	Name  string
	Links []string
}

func (op BacklinkIndexDeletion) Apply() {
	for _, link := range op.Links {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, op.Name)
		}
	}
}

type BacklinkIndexRenaming struct {
	OldName string
	NewName string
	Links   []string
}

func (op BacklinkIndexRenaming) Apply() {
	for _, link := range op.Links {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, op.OldName)
			backlinkIndex[link][op.NewName] = struct{}{}
		}
	}
}

var backlinkIndex = make(map[string]linkSet)
var backlinkConveyor = make(chan BacklinkIndexOperation, 64)

// I hope, the buffer size is enough -- chekoopa
//   Do we really need the buffer though? Dunno -- bouncepaw

// IndexBacklinks traverses all text hyphae, extracts links from them and forms an initial index
func IndexBacklinks() {
	// It is safe to ignore the mutex, because there is only one worker.
	src := FilterTextHyphae(YieldExistingHyphae())
	for h := range src {
		foundLinks := extractHyphaLinksFromContent(h.Name, fetchText(h))
		for _, link := range foundLinks {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][h.Name] = struct{}{}
		}
	}
}

// RunBacklinksConveyor runs an index operation processing loop
func RunBacklinksConveyor() {
	// It is supposed to run as a goroutine for all the time. So, don't blame the infinite loop.
	defer close(backlinkConveyor)
	for {
		(<-backlinkConveyor).Apply()
	}
}

// BacklinksCount returns the amount of backlinks to the hypha.
func BacklinksCount(h *Hypha) int {
	if _, exists := backlinkIndex[h.Name]; exists {
		return len(backlinkIndex[h.Name])
	}
	return 0
}

// BacklinksOnEdit is a creation/editing hook for backlinks index
func BacklinksOnEdit(h *Hypha, oldText string) {
	oldLinks := extractHyphaLinksFromContent(h.Name, oldText)
	newLinks := extractHyphaLinks(h)
	backlinkConveyor <- BacklinkIndexEdit{h.Name, oldLinks, newLinks}
}

// BacklinksOnDelete is a deletion hook for backlinks index
func BacklinksOnDelete(h *Hypha, oldText string) {
	oldLinks := extractHyphaLinksFromContent(h.Name, oldText)
	backlinkConveyor <- BacklinkIndexDeletion{h.Name, oldLinks}
}

// BacklinksOnRename is a renaming hook for backlinks index
func BacklinksOnRename(h *Hypha, oldName string) {
	actualLinks := extractHyphaLinks(h)
	backlinkConveyor <- BacklinkIndexRenaming{oldName, h.Name, actualLinks}
}

// YieldHyphaBacklinks gets backlinks for a desired hypha, sorts and iterates over them
func YieldHyphaBacklinks(query string) <-chan string {
	hyphaName := util.CanonicalName(query)
	out := make(chan string)
	sorted := PathographicSort(out)
	go func() {
		backlinks, exists := backlinkIndex[hyphaName]
		if exists {
			for link := range backlinks {
				out <- link
			}
		}
		close(out)
	}()
	return sorted
}

// extractHyphaLinks extracts hypha links from a desired hypha
func extractHyphaLinks(h *Hypha) []string {
	return extractHyphaLinksFromContent(h.Name, fetchText(h))
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
