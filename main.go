// Command mycorrhiza is a program that runs a mycorrhiza wiki.
//
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=history
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=mycoopts
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=auth
package main

import (
	"log"
	"os"

	"github.com/bouncepaw/mycorrhiza/categories"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
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
	parseCliArgs()

	if err := files.PrepareWikiRoot(); err != nil {
		log.Fatal(err)
	}

	if err := cfg.ReadConfigFile(files.ConfigPath()); err != nil {
		log.Fatal(err)
	}

	log.Println("Running Mycorrhiza Wiki", version.Short)
	if err := os.Chdir(files.HyphaeDir()); err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki directory is", cfg.WikiDir)

	// Init the subsystems:
	viewutil.Init()
	hyphae.Index(files.HyphaeDir())
	backlinks.IndexBacklinks()
	go backlinks.RunBacklinksConveyor()
	user.InitUserDatabase()
	history.Start()
	history.InitGitRepo()
	migration.MigrateRocketsMaybe()
	migration.MigrateHeadingsMaybe()
	shroom.SetHeaderLinks()
	categories.Init()
	interwiki.Init()

	// Static files:
	static.InitFS(files.StaticFiles())

	if !user.HasAnyAdmins() {
		log.Println("Your wiki has no admin yet. Run Mycorrhiza with -create-admin <username> option to create an admin.")
	}

	serveHTTP(web.Handler())
}
