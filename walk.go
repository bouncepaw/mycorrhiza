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

func matchNameToEverything(name string) (hyphaM bool, revTxtM bool, revBinM bool, metaJsonM bool) {
	simpleMatch := func(s string, p string) bool {
		m, _ := regexp.MatchString(p, s)
		return m
	}
	switch {
	case simpleMatch(name, revTxtPattern):
		revTxtM = true
	case simpleMatch(name, revBinPattern):
		revBinM = true
	case simpleMatch(name, metaJsonPattern):
		metaJsonM = true
	case simpleMatch(name, hyphaPattern):
		hyphaM = true
	}
	return
}

func stripLeadingInt(s string) string {
	return leadingInt.FindString(s)
}

func hyphaDirRevsValidate(dto map[string]map[string]string) (res bool) {
	for k, _ := range dto {
		switch k {
		case "0":
			delete(dto, "0")
		default:
			res = true
		}
	}
	return res
}

func scanHyphaDir(fullPath string) (valid bool, revs map[string]map[string]string, possibleSubhyphae []string, metaJsonPath string, err error) {
	revs = make(map[string]map[string]string)
	nodes, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return // implicit return values
	}

	for _, node := range nodes {
		hyphaM, revTxtM, revBinM, metaJsonM := matchNameToEverything(node.Name())
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

// Hypha name is rootWikiDir/{here}
func hyphaName(fullPath string) string {
	return fullPath[len(rootWikiDir)+1:]
}

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
