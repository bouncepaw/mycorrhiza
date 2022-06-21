// Package interwiki provides interwiki capabilities. Most of them, at least.
package interwiki

import (
	"encoding/json"
	"errors"
	"github.com/bouncepaw/mycomarkup/v5/options"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/util"
	"log"
	"os"
	"sync"
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
		if err := addEntry(&wiki); err != nil {
			log.Fatalln(err.Error())
		}
	}
	log.Printf("Loaded %d interwiki entries\n", len(listOfEntries))
}

func areNamesFree(names []string) (bool, string) {
	for _, name := range names {
		if _, found := entriesByName[name]; found {
			return false, name
		}
	}
	return true, ""
}

var mutex sync.Mutex

func addEntry(wiki *Wiki) error {
	mutex.Lock()
	defer mutex.Unlock()

	var (
		names    = append(wiki.Aliases, wiki.Name)
		ok, name = areNamesFree(names)
	)
	if !ok {
		log.Printf("There are multiple uses of the same name ‘%s’\n", name)
		return errors.New(name)
	}

	listOfEntries = append(listOfEntries, wiki)
	for _, name := range names {
		entriesByName[name] = wiki
	}
	return nil
}

func HrefLinkFormatFor(prefix string) (string, options.InterwikiError) {
	prefix = util.CanonicalName(prefix)
	if wiki, ok := entriesByName[prefix]; ok {
		return wiki.LinkHrefFormat, options.Ok
	}
	return "", options.UnknownPrefix
}

func ImgSrcFormatFor(prefix string) (string, options.InterwikiError) {
	prefix = util.CanonicalName(prefix)
	if wiki, ok := entriesByName[prefix]; ok {
		return wiki.ImgSrcFormat, options.Ok
	}
	return "", options.UnknownPrefix
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

func saveInterwikiJson() {
	// Trust me, wiki crashing when an admin takes an administrative action totally makes sense.
	if data, err := json.MarshalIndent(listOfEntries, "", "\t"); err != nil {
		log.Fatalln(err)
	} else if err = os.WriteFile(files.InterwikiJSON(), data, 0666); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Saved interwiki.json")
	}
}
