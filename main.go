package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var rootWikiDir string

func hyphaeAsMap(hyphae []*Hypha) map[string]*Hypha {
	mh := make(map[string]*Hypha)
	for _, h := range hyphae {
		mh[h.Name] = h
	}
	return mh
}

func main() {
	if len(os.Args) == 1 {
		panic("Expected a root wiki pages directory")
	}
	// Required so the rootWikiDir hereinbefore does not get redefined.
	var err error
	rootWikiDir, err = filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	hyphae := hyphaeAsMap(recurFindHyphae(rootWikiDir))
	setRelations(hyphae)

	fmt.Println(hyphae)
}
