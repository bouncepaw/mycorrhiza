// Package migration holds the utilities for migrating from older incompatible Mycomarkup versions.
//
// Migrations are meant to be removed couple of versions after being introduced.
//
// Available migrations:
//     * Rocket links
//     * Headings
package migration

import (
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
	"io"
	"log"
	"os"
	"strings"
)

func genericLineMigrator(
	commitMessage string,
	migrator func(string) string,
	commitErrorMessage string,
) {
	var (
		hop = history.
			Operation(history.TypeMarkupMigration).
			WithMsg(commitMessage).
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
			newText = migrator(oldText)
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
		log.Fatal(commitErrorMessage, hop.FirstErrorText())
	}

	log.Println("Migrated", len(mycoFiles), "Mycomarkup documents")
}
