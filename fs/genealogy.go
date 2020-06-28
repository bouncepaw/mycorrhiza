package fs

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

func (s *Storage) RenderHypha(h *Hypha) {
}

// If Name == "", the tree is empty.
type Tree struct {
	Name        string
	Ancestors   []string
	Siblings    []string
	Descendants []*Tree
	Root        bool
}

// GetTree generates a Tree for the given hypha name.
// It can also generate trees for non-existent hyphae, that's why we use `name string` instead of making it a method on `Hypha`.
// In `root` is `false`, siblings will not be fetched.
func (s *Storage) GetTree(name string, root bool) *Tree {
	t := &Tree{Name: name, Root: root}
	for hyphaName, _ := range s.paths {
		s.compareNamesAndAppend(t, hyphaName)
	}
	sort.Slice(t.Ancestors, func(i, j int) bool {
		return strings.Count(t.Ancestors[i], "/") < strings.Count(t.Ancestors[j], "/")
	})
	sort.Strings(t.Siblings)
	sort.Slice(t.Descendants, func(i, j int) bool {
		a := t.Descendants[i].Name
		b := t.Descendants[j].Name
		return len(a) < len(b)
	})
	log.Printf("Generate tree for %v: %v %v\n", t.Name, t.Ancestors, t.Siblings)
	return t
}

// Compares names appends name2 to an array of `t`:
func (s *Storage) compareNamesAndAppend(t *Tree, name2 string) {
	switch {
	case t.Name == name2:
	case strings.HasPrefix(t.Name, name2):
		t.Ancestors = append(t.Ancestors, name2)
	case t.Root && (strings.Count(t.Name, "/") == strings.Count(name2, "/") &&
		(filepath.Dir(t.Name) == filepath.Dir(name2))):
		t.Siblings = append(t.Siblings, name2)
	case strings.HasPrefix(name2, t.Name):
		t.Descendants = append(t.Descendants, s.GetTree(name2, false))
	}
}

// asHtml returns HTML representation of a tree.
// It recursively itself on the tree's children.
// TODO: redo with templates. I'm not in mood for it now.
func (t *Tree) AsHtml() (html string) {
	if t.Name == "" {
		return ""
	}
	html += `<ul class="navitree__node">`
	if t.Root {
		for _, ancestor := range t.Ancestors {
			html += navitreeEntry(ancestor, "navitree__ancestor")
		}
		for _, siblingName := range t.Siblings {
			html += navitreeEntry(siblingName, "navitree__sibling")
		}
		html += navitreeEntry(t.Name, "navitree__pagename")
	} else {
		html += navitreeEntry(t.Name, "navitree__name")
	}

	for _, subtree := range t.Descendants {
		html += subtree.AsHtml()
	}

	html += `</ul>`
	return html
}

// navitreeEntry is a small utility function that makes generating html easier.
// Someone please redo it in templates.
func navitreeEntry(name, class string) string {
	return fmt.Sprintf(`<li class="navitree__entry %s">
	<a class="navitree__link" href="/%s">%s</a>
</li>
`, class, name, filepath.Base(name))
}
