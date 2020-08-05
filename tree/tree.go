package tree

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

// If Name == "", the tree is empty.
type tree struct {
	name          string
	siblings      []string
	descendants   []*tree
	root          bool
	hyphaIterator func(func(string))
}

// TreeAsHtml generates a tree for `hyphaName`. `hyphaStorage` has this type because package `tree` has no access to `HyphaData` data type. One day it shall have it, I guess.
func TreeAsHtml(hyphaName string, hyphaIterator func(func(string))) string {
	t := &tree{name: hyphaName, root: true, hyphaIterator: hyphaIterator}
	t.fill()
	return t.asHtml()
}

// subtree adds a descendant tree to `t` and returns that tree.
func (t *tree) fork(descendantName string) *tree {
	subt := &tree{
		name:          descendantName,
		root:          false,
		hyphaIterator: t.hyphaIterator,
	}
	t.descendants = append(t.descendants, subt)
	return subt
}

// Compares names and does something with them, may generate a subtree.
func (t *tree) compareNamesAndAppend(name2 string) {
	switch {
	case t.name == name2:
	case t.root && path.Dir(t.name) == path.Dir(name2):
		t.siblings = append(t.siblings, name2)
	case t.name == path.Dir(name2):
		t.fork(name2).fill()
	}
}

// Fills t.siblings and t.descendants, sorts them and does the same to the descendants.
func (t *tree) fill() {
	t.hyphaIterator(func(hyphaName string) {
		t.compareNamesAndAppend(hyphaName)
	})
	sort.Strings(t.siblings)
	sort.Slice(t.descendants, func(i, j int) bool {
		return t.descendants[i].name < t.descendants[j].name
	})
}

// asHtml returns HTML representation of a tree.
// It applies itself recursively on the tree's children.
func (t *tree) asHtml() (html string) {
	if t.root {
		for _, siblingName := range t.siblings {
			html += navitreeEntry(siblingName, "navitree__sibling")
		}
		html += navitreeEntry(t.name, "navitree__pagename")
	} else {
		html += navitreeEntry(t.name, "navitree__name")
	}

	for _, subtree := range t.descendants {
		html += subtree.asHtml()
	}

	return `<ul class="navitree__node">` + html + `</ul>`
}

// Strip hypha name from all ancestor names, replace _ with spaces, title case
func beautifulName(uglyName string) string {
	return strings.Title(strings.ReplaceAll(path.Base(uglyName), "_", " "))
}

// navitreeEntry is a small utility function that makes generating html easier.
func navitreeEntry(name, class string) string {
	return fmt.Sprintf(`<li class="navitree__entry %s">
	<a class="navitree__link" href="/page/%s">%s</a>
</li>
`, class, name, beautifulName(name))
}
