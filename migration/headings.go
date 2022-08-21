package migration

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5/tools"
	"github.com/bouncepaw/mycorrhiza/files"
	"io/ioutil"
	"log"
	"os"
)

var headingMarkerPath string

func MigrateHeadingsMaybe() {
	headingMarkerPath = files.FileInRoot(".mycomarkup-heading-migration-marker.txt")
	if !shouldMigrateHeadings() {
		return
	}

	genericLineMigrator(
		"Migrate headings to the new syntax",
		tools.MigrateHeadings,
		"Something went wrong when commiting heading migration: ")
	createHeadingMarker()
}

func shouldMigrateHeadings() bool {
	file, err := os.Open(headingMarkerPath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		log.Fatalln("When checking if heading migration is needed:", err.Error())
	}
	_ = file.Close()
	return false
}

func createHeadingMarker() {
	err := ioutil.WriteFile(
		headingMarkerPath,
		[]byte(`This file is used to mark that the heading migration was successful. If this file is deleted, the migration might happen again depending on the version. You should probably not touch this file at all and let it be.`),
		0766,
	)
	if err != nil {
		log.Fatalln(err)
	}
}
