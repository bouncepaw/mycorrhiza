package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

const (
	hyphaPattern    = `[^\s\d:/?&\\][^:?&\\]*`
	hyphaUrl        = `/{hypha:` + hyphaPattern + `}`
	revisionPattern = `[\d]+`
	revQuery        = `{rev:` + revisionPattern + `}`
	revTxtPattern   = revisionPattern + `\.txt`
	revBinPattern   = revisionPattern + `\.bin`
	metaJsonPattern = `meta\.json`
)

var (
	leadingInt = regexp.MustCompile(`^[-+]?\d+`)
)

// matchNameToEverything matches `name` to all filename patterns and returns 4 boolean results.
func matchNameToEverything(name string) (revTxtM, revBinM, metaJsonM, hyphaM bool) {
	// simpleMatch reduces boilerplate. Errors are ignored because I trust my regex skills.
	simpleMatch := func(s string, p string) bool {
		m, _ := regexp.MatchString(p, s)
		return m
	}
	return simpleMatch(name, revTxtPattern),
		simpleMatch(name, revBinPattern),
		simpleMatch(name, metaJsonPattern),
		simpleMatch(name, hyphaPattern)
}

// stripLeadingInt finds number in the beginning of `s` and returns it.
func stripLeadingInt(s string) string {
	return leadingInt.FindString(s)
}

// hyphaDirRevsValidate checks if `dto` is ok.
// It also deletes pair with "0" as key so there is no revision with this id.
func hyphaDirRevsValidate(dto map[string]map[string]string) (res bool) {
	if _, ok := dto["0"]; ok {
		delete(dto, "0")
	}
	return len(dto) > 0
}

// scanHyphaDir scans directory at `fullPath` and tells what it has found.
func scanHyphaDir(fullPath string) (valid bool, revs map[string]map[string]string, possibleSubhyphae []string, metaJsonPath string, err error) {
	revs = make(map[string]map[string]string)
	nodes, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return // implicit return values
	}

	for _, node := range nodes {
		revTxtM, revBinM, metaJsonM, hyphaM := matchNameToEverything(node.Name())
		switch {
		case hyphaM && node.IsDir():
			possibleSubhyphae = append(possibleSubhyphae, filepath.Join(fullPath, node.Name()))
		case revTxtM && !node.IsDir():
			revId := stripLeadingInt(node.Name())
			if _, ok := revs[revId]; !ok {
				revs[revId] = make(map[string]string)
			}
			revs[revId]["txt"] = filepath.Join(fullPath, node.Name())
		case revBinM && !node.IsDir():
			revId := stripLeadingInt(node.Name())
			if _, ok := revs[revId]; !ok {
				revs[revId] = make(map[string]string)
			}
			revs[revId]["bin"] = filepath.Join(fullPath, node.Name())
		case metaJsonM && !node.IsDir():
			metaJsonPath = filepath.Join(fullPath, "meta.json")
			// Other nodes are ignored. It is not promised they will be ignored in future versions
		}
	}

	valid = hyphaDirRevsValidate(revs)
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
	valid, revs, possibleSubhyphae, metaJsonPath, err := scanHyphaDir(fullPath)
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
		return hyphae
	}

	// Fill in every revision paths
	for id, paths := range revs {
		if r, ok := h.Revisions[id]; ok {
			r.FullName = filepath.Join(h.parentName, r.ShortName)
			for fType, fPath := range paths {
				switch fType {
				case "bin":
					r.BinaryPath = fPath
				case "txt":
					r.TextPath = fPath
				}
			}
		} else {
			log.Printf("Error when reading hyphae from disk: hypha `%s`'s meta.json provided no information about revision `%s`, but files %s are provided; skipping\n", h.FullName, id, paths)
		}
	}

	// Now the hypha should be ok, gotta send structs
	hyphae[h.FullName] = &h
	return hyphae
}
