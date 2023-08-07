// Package interwiki provides interwiki capabilities. Most of them, at least.
package interwiki

import (
	"encoding/json"
	"errors"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
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

func dropEmptyStrings(ss []string) (clean []string) {
	for _, s := range ss {
		if s != "" {
			clean = append(clean, s)
		}
	}
	return clean
}

// difference returns the elements in `a` that aren't in `b`.
// Taken from https://stackoverflow.com/a/45428032
// CC BY-SA 4.0, no changes made
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
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

func replaceEntry(oldWiki *Wiki, newWiki *Wiki) error {
	diff := difference(
		append(newWiki.Aliases, newWiki.Name),
		append(oldWiki.Aliases, oldWiki.Name),
	)
	if ok, name := areNamesFree(diff); !ok {
		return errors.New(name)
	}
	deleteEntry(oldWiki)
	return addEntry(newWiki)
}

func deleteEntry(wiki *Wiki) {
	mutex.Lock()
	defer mutex.Unlock()

	// I'm being fancy here. Come on, the code here is already a mess.
	// Let me have some fun.
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		names := append(wiki.Aliases, wiki.Name)
		for _, name := range names {
			name := name // I guess we need that
			delete(entriesByName, name)
		}
		wg.Done()
	}()

	go func() {
		for i, w := range listOfEntries {
			i, w := i, w
			if w.Name == wiki.Name {
				log.Println("It came to delete")
				// Drop ith element.
				listOfEntries[i] = listOfEntries[len(listOfEntries)-1]
				listOfEntries = listOfEntries[:len(listOfEntries)-1]
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

func addEntry(wiki *Wiki) error {
	mutex.Lock()
	defer mutex.Unlock()
	wiki.Aliases = dropEmptyStrings(wiki.Aliases)
	var (
		names    = append(wiki.Aliases, wiki.Name)
		ok, name = areNamesFree(names)
	)
	if !ok {
		log.Printf("There are multiple uses of the same name ‘%s’\n", name)
		return errors.New(name)
	}
	if len(names) == 0 {
		log.Println("No names passed for a new interwiki entry")
		// There is something clearly wrong with error-returning in this function.
		return errors.New("")
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
