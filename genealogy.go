/* Genealogy is all about relationships between hyphae. For now, the only goal of this file is to help find children of hyphae as they are not marked during the hypha search phase.

TODO: make use of family relations.
*/
package main

// setRelations fills in all children names based on what hyphae call their parents.
func setRelations(hyphae map[string]*Hypha) {
	for name, h := range hyphae {
		if _, ok := hyphae[h.parentName]; ok && h.parentName != "." {
			hyphae[h.parentName].ChildrenNames = append(hyphae[h.parentName].ChildrenNames, name)
		}
	}
}

// AddChild adds a name to the list of children names of the hypha.
func (h *Hypha) AddChild(childName string) {
	h.ChildrenNames = append(h.ChildrenNames, childName)
}
