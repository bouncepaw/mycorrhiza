// Package migration holds the utilities for migrating from older incompatible Mycomarkup versions.
//
// As of, there is rocket link migration only. Migrations are meant to be removed couple of versions after being introduced.
package migration

import (
	"github.com/bouncepaw/mycomarkup/v4/tools"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
)

// TODO: add heading migration too.

var rocketMarkerPath string

// MigrateRocketsMaybe checks if the rocket link migration marker exists. If it exists, nothing is done. If it does not, the migration takes place.
//
// This function writes logs and might terminate the program. Tons of side-effects, stay safe.
func MigrateRocketsMaybe() {
	rocketMarkerPath = files.FileInRoot(".mycomarkup-rocket-link-migration-marker.txt")
	if !shouldMigrateRockets() {
		return
	}

	var (
		hop = history.
			Operation(history.TypeMarkupMigration).
			WithMsg("Migrate rocket links to the new syntax").
			WithUser(user.WikimindUser())
		mycoFiles = []string{}
	)

	for hypha := range hyphae.FilterHyphaeWithText(hyphae.YieldExistingHyphae()) {
		/// Open file, read from file, modify file. If anything goes wrong, scream and shout.

		file, err := os.OpenFile(hypha.TextFilePath(), os.O_RDWR, 0766)
		if err != nil {
			hop.WithErrAbort(err)
			log.Fatal("Something went wrong when opening ", hypha.TextFilePath(), ": ", err.Error())
		}

		var buf strings.Builder
		_, err = io.Copy(&buf, file)
		if err != nil {
			hop.WithErrAbort(err)
			_ = file.Close()
			log.Fatal("Something went wrong when reading ", hypha.TextFilePath(), ": ", err.Error())
		}

		var (
			oldText = buf.String()
			newText = tools.MigrateRocketLinks(oldText)
		)
		if oldText != newText { // This file right here is being migrated for real.
			mycoFiles = append(mycoFiles, hypha.TextFilePath())

			err = file.Truncate(0)
			if err != nil {
				hop.WithErrAbort(err)
				_ = file.Close()
				log.Fatal("Something went wrong when truncating ", hypha.TextFilePath(), ": ", err.Error())
			}

			_, err = file.Seek(0, 0)
			if err != nil {
				hop.WithErrAbort(err)
				_ = file.Close()
				log.Fatal("Something went wrong when seeking in  ", hypha.TextFilePath(), ": ", err.Error())
			}

			_, err = file.WriteString(newText)
			if err != nil {
				hop.WithErrAbort(err)
				_ = file.Close()
				log.Fatal("Something went wrong when writing to ", hypha.TextFilePath(), ": ", err.Error())
			}
		}
		_ = file.Close()
	}

	if len(mycoFiles) == 0 {
		hop.Abort()
		return
	}

	if hop.WithFiles(mycoFiles...).Apply().HasErrors() {
		log.Fatal("Something went wrong when commiting rocket link migration: ", hop.FirstErrorText())
	}
	log.Println("Migrated", len(mycoFiles), "Mycomarkup documents")
	createRocketLinkMarker()
}

func shouldMigrateRockets() bool {
	file, err := os.Open(rocketMarkerPath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		log.Fatalln("When checking if rocket migration is needed:", err.Error())
	}
	_ = file.Close()
	return false
}

func createRocketLinkMarker() {
	err := ioutil.WriteFile(
		rocketMarkerPath,
		[]byte(`This file is used to mark that the rocket link migration was made successfully. If this file is deleted, the migration might happen again depending on the version. You should probably not touch this file at all and let it be.`),
		0766,
	)
	if err != nil {
		log.Fatalln(err)
	}
}
