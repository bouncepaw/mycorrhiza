package hyphae

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/mimetype"
)

// Index finds all hypha files in the full `path` and saves them to the hypha storage.
func Index(path string) {
	byNames = make(map[string]Hypher)
	ch := make(chan Hypher, 5)

	go func(ch chan Hypher) {
		indexHelper(path, 0, ch)
		close(ch)
	}(ch)

	for nh := range ch {
		switch oh := ByName(nh.CanonicalName()).(type) {
		case *EmptyHypha:
			Insert(nh)
		default:
			// In case of conflicts the newer hypha overwrites the previous
			switch nh, oh := nh.(*NonEmptyHypha), oh.(*NonEmptyHypha); {
			case (nh.Kind() == HyphaText) && (oh.Kind() == HyphaMedia):
				oh.TextPath = nh.TextPartPath()

			case (nh.Kind() == HyphaText) && (oh.Kind() == HyphaText):
				log.Printf("File collision for hypha ‘%s’, using %s rather than %s\n", nh.CanonicalName(), nh.TextPartPath(), oh.TextPartPath())
				oh.TextPath = nh.TextPartPath()

			case (nh.Kind() == HyphaMedia) && (oh.Kind() == HyphaMedia):
				log.Printf("File collision for hypha ‘%s’, using %s rather than %s\n", nh.CanonicalName(), nh.BinaryPath(), oh.BinaryPath())
				oh.SetBinaryPath(nh.BinaryPath())

			case (nh.Kind() == HyphaMedia) && (oh.Kind() == HyphaText):
				oh.SetBinaryPath(nh.BinaryPath())
			}
		}
	}
	log.Println("Indexed", Count(), "hyphae")
}

// indexHelper finds all hypha files in the full `path` and sends them to the
// channel. Handling of duplicate entries and attachment and counting them is
// up to the caller.
func indexHelper(path string, nestLevel uint, ch chan Hypher) {
	nodes, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range nodes {
		// If this hypha looks like it can be a hypha path, go deeper. Do not
		// touch the .git folders for it has an administrative importance!
		if node.IsDir() && IsValidName(node.Name()) && node.Name() != ".git" {
			indexHelper(filepath.Join(path, node.Name()), nestLevel+1, ch)
			continue
		}

		var (
			hyphaPartPath           = filepath.Join(path, node.Name())
			hyphaName, isText, skip = mimetype.DataFromFilename(hyphaPartPath)
			hypha                   = &NonEmptyHypha{name: hyphaName, Exists: true}
		)
		if !skip {
			if isText {
				hypha.TextPath = hyphaPartPath
			} else {
				hypha.SetBinaryPath(hyphaPartPath)
			}
			ch <- hypha
		}
	}
}
