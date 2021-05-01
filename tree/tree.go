package tree

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"
)

func findSiblingsAndDescendants(hyphaName string) ([]*sibling, map[string]bool) {
	var (
		siblings     = []*sibling{&sibling{hyphaName, 0, 0}}
		siblingCheck = func(h *hyphae.Hypha) hyphae.CheckResult {
			if path.Dir(hyphaName) == path.Dir(h.Name) && h.Name != hyphaName {
				siblings = append(siblings, &sibling{h.Name, 0, 0})
			}
			return hyphae.CheckContinue
		}

		descendantsPool = make(map[string]bool, 0)
		descendantCheck = func(h *hyphae.Hypha) hyphae.CheckResult {
			if strings.HasPrefix(h.Name, hyphaName+"/") {
				descendantsPool[h.Name] = true
			}
			return hyphae.CheckContinue
		}

		i7n = hyphae.NewIteration()
	)
	i7n.AddCheck(siblingCheck)
	i7n.AddCheck(descendantCheck)
	i7n.Ignite()
	sort.Slice(siblings, func(i, j int) bool {
		return siblings[i].name < siblings[j].name
	})
	return siblings, descendantsPool
}

func countSubhyphae(siblings []*sibling) {
	var (
		subhyphaCheck = func(h *hyphae.Hypha) hyphae.CheckResult {
			for _, s := range siblings {
				if path.Dir(h.Name) == s.name {
					s.directSubhyphaeCount++
					return hyphae.CheckContinue
				} else if strings.HasPrefix(h.Name, s.name+"/") {
					s.indirectSubhyphaeCount++
					return hyphae.CheckContinue
				}
			}
			return hyphae.CheckContinue
		}
		i7n = hyphae.NewIteration()
	)
	i7n.AddCheck(subhyphaCheck)
	i7n.Ignite()
}

// Tree generates a tree for `hyphaName` as html and returns next and previous hyphae if any.
func Tree(hyphaName string) (siblingsHTML, childrenHTML, prev, next string) {
	children := make([]child, 0)
	I := 0
	// The tree is generated in two iterations of hyphae storage:
	// 1. Find all siblings (sorted) and descendants' names
	// 2. Count how many subhyphae siblings have
	//
	// We also have to figure out what is going on with the descendants: who is a child of whom. We do that in parallel with (2) because we can.
	// One of the siblings is the hypha with name `hyphaName`
	siblings, descendantsPool := findSiblingsAndDescendants(hyphaName)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		countSubhyphae(siblings)
		wg.Done()
	}()
	go func() {
		children = figureOutChildren(hyphaName, descendantsPool).children
		wg.Done()
	}()
	wg.Wait()

	for i, s := range siblings {
		if s.name == hyphaName {
			I = i
			siblingsHTML += fmt.Sprintf(`<li class="relative-hyphae__entry relative-hyphae__entry_this"><span>%s</span></li>`, util.BeautifulName(path.Base(hyphaName)))
		} else {
			siblingsHTML += siblingHTML(s)
		}
	}
	if I != 0 && len(siblings) > 1 {
		prev = siblings[I-1].name
	}
	if I != len(siblings)-1 && len(siblings) > 1 {
		next = siblings[I+1].name
	}
	return fmt.Sprintf(`<ul class="relative-hyphae__list">%s</ul>`, siblingsHTML), subhyphaeMatrix(children), prev, next
}

type child struct {
	name     string
	children []child
}

func figureOutChildren(hyphaName string, subhyphaePool map[string]bool) child {
	var (
		nestLevel = strings.Count(hyphaName, "/")
		adopted   = make([]child, 0)
	)
	for subhyphaName, _ := range subhyphaePool {
		subnestLevel := strings.Count(subhyphaName, "/")
		if subnestLevel-1 == nestLevel && path.Dir(subhyphaName) == hyphaName {
			delete(subhyphaePool, subhyphaName)
			adopted = append(adopted, figureOutChildren(subhyphaName, subhyphaePool))
		}
	}
	return child{hyphaName, adopted}
}

type sibling struct {
	name                   string
	directSubhyphaeCount   int
	indirectSubhyphaeCount int
}

func subhyphaeMatrix(children []child) (html string) {
	sort.Slice(children, func(i, j int) bool {
		return children[i].name < children[j].name
	})
	for _, child := range children {
		html += childHTML(&child)
	}
	return html
}
