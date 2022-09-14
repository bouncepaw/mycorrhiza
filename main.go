// Command mycorrhiza is a program that runs a mycorrhiza wiki.
//
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=tree
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=history
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=mycoopts
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=auth
//go:generate go run github.com/valyala/quicktemplate/qtc -dir=hypview
package main

import (
	"github.com/bouncepaw/mycorrhiza/backlinks"
	"github.com/bouncepaw/mycorrhiza/categories"
	"github.com/bouncepaw/mycorrhiza/interwiki"
	"github.com/bouncepaw/mycorrhiza/migration"
	"github.com/bouncepaw/mycorrhiza/version"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"log"
	"os"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/static"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/web"
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

	serveHTTP(web.Handler())
}
