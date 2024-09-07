package categories

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"
	"sort"
	"sync"

	"github.com/bouncepaw/mycorrhiza/internal/files"
	"github.com/bouncepaw/mycorrhiza/util"
)

var categoryToHyphae = map[string]*categoryNode{}
var hyphaToCategories = map[string]*hyphaNode{}

// Init initializes the category system. Call it after the Structure is initialized. This function might terminate the program in case of a bad mood or filesystem faults.
func Init() error {
	record, err := readCategoriesFromDisk()
	if err != nil {
		slog.Error("Failed to read categories from disk", "err", err)
		return err
	}

	for _, cat := range record.Categories {
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

type catFileRecord struct {
	Categories []catRecord `json:"categories"`
}

type catRecord struct {
	Name   string   `json:"name"`
	Hyphae []string `json:"hyphae"`
}

func readCategoriesFromDisk() (catFileRecord, error) {
	var (
		record            catFileRecord
		categoriesFile    = files.CategoriesJSON()
		fileContents, err = os.ReadFile(categoriesFile)
	)
	if os.IsNotExist(err) {
		return record, nil
	}
	if err != nil {
		return record, err
	}

	err = json.Unmarshal(fileContents, &record)
	if err != nil {
		return record, err
	}

	return record, nil
}

var fileMutex sync.Mutex

func saveToDisk() {
	var record catFileRecord
	for name, node := range categoryToHyphae {
		record.Categories = append(record.Categories, catRecord{
			Name:   name,
			Hyphae: node.hyphaList,
		})
	}
	data, err := json.MarshalIndent(record, "", "\t")
	if err != nil {
		slog.Error("Failed to marshal categories record", "err", err)
		os.Exit(1) // Better fail now, than later
	}

	// TODO: make the data safer somehow?? Back it up before overwriting?
	fileMutex.Lock()
	err = os.WriteFile(files.CategoriesJSON(), data, 0666)
	if err != nil {
		slog.Error("Failed to write categories.json", "err", err)
		os.Exit(1)
	}
	fileMutex.Unlock()
}
