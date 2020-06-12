package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
)

func scanHyphaDir(fullPath string) (structureMet bool, possibleRevisionPaths []string, possibleHyphaPaths []string, reterr error) {
	nodes, err := ioutil.ReadDir(fullPath)
	if err != nil {
		reterr = err
		return // implicit return values
	}

	var (
		mmJsonPresent  bool
		zeroDirPresent bool
	)

	for _, node := range nodes {
		matchedHypha, _ := regexp.MatchString(hyphaPattern, node.Name())
		matchedRev, _ := regexp.MatchString(revisionPattern, node.Name())
		switch {
		case matchedRev && node.IsDir():
			if node.Name() == "0" {
				zeroDirPresent = true
			}
			possibleRevisionPaths = append(
				possibleRevisionPaths,
				filepath.Join(fullPath, node.Name()),
			)
		case (node.Name() == "mm.json") && !node.IsDir():
			mmJsonPresent = true
		case matchedHypha && node.IsDir():
			possibleHyphaPaths = append(
				possibleHyphaPaths,
				filepath.Join(fullPath, node.Name()),
			)
			// Other nodes are ignored. It is not promised they will be ignored in future versions
		}
	}

	if mmJsonPresent && zeroDirPresent {
		structureMet = true
	}

	return // implicit return values
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Hypha name is rootWikiDir/{here}
func hyphaName(fullPath string) string {
	return fullPath[len(rootWikiDir)+1:]
}

const (
	hyphaPattern    = `[^\s\d:/?&\\][^:?&\\]*`
	revisionPattern = `[\d]+`
)

// Sends found hyphae to the `ch`. `fullPath` is tested for hyphaness, then its subdirs with hyphaesque names are tested too using goroutines for each subdir. The function is recursive.
func recurFindHyphae(fullPath string) (hyphae []*Hypha) {

	structureMet, possibleRevisionPaths, possibleHyphaPaths, err := scanHyphaDir(fullPath)
	if err != nil {
		return hyphae
	}

	// First, let's process inner hyphae
	for _, possibleHyphaPath := range possibleHyphaPaths {
		hyphae = append(hyphae, recurFindHyphae(possibleHyphaPath)...)
	}

	// This folder is not a hypha itself, nothing to do here
	if !structureMet {
		return hyphae
	}

	// Template hypha struct. Other fields are default jsont values.
	h := Hypha{
		Path:       fullPath,
		Name:       hyphaName(fullPath),
		ParentName: filepath.Dir(hyphaName(fullPath)),
		// Children names are unknown now
	}

	// Fill in every revision
	for _, possibleRevisionPath := range possibleRevisionPaths {
		rev, err := makeRevision(possibleRevisionPath)
		if err == nil {
			h.Revisions = append(h.Revisions, rev)
		}
	}

	mmJsonPath := filepath.Join(fullPath, "mm.json")
	mmJsonContents, err := ioutil.ReadFile(mmJsonPath)
	if err != nil {
		fmt.Println(fullPath, ">\tError:", err)
		return hyphae
	}
	err = json.Unmarshal(mmJsonContents, &h)
	if err != nil {
		fmt.Println(fullPath, ">\tError:", err)
		return hyphae
	}

	// Now the hypha should be ok, gotta send structs
	hyphae = append(hyphae, &h)
	return hyphae
}

func makeRevision(fullPath string) (r Revision, err error) {
	// fullPath is expected to be a path to a dir.
	// Revision directory must have at least `m.json` and `t.txt` files.
	var (
		mJsonPresent bool
		tTxtPresent  bool
		bPresent     bool
	)

	nodes, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return r, err
	}

	for _, node := range nodes {
		if node.IsDir() {
			continue
		}
		switch node.Name() {
		case "m.json":
			mJsonPresent = true
		case "t.txt":
			tTxtPresent = true
		case "b":
			bPresent = true
		}
	}

	if !(mJsonPresent && tTxtPresent) {
		return r, errors.New("makeRevision: m.json and t.txt files are not found")
	}

	// If all the flags are true, this directory is assumed to be a revision. Gotta check further. This is template Revision struct. Other fields fall back to default init values.
	mJsonPath := filepath.Join(fullPath, "m.json")
	mJsonContents, err := ioutil.ReadFile(mJsonPath)
	if err != nil {
		fmt.Println(fullPath, ">\tError:", err)
		return r, err
	}

	r = Revision{}
	err = json.Unmarshal(mJsonContents, &r)
	if err != nil {
		fmt.Println(fullPath, ">\tError:", err)
		return r, err
	}

	// Now, let's fill in t.txt path
	r.TextPath = filepath.Join(fullPath, "t.txt")

	// There's sense in reading binary file only if the hypha is marked as such
	if r.MimeType != "application/x-hypha" {
		// Do not check for binary file presence, attempt to read it will fail anyway
		if bPresent {
			r.BinaryPath = filepath.Join(fullPath, "b")
		} else {
			return r, errors.New("makeRevision: b file not present")
		}
	}

	// So far, so good. Let's fill in id. It is guaranteed to be correct, so no error checking
	id, _ := strconv.Atoi(filepath.Base(fullPath))
	r.Id = id

	// It is safe now to return, I guess
	return r, nil
}
