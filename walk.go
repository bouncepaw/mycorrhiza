package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	hyphaPattern    = `[^\s\d:/?&\\][^:?&\\]*`
	hyphaUrl        = `/{hypha:` + hyphaPattern + `}`
	revisionPattern = `[\d]+`
	revQuery        = `{rev:` + revisionPattern + `}`
	metaJsonPattern = `meta\.json`
)

var (
	leadingInt = regexp.MustCompile(`^[-+]?\d+`)
)

// matchNameToEverything matches `name` to all filename patterns and returns 4 boolean results.
func matchNameToEverything(name string) (metaJsonM, hyphaM bool) {
	// simpleMatch reduces boilerplate. Errors are ignored because I trust my regex skills.
	simpleMatch := func(s string, p string) bool {
		m, _ := regexp.MatchString(p, s)
		return m
	}
	return simpleMatch(name, metaJsonPattern),
		simpleMatch(name, hyphaPattern)
}

// scanHyphaDir scans directory at `fullPath` and tells what it has found.
func scanHyphaDir(fullPath string) (valid bool, possibleSubhyphae []string, metaJsonPath string, err error) {
	nodes, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return // implicit return values
	}

	for _, node := range nodes {
		metaJsonM, hyphaM := matchNameToEverything(node.Name())
		switch {
		case hyphaM && node.IsDir():
			possibleSubhyphae = append(possibleSubhyphae, filepath.Join(fullPath, node.Name()))
		case metaJsonM && !node.IsDir():
			metaJsonPath = filepath.Join(fullPath, "meta.json")
			// Other nodes are ignored. It is not promised they will be ignored in future versions
		}
	}

	if metaJsonPath != "" {
		valid = true
	}
	return // implicit return values
}

// hyphaName gets name of a hypha by stripping path to the hypha in `fullPath`
func hyphaName(fullPath string) string {
	// {rootWikiDir}/{the name}
	return fullPath[len(rootWikiDir)+1:]
}

// recurFindHyphae recursively searches for hyphae in passed directory path.
func recurFindHyphae(fullPath string) map[string]*Hypha {
	hyphae := make(map[string]*Hypha)
	valid, possibleSubhyphae, metaJsonPath, err := scanHyphaDir(fullPath)
	if err != nil {
		return hyphae
	}

	// First, let's process subhyphae
	for _, possibleSubhypha := range possibleSubhyphae {
		for k, v := range recurFindHyphae(possibleSubhypha) {
			hyphae[k] = v
		}
	}

	// This folder is not a hypha itself, nothing to do here
	if !valid {
		return hyphae
	}

	// Template hypha struct. Other fields are default json values.
	h := Hypha{
		FullName:   hyphaName(fullPath),
		Path:       fullPath,
		parentName: filepath.Dir(hyphaName(fullPath)),
		// Children names are unknown now
	}

	metaJsonContents, err := ioutil.ReadFile(metaJsonPath)
	if err != nil {
		log.Printf("Error when reading `%s`; skipping", metaJsonPath)
		return hyphae
	}
	err = json.Unmarshal(metaJsonContents, &h)
	if err != nil {
		log.Printf("Error when unmarshaling `%s`; skipping", metaJsonPath)
		log.Println(err)
		return hyphae
	}

	// fill in rooted paths to content files and full names
	for idStr, rev := range h.Revisions {
		rev.FullName = filepath.Join(h.parentName, rev.ShortName)
		rev.Id, _ = strconv.Atoi(idStr)
		if rev.BinaryName != "" {
			rev.BinaryPath = filepath.Join(fullPath, rev.BinaryName)
		}
		rev.TextPath = filepath.Join(fullPath, rev.TextName)
	}

	// Now the hypha should be ok, gotta send structs
	hyphae[h.FullName] = &h
	return hyphae
}
