package hyphae

import (
	"github.com/bouncepaw/mycorrhiza/internal/mimetype"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

// Index finds all hypha files in the full `path` and saves them to the hypha storage.
func Index(path string) {
	byNames = make(map[string]ExistingHypha)
	ch := make(chan ExistingHypha, 5)

	go func(ch chan ExistingHypha) {
		indexHelper(path, 0, ch)
		close(ch)
	}(ch)

	for foundHypha := range ch {
		switch storedHypha := ByName(foundHypha.CanonicalName()).(type) {
		case *EmptyHypha:
			Insert(foundHypha)

		case *TextualHypha:
			switch foundHypha := foundHypha.(type) {
			case *TextualHypha: // conflict! overwrite
				storedHypha.mycoFilePath = foundHypha.mycoFilePath
				slog.Info("File collision",
					"hypha", foundHypha.CanonicalName(),
					"usingFile", foundHypha.TextFilePath(),
					"insteadOf", storedHypha.TextFilePath(),
				)
			case *MediaHypha: // no conflict
				Insert(ExtendTextualToMedia(storedHypha, foundHypha.mediaFilePath))
			}

		case *MediaHypha:
			switch foundHypha := foundHypha.(type) {
			case *TextualHypha: // no conflict
				storedHypha.mycoFilePath = foundHypha.mycoFilePath
			case *MediaHypha: // conflict! overwrite
				storedHypha.mediaFilePath = foundHypha.mediaFilePath

				slog.Info("File collision",
					"hypha", foundHypha.CanonicalName(),
					"usingFile", foundHypha.MediaFilePath(),
					"insteadOf", storedHypha.MediaFilePath(),
				)
			}
		}
	}
	log.Println("Indexed", Count(), "hyphae")
}

// indexHelper finds all hypha files in the full `path` and sends them to the
// channel. Handling of duplicate entries and media and counting them is
// up to the caller.
func indexHelper(path string, nestLevel uint, ch chan ExistingHypha) {
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
		)
		if !skip {
			if isText {
				ch <- &TextualHypha{
					canonicalName: hyphaName,
					mycoFilePath:  hyphaPartPath,
				}
			} else {
				ch <- &MediaHypha{
					canonicalName: hyphaName,
					mycoFilePath:  "",
					mediaFilePath: hyphaPartPath,
				}
			}
		}
	}
}
