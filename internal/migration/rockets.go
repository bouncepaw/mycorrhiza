package migration

import (
	"io/ioutil"
	"log/slog"
	"os"

	"github.com/bouncepaw/mycorrhiza/internal/files"

	"git.sr.ht/~bouncepaw/mycomarkup/v5/tools"
)

var rocketMarkerPath string

// MigrateRocketsMaybe checks if the rocket link migration marker exists. If it exists, nothing is done. If it does not, the migration takes place.
//
// This function writes logs and might terminate the program. Tons of side-effects, stay safe.
func MigrateRocketsMaybe() {
	rocketMarkerPath = files.FileInRoot(".mycomarkup-rocket-link-migration-marker.txt")
	if !shouldMigrateRockets() {
		return
	}

	genericLineMigrator(
		"Migrate rocket links to the new syntax",
		tools.MigrateRocketLinks,
		"Something went wrong when commiting rocket link migration: ",
	)
	createRocketLinkMarker()
}

func shouldMigrateRockets() bool {
	file, err := os.Open(rocketMarkerPath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		slog.Error("Failed to check if rocket migration is needed", "err", err)
		os.Exit(1)
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
		slog.Error("Failed to create rocket link migration marker")
		os.Exit(1)
	}
}
