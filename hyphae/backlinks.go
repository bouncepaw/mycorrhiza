package hyphae

import (
	"os"
	"sync"

	"github.com/bouncepaw/mycorrhiza/util"

	"github.com/bouncepaw/mycomarkup"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/links"
	"github.com/bouncepaw/mycomarkup/mycocontext"
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

var backlinkIndex = make(map[string]linkSet)
var backlinkIndexMutex = sync.Mutex{}

// IndexBacklinks traverses all text hyphae, extracts links from them and forms an initial index
func IndexBacklinks() {
	// It is safe to ignore the mutex, because there is only one worker.
	src := FilterTextHyphae(YieldExistingHyphae())
	for h := range src {
		links := ExtractHyphaLinksFromContent(h.Name, fetchText(h))
		for _, link := range links {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][h.Name] = struct{}{}
		}
	}
}

func BacklinksCount(h *Hypha) int {
	if _, exists := backlinkIndex[h.Name]; exists {
		return len(backlinkIndex[h.Name])
	}
	return 0
}

func BacklinksOnEdit(h *Hypha, oldText string) {
	backlinkIndexMutex.Lock()
	newLinks := toLinkSet(ExtractHyphaLinks(h))
	oldLinks := toLinkSet(ExtractHyphaLinksFromContent(h.Name, oldText))
	for link := range oldLinks {
		if _, exists := newLinks[link]; !exists {
			delete(backlinkIndex[link], h.Name)
		}
	}
	for link := range newLinks {
		if _, exists := oldLinks[link]; !exists {
			if _, exists := backlinkIndex[link]; !exists {
				backlinkIndex[link] = make(linkSet)
			}
			backlinkIndex[link][h.Name] = struct{}{}
		}
	}
	backlinkIndexMutex.Unlock()
}

func BacklinksOnDelete(h *Hypha, oldText string) {
	backlinkIndexMutex.Lock()
	oldLinks := ExtractHyphaLinksFromContent(h.Name, oldText)
	for _, link := range oldLinks {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, h.Name)
		}
	}
	backlinkIndexMutex.Unlock()
}

func BacklinksOnRename(h *Hypha, oldName string) {
	backlinkIndexMutex.Lock()
	actualLinks := ExtractHyphaLinks(h)
	for _, link := range actualLinks {
		if lSet, exists := backlinkIndex[link]; exists {
			delete(lSet, oldName)
			backlinkIndex[link][h.Name] = struct{}{}
		}
	}
	backlinkIndexMutex.Unlock()
}

// YieldHyphaBacklinks gets backlinks for a desired hypha, sorts and iterates over them
func YieldHyphaBacklinks(query string) <-chan string {
	hyphaName := util.CanonicalName(query)
	out := make(chan string)
	sorted := PathographicSort(out)
	go func() {
		links, exists := backlinkIndex[hyphaName]
		if exists {
			for link := range links {
				out <- link
			}
		}
		close(out)
	}()
	return sorted
}

// YieldHyphaLinks extracts hypha links from a desired hypha, sorts and iterates over them
func YieldHyphaLinks(query string) <-chan string {
	// That is merely a debug function, but it could be useful.
	// Should we extract them into link-specific subfile? -- chekoopa
	hyphaName := util.CanonicalName(query)
	out := make(chan string)
	go func() {
		var h = ByName(hyphaName)
		links := ExtractHyphaLinks(h)
		for _, link := range links {
			out <- link
		}
		close(out)
	}()
	return out
}

// ExtractHyphaLinks extracts hypha links from a desired hypha
func ExtractHyphaLinks(h *Hypha) []string {
	return ExtractHyphaLinksFromContent(h.Name, fetchText(h))
}

// ExtractHyphaLinksFromContent extracts hypha links from a provided text
func ExtractHyphaLinksFromContent(hyphaName string, contents string) []string {
	ctx, _ := mycocontext.ContextFromStringInput(hyphaName, contents)
	linkVisitor, getLinks := LinkVisitor(ctx)
	mycomarkup.BlockTree(ctx, linkVisitor)
	foundLinks := getLinks()
	var result []string
	for _, link := range foundLinks {
		if link.OfKind(links.LinkLocalHypha) {
			result = append(result, link.TargetHypha())
		}
	}
	return result
}

// LinkVisitor creates a visitor which extracts all the links
func LinkVisitor(ctx mycocontext.Context) (
	visitor func(block blocks.Block),
	result func() []links.Link,
) {
	var (
		collected []links.Link
	)
	var extractBlock func(block blocks.Block)
	extractBlock = func(block blocks.Block) {
		// fmt.Println(reflect.TypeOf(block))
		switch b := block.(type) {
		case blocks.Paragraph:
			extractBlock(b.Formatted)
		case blocks.Heading:
			extractBlock(b.GetContents())
		case blocks.List:
			for _, item := range b.Items {
				for _, sub := range item.Contents {
					extractBlock(sub)
				}
			}
		case blocks.Img:
			for _, entry := range b.Entries {
				extractBlock(entry)
			}
		case blocks.ImgEntry:
			collected = append(collected, *b.Srclink)
		case blocks.Transclusion:
			link := *links.From(b.Target, "", ctx.HyphaName())
			collected = append(collected, link)
		case blocks.LaunchPad:
			for _, rocket := range b.Rockets {
				extractBlock(rocket)
			}
		case blocks.Formatted:
			for _, line := range b.Lines {
				for _, span := range line {
					switch s := span.(type) {
					case blocks.InlineLink:
						collected = append(collected, *s.Link)
					}
				}
			}
		case blocks.RocketLink:
			if !b.IsEmpty {
				collected = append(collected, b.Link)
			}
		}
	}
	visitor = func(block blocks.Block) {
		extractBlock(block)
	}
	result = func() []links.Link {
		return collected
	}
	return
}
