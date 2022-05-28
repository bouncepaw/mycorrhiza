// Package interwiki provides interwiki capabilities. Most of them, at least.
package interwiki

import (
	"encoding/json"
	"github.com/bouncepaw/mycorrhiza/files"
	"log"
	"os"
)

func Init() {
	var (
		record, err = readInterwiki()
	)
	if err != nil {
		log.Fatalln(err)
	}
	for _, wiki := range record {
		wiki := wiki // This line is required
		wiki.canonize()
		theMap.list = append(theMap.list, &wiki)
		for _, prefix := range wiki.Names {
			if _, found := theMap.byName[prefix]; found {
				log.Fatalf("There are multiple uses of the same prefix ‘%s’\n", prefix)
			} else {
				theMap.byName[prefix] = &wiki
			}
		}
	}
	log.Printf("Loaded %d interwiki entries\n", len(theMap.list))
}

func HrefLinkFormatFor(prefix string) string {
	if wiki, ok := theMap.byName[prefix]; ok {
		return wiki.LinkHrefFormat
	}
	return "{NAME}" // TODO: error
}

func ImgSrcFormatFor(prefix string) string {
	if wiki, ok := theMap.byName[prefix]; ok {
		return wiki.ImgSrcFormat
	}
	return "{NAME}" // TODO: error
}

func readInterwiki() ([]Wiki, error) {
	var (
		record            []Wiki
		fileContents, err = os.ReadFile(files.InterwikiJSON())
	)
	if os.IsNotExist(err) {
		return record, nil
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileContents, &record)
	if err != nil {
		return nil, err
	}
	return record, nil
}
