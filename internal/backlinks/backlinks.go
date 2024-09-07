// Package backlinks maintains the index of backlinks and lets you update it and query it.
package backlinks

import (
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"log/slog"
	"os"
	"sort"

	"github.com/bouncepaw/mycorrhiza/util"
)

// yieldHyphaBacklinks gets backlinks for the desired hypha, sorts and yields them one by one.
func yieldHyphaBacklinks(hyphaName string) <-chan string {
	hyphaName = util.CanonicalName(hyphaName)
	out := make(chan string)
	sorted := hyphae.PathographicSort(out)
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

var backlinkConveyor = make(chan backlinkIndexOperation) // No need to buffer because these operations are rare.

// RunBacklinksConveyor runs an index operation processing loop. Call it somewhere in main.
func RunBacklinksConveyor() {
	// It is supposed to run as a goroutine for all the time. So, don't blame the infinite loop.
	defer close(backlinkConveyor)
	for {
		(<-backlinkConveyor).apply()
	}
}

var backlinkIndex = make(map[string]linkSet)

// IndexBacklinks traverses all text hyphae, extracts links from them and forms an initial index. Call it when indexing and reindexing hyphae.
func IndexBacklinks() {
	// It is safe to ignore the mutex, because there is only one worker.
	for h := range hyphae.FilterHyphaeWithText(hyphae.YieldExistingHyphae()) {
		foundLinks := extractHyphaLinksFromContent(h.CanonicalName(), fetchText(h))
		for _, link := range foundLinks {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][h.CanonicalName()] = struct{}{}
		}
	}
}

// BacklinksCount returns the amount of backlinks to the hypha. Pass canonical names.
func BacklinksCount(hyphaName string) int {
	if links, exists := backlinkIndex[hyphaName]; exists {
		return len(links)
	}
	return 0
}

func BacklinksFor(hyphaName string) []string {
	var backlinks []string
	for b := range yieldHyphaBacklinks(hyphaName) {
		backlinks = append(backlinks, b)
	}
	return backlinks
}

func Orphans() []string {
	var orphans []string
	for h := range hyphae.YieldExistingHyphae() {
		if BacklinksCount(h.CanonicalName()) == 0 {
			orphans = append(orphans, h.CanonicalName())
		}
	}
	sort.Strings(orphans)
	return orphans
}

// Using set here seems like the most appropriate solution
type linkSet map[string]struct{}

func toLinkSet(xs []string) linkSet {
	result := make(linkSet)
	for _, x := range xs {
		result[x] = struct{}{}
	}
	return result
}

func fetchText(h hyphae.Hypha) string {
	var path string
	switch h := h.(type) {
	case *hyphae.EmptyHypha:
		return ""
	case *hyphae.TextualHypha:
		path = h.TextFilePath()
	case *hyphae.MediaHypha:
		if !h.HasTextFile() {
			return ""
		}
		path = h.TextFilePath()
	}

	text, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read file", "path", path, "err", err, "hyphaName", h.CanonicalName())
		return ""
	}
	return string(text)
}

// backlinkIndexOperation is an operation for the backlink index. This operation is executed async-safe.
type backlinkIndexOperation interface {
	apply()
}

// backlinkIndexEdit contains data for backlink index update after a hypha edit
type backlinkIndexEdit struct {
	name     string
	oldLinks []string
	newLinks []string
}

// apply changes backlink index respective to the operation data
func (op backlinkIndexEdit) apply() {
	oldLinks := toLinkSet(op.oldLinks)
	newLinks := toLinkSet(op.newLinks)
	for link := range oldLinks {
		if _, exists := newLinks[link]; !exists {
			delete(backlinkIndex[link], op.name)
		}
	}
	for link := range newLinks {
		if _, exists := oldLinks[link]; !exists {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][op.name] = struct{}{}
		}
	}
}

// backlinkIndexDeletion contains data for backlink index update after a hypha deletion
type backlinkIndexDeletion struct {
	name  string
	links []string
}

// apply changes backlink index respective to the operation data
func (op backlinkIndexDeletion) apply() {
	for _, link := range op.links {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, op.name)
		}
	}
}

// backlinkIndexRenaming contains data for backlink index update after a hypha renaming
type backlinkIndexRenaming struct {
	oldName string
	newName string
	links   []string
}

// apply changes backlink index respective to the operation data
func (op backlinkIndexRenaming) apply() {
	for _, link := range op.links {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, op.oldName)
			backlinkIndex[link][op.newName] = struct{}{}
		}
	}
}
