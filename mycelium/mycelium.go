package mycelium

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

var (
	MainMycelium   string
	SystemMycelium string
)

func gatherDirNames(path string) map[string]struct{} {
	res := make(map[string]struct{})
	nodes, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range nodes {
		if node.IsDir() {
			res[node.Name()] = struct{}{}
		}
	}
	return res
}

// Add values to the set. If a value is already there, return false.
func addInUniqueSet(set map[string]struct{}, names []string) bool {
	ok := true
	for _, name := range names {
		if _, present := set[name]; present {
			ok = false
		}
		set[name] = struct{}{}
	}
	return ok
}

func Init() {
	var (
		// Used to check if there are no duplicates
		foundNames    = make(map[string]struct{})
		dirs          = gatherDirNames(cfg.WikiDir)
		mainPresent   bool
		systemPresent bool
	)
	for _, mycelium := range cfg.Mycelia {
		switch mycelium.Type {
		case "main":
			mainPresent = true
			MainMycelium = mycelium.Names[0]
		case "system":
			systemPresent = true
			SystemMycelium = mycelium.Names[0]
		}
		// Check if there is a dir corresponding to the mycelium
		if _, ok := dirs[mycelium.Names[0]]; !ok {
			log.Fatal("No directory found for mycelium " + mycelium.Names[0])
		}
		// Confirm uniqueness of names
		if ok := addInUniqueSet(foundNames, mycelium.Names); !ok {
			log.Fatal("At least one name was used more than once for mycelia")
		}
	}
	if !mainPresent {
		log.Fatal("No `main` mycelium given in config.json")
	}
	if !systemPresent {
		log.Fatal("No `system` mycelium given in config.json")
	}
	log.Println("Mycelial dirs are present")
}

func NameWithMyceliumInMap(m map[string]string) (res string) {
	var (
		hyphaName, okH = m["hypha"]
		mycelName, okM = m["mycelium"]
	)
	log.Println(m)
	if !okH {
		// It will result in an error when trying to open a hypha with such name
		return ":::"
	}
	if okM {
		res = canonicalMycelium(mycelName)
	} else {
		res = MainMycelium
	}
	return res + "/" + hyphaName
}

func canonicalMycelium(name string) string {
	log.Println("Determining canonical mycelial name for", name)
	name = strings.ToLower(name)
	for _, mycel := range cfg.Mycelia {
		for _, mycelName := range mycel.Names {
			if mycelName == name {
				return mycel.Names[0]
			}
		}
	}
	// This is a nonexistent mycelium. Return a name that will trigger an error
	return ":error:"
}
