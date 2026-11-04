package catsvc

import (
	"log/slog"
	"os"
	"slices"
	"sort"
	"sync"

	"github.com/bouncepaw/mycorrhiza/internal/categories/domain"
	"github.com/bouncepaw/mycorrhiza/internal/categories/storage"
	"github.com/bouncepaw/mycorrhiza/util"
)

var (
	categoryToHyphae  = map[string]*categoryNode{}
	hyphaToCategories = map[string]*hyphaNode{}

	store *catstorage.Storage
)

// Init initializes the category system. Call it after the Structure is initialized. This function might terminate the program in case of a bad mood or filesystem faults.
func Init(s *catstorage.Storage) error {
	store = s
	record, err := store.LoadCategories()
	if err != nil {
		slog.Error("Failed to read categories from disk", "err", err)
		return err
	}

	for _, cat := range record {
		if len(cat.Hyphae) == 0 {
			continue
		}
		cat.Name = util.CanonicalName(cat.Name)
		for i, hyphaName := range cat.Hyphae {
			cat.Hyphae[i] = util.CanonicalName(hyphaName)
		}
		sort.Strings(cat.Hyphae)
		categoryToHyphae[cat.Name] = &categoryNode{hyphaList: cat.Hyphae}
	}

	for cat, hyphaeInCat := range categoryToHyphae {
		for _, hyphaName := range hyphaeInCat.hyphaList {
			if node, ok := hyphaToCategories[hyphaName]; ok {
				node.storeCategory(cat)
			} else {
				hyphaToCategories[hyphaName] = &hyphaNode{categoryList: []string{cat}}
			}
		}
	}

	slog.Info("Indexed categories", "n", len(categoryToHyphae))
	return nil
}

type categoryNode struct {
	hyphaList []string
}

func (cn *categoryNode) storeHypha(hypname string) {
	i, found := slices.BinarySearch(cn.hyphaList, hypname)
	if found {
		return
	}
	cn.hyphaList = slices.Insert(cn.hyphaList, i, hypname)
}

func (cn *categoryNode) removeHypha(hypname string) {
	i, found := slices.BinarySearch(cn.hyphaList, hypname)
	if !found {
		return
	}
	cn.hyphaList = slices.Delete(cn.hyphaList, i, i+1)
}

type hyphaNode struct {
	categoryList []string
}

// inserts sorted
func (hn *hyphaNode) storeCategory(cat string) {
	i, found := slices.BinarySearch(hn.categoryList, cat)
	if found {
		return
	}
	hn.categoryList = slices.Insert(hn.categoryList, i, cat)
}

func (hn *hyphaNode) removeCategory(cat string) {
	i, found := slices.BinarySearch(hn.categoryList, cat)
	if !found {
		return
	}
	hn.categoryList = slices.Delete(hn.categoryList, i, i+1)
}

var fileMutex sync.Mutex

func saveToDisk() {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	var record []catdomain.Category
	mutex.RLock()
	for name, node := range categoryToHyphae {
		record = append(record, catdomain.Category{
			Name:   name,
			Hyphae: slices.Clone(node.hyphaList),
		})
	}
	mutex.RUnlock()

	if err := store.SaveCategories(record); err != nil {
		slog.Error("Failed to write categories.json", "err", err)
		os.Exit(1)
	}
}
