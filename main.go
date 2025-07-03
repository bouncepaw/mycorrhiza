// Command mycorrhiza is a program that runs a mycorrhiza wiki.
//
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=mycoopts
package main

import (
	"log/slog"
	"os"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/categories"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/migration"
	"github.com/bouncepaw/mycorrhiza/internal/shroom"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/internal/version"
	"github.com/bouncepaw/mycorrhiza/interwiki"
	"github.com/bouncepaw/mycorrhiza/web"
	"github.com/bouncepaw/mycorrhiza/web/static"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

func main() {
	if err := parseCliArgs(); err != nil {
		os.Exit(1)
	}

	if err := files.PrepareWikiRoot(); err != nil {
		slog.Error("Failed to prepare wiki root", "err", err)
		os.Exit(1)
	}

	if err := cfg.ReadConfigFile(files.ConfigPath()); err != nil {
		slog.Error("Failed to read config", "err", err)
		os.Exit(1)
	}

	if err := os.Chdir(files.HyphaeDir()); err != nil {
		slog.Error("Failed to chdir to hyphae dir",
			"err", err, "hyphaeDir", files.HyphaeDir())
		os.Exit(1)
	}
	slog.Info("Running Mycorrhiza Wiki",
		"version", version.Short, "wikiDir", cfg.WikiDir)

	// Init the subsystems:
	// TODO: keep all crashes in main rather than somewhere there
	viewutil.Init()
	hyphae.Index(files.HyphaeDir())
	backlinks.IndexBacklinks()
	go backlinks.RunBacklinksConveyor()
	user.InitUserDatabase()
	if err := history.Start(); err != nil {
		os.Exit(1)
	}
	history.InitGitRepo()
	migration.MigrateRocketsMaybe()
	migration.MigrateHeadingsMaybe()
	shroom.SetHeaderLinks()
	if err := categories.Init(); err != nil {
		os.Exit(1)
	}
	if err := interwiki.Init(); err != nil {
		os.Exit(1)
	}

	// Static files:
	static.InitFS(files.StaticFiles())

	if !user.HasAnyAdmins() {
		slog.Error("Your wiki has no admin yet. Run Mycorrhiza with -create-admin <username> option to create an admin.")
	}

	err := serveHTTP(web.Handler())
	if err != nil {
		os.Exit(1)
	}
}
