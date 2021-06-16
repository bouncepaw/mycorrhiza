//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=assets
//go:generate qtc -dir=views
//go:generate qtc -dir=tree
// Command mycorrhiza is a program that runs a mycorrhiza wiki.
package main

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/static"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/web"
	"log"
	"net/http"
	"os"
)

func main() {
	parseCliArgs()

	// It is ok if the path is ""
	cfg.ReadConfigFile()

	if err := files.CalculatePaths(); err != nil {
		log.Fatal(err)
	}

	log.Println("Running Mycorrhiza Wiki 1.2.0 indev")
	if err := os.Chdir(cfg.WikiDir); err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki storage directory is", cfg.WikiDir)

	// Init the subsystems:
	hyphae.Index(cfg.WikiDir)
	user.InitUserDatabase()
	history.Start()
	shroom.SetHeaderLinks()

	// Static files:
	static.InitFS(cfg.WikiDir + "/static")

	// Network:
	go handleGemini()
	web.Init()
	log.Fatal(http.ListenAndServe("0.0.0.0:"+cfg.HTTPPort, nil))
}
