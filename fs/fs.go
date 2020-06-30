package fs

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

type Storage struct {
	// hypha name => path
	paths map[string]string
	root  string
}

var Hs *Storage

// InitStorage initiates filesystem-based hypha storage. It has to be called after configuration was inited.
func InitStorage() {
	Hs = &Storage{
		paths: make(map[string]string),
		root:  cfg.WikiDir,
	}
	Hs.indexHyphae(Hs.root)
	log.Printf("Indexed %v hyphae\n", len(Hs.paths))
}

// hyphaName gets name of a hypha by stripping path to the hypha in `fullPath`
func hyphaName(fullPath string) string {
	// {cfg.WikiDir}/{the name}
	return fullPath[len(cfg.WikiDir)+1:]
}

// indexHyphae searches for all hyphae that seem valid in `path` and saves their absolute paths to `s.paths`. This function is recursive.
func (s *Storage) indexHyphae(path string) {
	nodes, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Error when checking", path, ":", err, "; skipping")
		return
	}

	for _, node := range nodes {
		matchesHypha, err := regexp.MatchString(cfg.HyphaPattern, node.Name())
		if err != nil {
			log.Fatal("Error when matching", node.Name(), err, "\n")
			return
		}
		switch name := filepath.Join(path, node.Name()); {
		case matchesHypha && node.IsDir():
			s.indexHyphae(name)
		case node.Name() == "meta.json" && !node.IsDir():
			s.paths[hyphaName(path)] = path
		}
	}
}
