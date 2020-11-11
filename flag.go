package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	flag.StringVar(&util.ServerPort, "port", "1737", "Port to serve the wiki at")
	flag.StringVar(&util.HomePage, "home", "home", "The home page")
	flag.StringVar(&util.SiteTitle, "title", "🍄", "How to call your wiki in the navititle")
}

// Do the things related to cli args and die maybe
func parseCliArgs() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Error: pass a wiki directory")
	}

	var err error
	WikiDir, err = filepath.Abs(args[0])
	util.WikiDir = WikiDir
	if err != nil {
		log.Fatal(err)
	}

	if !isCanonicalName(util.HomePage) {
		log.Fatal("Error: you must use a proper name for the homepage")
	}
}
