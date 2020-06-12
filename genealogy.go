/* Genealogy is all about relationships between hyphae. For now, the only goal of this file is to help find children of hyphae as they are not marked during the hypha search phase.
 */
package main

type Genealogy struct {
	parent string
	child  string
}

func setRelations(hyphae map[string]*Hypha) {
	for name, h := range hyphae {
		if _, ok := hyphae[h.ParentName]; ok && h.ParentName != "." {
			hyphae[h.ParentName].ChildrenNames = append(hyphae[h.ParentName].ChildrenNames, name)
		}
	}
}
