// Package categories provides category management.
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

// For WithHypha and Contents, should the results be sorted?

// WithHypha returns what categories have the given hypha.
func WithHypha(hyphaName string) (categoryList []string) {
	panic("todo")
	return
}

// Contents returns what hyphae are in the category. If the returned slice is empty, the category does not exist, and vice versa.
func Contents(catName string) (hyphaList []string) {
	panic("todo")
	return
}

// AddHyphaToCategory adds the hypha to the category and updates the records on the disk. If the hypha is already in the category, nothing happens. This operation is async-safe.
func AddHyphaToCategory(hyphaName, catName string) {
	for _, cat := range WithHypha(hyphaName) {
		if cat == catName {
			return
		}
	}
	panic("todo")
}

// RemoveHyphaFromCategory removes the hypha from the category and updates the records on the disk. If the hypha is not in the category, nothing happens. This operation is async-safe.
func RemoveHyphaFromCategory(hyphaName, catName string) {
	panic("todo")
}
