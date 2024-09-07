// Package migration holds the utilities for migrating from older incompatible Mycomarkup versions.
//
// Migrations are meant to be removed couple of versions after being introduced.
//
// Available migrations:
//   - Rocket links
//   - Headings
package migration

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/user"
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
			slog.Error("Failed to open text part file", "path", hypha.TextFilePath(), "err", err)
			os.Exit(1)
		}

		var buf strings.Builder
		_, err = io.Copy(&buf, file)
		if err != nil {
			hop.WithErrAbort(err)
			_ = file.Close()
			slog.Error("Failed to read text part file", "path", hypha.TextFilePath(), "err", err)
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
				slog.Error("Failed to truncate text part file", "path", hypha.TextFilePath(), "err", err)
				os.Exit(1)
			}

			_, err = file.Seek(0, 0)
			if err != nil {
				hop.WithErrAbort(err)
				_ = file.Close()
				slog.Error("Failed to seek in text part file", "path", hypha.TextFilePath(), "err", err)
				os.Exit(1)
			}

			_, err = file.WriteString(newText)
			if err != nil {
				hop.WithErrAbort(err)
				_ = file.Close()
				slog.Error("Failed to write to text part file", "path", hypha.TextFilePath(), "err", err)
				os.Exit(1)
			}
		}
		_ = file.Close()
	}

	if len(mycoFiles) == 0 {
		hop.Abort()
		return
	}

	if hop.WithFiles(mycoFiles...).Apply().HasErrors() {
		slog.Error(commitErrorMessage + hop.FirstErrorText())
	}

	slog.Info("Migrated Mycomarkup documents", "n", len(mycoFiles))
}
