package categories

import (
	"encoding/json"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/util"
	"log"
	"os"
)

var categoryToHyphae = map[string]*categoryNode{}
var hyphaToCategories = map[string]*hyphaNode{}

// InitCategories initializes the category system. Call it after the Structure is initialized. This function might terminate the program in case of a bad mood or filesystem faults.
func InitCategories() {
	var (
		record, err = readCategoriesFromDisk()
	)
	if err != nil {
		log.Fatalln(err)
	}

	for _, cat := range record.Categories {
		if len(cat.Hyphae) == 0 {
			continue
		}
		cat.Name = util.CanonicalName(cat.Name)
		for i, hyphaName := range cat.Hyphae {
			cat.Hyphae[i] = util.CanonicalName(hyphaName)
		}
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

	log.Println("Found", len(categoryToHyphae), "categories")
}

type categoryNode struct {
	// TODO: ensure this is sorted
	hyphaList []string
}

func (cn *categoryNode) storeHypha(hypname string) {
	for _, hyphaName := range cn.hyphaList {
		if hyphaName == hypname {
			return
		}
	}
	cn.hyphaList = append(cn.hyphaList, hypname)
}

func (cn *categoryNode) removeHypha(hypname string) {
	for i, hyphaName := range cn.hyphaList {
		if hyphaName == hypname {
			cn.hyphaList[i] = cn.hyphaList[len(cn.hyphaList)-1]
			cn.hyphaList = cn.hyphaList[:len(cn.hyphaList)-1]
		}
	}
}

type hyphaNode struct {
	// TODO: ensure this is sorted
	categoryList []string
}

func (hn *hyphaNode) storeCategory(cat string) {
	for _, category := range hn.categoryList {
		if category == cat {
			return
		}
	}
	hn.categoryList = append(hn.categoryList, cat)
}

func (hn *hyphaNode) removeCategory(cat string) {
	for i, category := range hn.categoryList {
		if category == cat {
			hn.categoryList[i] = hn.categoryList[len(hn.categoryList)-1]
			hn.categoryList = hn.categoryList[:len(hn.categoryList)-1]
		}
	}
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
