// Package categories provides category management. All operations in this package are mutexed.
//
// As per the long pondering, this is how categories (cats for short)
// work in Mycorrhiza:
//
//     - Cats are not hyphae. Cats are separate entities. This is not as
//       vibeful as I would have wanted, but seems to be more practical
//       due to //the reasons//.
//     - Cats are stored outside of git. Instead, they are stored in a
//       JSON file, path to which is determined by files.CategoriesJSON.
//     - Due to not being stored in git, no cat history is tracked, and
//       cat operations are not mentioned on the recent changes page.
//     - For cat A, if there are 0 hyphae in the cat, cat A does not
//       exist. If there are 1 or more hyphae in the cat, cat A exists.
package categories

import "sync"

// List returns names of all categories.
func List() (categoryList []string) {
	mutex.RLock()
	for cat, _ := range categoryToHyphae {
		categoryList = append(categoryList, cat)
	}
	mutex.RUnlock()
	return categoryList
}

// WithHypha returns what categories have the given hypha. The hypha name must be canonical.
func WithHypha(hyphaName string) (categoryList []string) {
	mutex.RLock()
	defer mutex.RUnlock()
	if node, ok := hyphaToCategories[hyphaName]; ok {
		return node.categoryList
	} else {
		return nil
	}
}

// Contents returns what hyphae are in the category. If the returned slice is empty, the category does not exist, and vice versa. The category name must be canonical.
func Contents(catName string) (hyphaList []string) {
	mutex.RLock()
	defer mutex.RUnlock()
	if node, ok := categoryToHyphae[catName]; ok {
		return node.hyphaList
	} else {
		return nil
	}
}

var mutex sync.RWMutex

// AddHyphaToCategory adds the hypha to the category and updates the records on the disk. If the hypha is already in the category, nothing happens. Pass canonical names.
func AddHyphaToCategory(hyphaName, catName string) {
	mutex.Lock()
	if node, ok := hyphaToCategories[hyphaName]; ok {
		node.storeCategory(catName)
	} else {
		hyphaToCategories[hyphaName] = &hyphaNode{categoryList: []string{catName}}
	}

	if node, ok := categoryToHyphae[catName]; ok {
		node.storeHypha(hyphaName)
	} else {
		categoryToHyphae[catName] = &categoryNode{hyphaList: []string{hyphaName}}
	}
	mutex.Unlock()
	go saveToDisk()
}

// RemoveHyphaFromCategory removes the hypha from the category and updates the records on the disk. If the hypha is not in the category, nothing happens. Pass canonical names.
func RemoveHyphaFromCategory(hyphaName, catName string) {
	mutex.Lock()
	if node, ok := hyphaToCategories[hyphaName]; ok {
		node.removeCategory(catName)
	}

	if node, ok := categoryToHyphae[catName]; ok {
		node.removeHypha(hyphaName)
	}
	mutex.Unlock()
	go saveToDisk()
}
