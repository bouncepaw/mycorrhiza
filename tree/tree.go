package tree

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/bouncepaw/mycorrhiza/util"
)

type sibling struct {
	name        string
	hasChildren bool
}

func (s *sibling) checkThisChild(hyphaName string) {
	if !s.hasChildren && path.Dir(hyphaName) == s.name {
		s.hasChildren = true
	}
}

func (s *sibling) asHTML() string {
	class := "navitree__entry navitree__sibling"
	if s.hasChildren {
		class += " navitree__sibling_fertile navitree__entry_fertile"
	} else {
		class += " navitree__sibling_infertile navitree__entry_infertile"
	}
	return fmt.Sprintf(
		`<li class="%s"><a class="navitree__link" href="/hypha/%s">%s</a></li>`,
		class,
		s.name,
		util.BeautifulName(path.Base(s.name)),
	)
}

type mainFamilyMember struct {
	name     string
	children []*mainFamilyMember
}

func (m *mainFamilyMember) checkThisChild(hyphaName string) (adopted bool) {
	if path.Dir(hyphaName) == m.name {
		m.children = append(m.children, &mainFamilyMember{
			name:     hyphaName,
			children: make([]*mainFamilyMember, 0),
		})
		return true
	}
	return false
}

func (m *mainFamilyMember) asHTML() string {
	if len(m.children) == 0 {
		return fmt.Sprintf(`<li class="navitree__entry navitree__entry_infertile navitree__trunk navitree__trunk_infertile"><a class="navitree__link" href="/hypha/%s">%s</a></li>`, m.name, util.BeautifulName(path.Base(m.name)))
	}
	sort.Slice(m.children, func(i, j int) bool {
		return m.children[i].name < m.children[j].name
	})
	html := fmt.Sprintf(`<li class="navitree__entry navitree__entry_fertile navitree__trunk navitree__trunk_fertile"><a class="navitree__link" href="/hypha/%s">%s</a><ul>`, m.name, util.BeautifulName(path.Base(m.name)))
	for _, child := range m.children {
		html += child.asHTML()
	}
	return html + `</li></ul></li>`
}

func mainFamilyFromPool(hyphaName string, subhyphaePool map[string]bool) *mainFamilyMember {
	var (
		nestLevel = strings.Count(hyphaName, "/")
		adopted   = make([]*mainFamilyMember, 0)
	)
	for subhyphaName, _ := range subhyphaePool {
		subnestLevel := strings.Count(subhyphaName, "/")
		if subnestLevel-1 == nestLevel && path.Dir(subhyphaName) == hyphaName {
			delete(subhyphaePool, subhyphaName)
			adopted = append(adopted, mainFamilyFromPool(subhyphaName, subhyphaePool))
		}
	}
	return &mainFamilyMember{name: hyphaName, children: adopted}
}

// Tree generates a tree for `hyphaName` as html and returns next and previous hyphae if any.
func Tree(hyphaName string, hyphaIterator func(func(string))) (html, prev, next string) {
	var (
		// One of the siblings is the hypha with name `hyphaName`
		siblings      = findSiblings(hyphaName, hyphaIterator)
		subhyphaePool = make(map[string]bool)
		I             int
	)
	hyphaIterator(func(otherHyphaName string) {
		for _, s := range siblings {
			s.checkThisChild(otherHyphaName)
		}
		if strings.HasPrefix(otherHyphaName, hyphaName+"/") {
			subhyphaePool[otherHyphaName] = true
		}
	})
	for i, s := range siblings {
		if s.name == hyphaName {
			I = i
			break
		}
		html += s.asHTML()
	}
	html += mainFamilyFromPool(hyphaName, subhyphaePool).asHTML()
	for _, s := range siblings[I+1:] {
		html += s.asHTML()
	}
	if I != 0 {
		prev = siblings[I-1].name
	}
	if I != len(siblings)-1 {
		next = siblings[I+1].name
	}
	return fmt.Sprintf(`<ul class="navitree">%s</ul>`, html), prev, next
}

func findSiblings(hyphaName string, hyphaIterator func(func(string))) []*sibling {
	siblings := []*sibling{&sibling{name: hyphaName, hasChildren: true}}
	hyphaIterator(func(otherHyphaName string) {
		if path.Dir(hyphaName) == path.Dir(otherHyphaName) && hyphaName != otherHyphaName {
			siblings = append(siblings, &sibling{name: otherHyphaName, hasChildren: false})
		}
	})
	sort.Slice(siblings, func(i, j int) bool {
		return siblings[i].name < siblings[j].name
	})
	return siblings
}
