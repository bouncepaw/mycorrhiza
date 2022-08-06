package tree

import (
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"path"
	"sort"
	"strings"
)

// Tree returns the subhypha matrix as HTML and names of the next and previous hyphae (or empty strings).
func Tree(hyphaName string) (childrenHTML, prev, next string) {
	var (
		root             = child{hyphaName, true, make([]child, 0)}
		descendantPrefix = hyphaName + "/"
		parent           = path.Dir(hyphaName) // Beware, it might be . and whatnot.
		slashCount       = strings.Count(hyphaName, "/")
	)
	for h := range hyphae.YieldExistingHyphae() {
		name := h.CanonicalName()
		if strings.HasPrefix(name, descendantPrefix) {
			var subPath = strings.TrimPrefix(name, descendantPrefix)
			addHyphaToChild(name, subPath, &root)
			// A child is not a sibling, so we skip the rest.
			continue
		}

		// Skipping non-siblings.
		if !(path.Dir(name) == parent && slashCount == strings.Count(name, "/")) {
			continue
		}

		if (name < hyphaName) && (name > prev) {
			prev = name
		} else if (name > hyphaName) && (name < next || next == "") {
			next = name
		}
	}
	return subhyphaeMatrix(root.children), prev, next
}

type child struct {
	name     string
	exists   bool
	children []child
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

func subhyphaeMatrix(children []child) (html string) {
	sort.Slice(children, func(i, j int) bool {
		return children[i].name < children[j].name
	})
	for _, child := range children {
		html += childHTML(&child)
	}
	return html
}
