//go:generate qtc -dir=views
//go:generate qtc -dir=tree
// Command mycorrhiza is a program that runs a mycorrhiza wiki.
package main

import (
	"log"
	"net/http"
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

	log.Println("Running Mycorrhiza Wiki 1.2.0")
	if err := os.Chdir(files.HyphaeDir()); err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki directory is", cfg.WikiDir)
	log.Println("Using Git storage at", files.HyphaeDir())

	// Init the subsystems:
	hyphae.Index(files.HyphaeDir())
	user.InitUserDatabase()
	history.Start()
	history.InitGitRepo()
	shroom.SetHeaderLinks()

	// Static files:
	static.InitFS(files.StaticFiles())

	// Network:
	web.Init()
	log.Fatal(http.ListenAndServe("0.0.0.0:"+cfg.HTTPPort, nil))
}
