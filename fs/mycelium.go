package fs

import (
	"io/ioutil"
	"log"

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

func VerifyMycelia() {
	var (
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
	}
	if !mainPresent {
		log.Fatal("No `main` mycelium given in config.json")
	}
	if !systemPresent {
		log.Fatal("No `system` mycelium given in config.json")
	}
	log.Println("Mycelial dirs are present")
}
