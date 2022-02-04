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

func findSiblings(hyphaName string) []*sibling {
	parentHyphaName := ""
	if hyphaRawDir := path.Dir(hyphaName); hyphaRawDir != "." {
		parentHyphaName = hyphaRawDir
	}
	var (
		siblingsMap  = make(map[string]bool)
		siblingCheck = func(h hyphae.Hypher) hyphae.CheckResult {
			switch {
			case h.CanonicalName() == hyphaName, // MediaHypha is no sibling of itself
				h.CanonicalName() == parentHyphaName: // Parent hypha is no sibling of its child
				return hyphae.CheckContinue
			}
			if (parentHyphaName != "" && strings.HasPrefix(h.CanonicalName(), parentHyphaName+"/")) ||
				(parentHyphaName == "") {
				var (
					rawSubPath = strings.TrimPrefix(h.CanonicalName(), parentHyphaName)[1:]
					slashIdx   = strings.IndexRune(rawSubPath, '/')
				)
				if slashIdx > -1 {
					var sibPath = h.CanonicalName()[:slashIdx+len(parentHyphaName)+1]
					if _, exists := siblingsMap[sibPath]; !exists {
						siblingsMap[sibPath] = false
					}
				} else { // it is a straight sibling
					siblingsMap[h.CanonicalName()] = true
				}
			}
			return hyphae.CheckContinue
		}

		i7n = hyphae.NewIteration()
	)
	siblingsMap[hyphaName] = true

	i7n.AddCheck(siblingCheck)
	i7n.Ignite()

	siblings := make([]*sibling, len(siblingsMap))
	sibIdx := 0
	for sibName, exists := range siblingsMap {
		siblings[sibIdx] = &sibling{sibName, 0, 0, exists}
		sibIdx++
	}
	sort.Slice(siblings, func(i, j int) bool {
		return siblings[i].name < siblings[j].name
	})
	return siblings
}

func countSubhyphae(siblings []*sibling) {
	var (
		subhyphaCheck = func(h hyphae.Hypher) hyphae.CheckResult {
			for _, s := range siblings {
				if path.Dir(h.CanonicalName()) == s.name {
					s.directSubhyphaeCount++
					return hyphae.CheckContinue
				} else if strings.HasPrefix(h.CanonicalName(), s.name+"/") {
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
	// 1. Find all siblings (sorted)
	// 2. Count how many subhyphae siblings have
	//
	// We also have to figure out what is going on with the descendants: who is a child of whom. We do that in parallel with (2) because we can.
	// One of the siblings is the hypha with name `hyphaName`
	var siblings []*sibling

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		siblings = findSiblings(hyphaName)
		countSubhyphae(siblings)
		wg.Done()
	}()
	go func() {
		children = figureOutChildren(hyphaName).children
		wg.Done()
	}()
	wg.Wait()

	for i, s := range siblings {
		if s.name == hyphaName {
			I = i
			siblingsHTML += fmt.Sprintf(`<li class="sibling-hyphae__entry sibling-hyphae__entry_this"><span>%s</span></li>`, util.BeautifulName(path.Base(hyphaName)))
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
	return fmt.Sprintf(`<ul class="sibling-hyphae__list">%s</ul>`, siblingsHTML), subhyphaeMatrix(children), prev, next
}

type child struct {
	name     string
	exists   bool
	children []child
}

func figureOutChildren(hyphaName string) child {
	var (
		descPrefix = hyphaName + "/"
		child      = child{hyphaName, true, make([]child, 0)}
	)

	for desc := range hyphae.YieldExistingHyphae() {
		var descName = desc.CanonicalName()
		if strings.HasPrefix(descName, descPrefix) {
			var subPath = strings.TrimPrefix(descName, descPrefix)
			addHyphaToChild(descName, subPath, &child)
		}
	}

	return child
}

func addHyphaToChild(hyphaName, subPath string, child *child) {
	// when hyphaName = "root/a/b", subPath = "a/b", and child.name = "root"
	// addHyphaToChild("root/a/b", "b", child{"root/a"})
	// when hyphaName = "root/a/b", subPath = "b", and child.name = "root/a"
	// set .exists=true for "root/a/b", and create it if it isn't there already
	var exists = !strings.Contains(subPath, "/")
	if exists {
		var subchild = findOrCreateSubchild(subPath, child)
		subchild.exists = true
	} else {
		var (
			firstSlash = strings.IndexRune(subPath, '/')
			firstDir   = subPath[:firstSlash]
			restOfPath = subPath[firstSlash+1:]
			subchild   = findOrCreateSubchild(firstDir, child)
		)
		addHyphaToChild(hyphaName, restOfPath, subchild)
	}
}

func findOrCreateSubchild(name string, baseChild *child) *child {
	// when name = "a", and baseChild.name = "root"
	// if baseChild.children contains "root/a", return it
	// else create it and return that
	var fullName = baseChild.name + "/" + name
	for i := range baseChild.children {
		if baseChild.children[i].name == fullName {
			return &baseChild.children[i]
		}
	}
	baseChild.children = append(baseChild.children, child{fullName, false, make([]child, 0)})
	return &baseChild.children[len(baseChild.children)-1]
}

type sibling struct {
	name                   string
	directSubhyphaeCount   int
	indirectSubhyphaeCount int
	exists                 bool
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
