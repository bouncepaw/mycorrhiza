package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	flag.StringVar(&util.URL, "url", "http://0.0.0.0:$port", "URL at which your wiki can be found. Used to generate feeds and social media previews")
	flag.StringVar(&util.ServerPort, "port", "1737", "Port to serve the wiki at using HTTP")
	flag.StringVar(&util.HomePage, "home", "home", "The home page name")
	flag.StringVar(&util.SiteNavIcon, "icon", "🍄", "What to show in the navititle in the beginning, before the colon")
	flag.StringVar(&util.SiteName, "name", "wiki", "What is the name of your wiki")
	flag.StringVar(&util.UserTree, "user-tree", "u", "Hypha which is a superhypha of all user pages")
	flag.StringVar(&util.AuthMethod, "auth-method", "none", "What auth method to use. Variants: \"none\", \"fixed\"")
	flag.StringVar(&util.FixedCredentialsPath, "fixed-credentials-path", "mycocredentials.json", "Used when -auth-method=fixed. Path to file with user credentials.")
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

	if util.URL == "http://0.0.0.0:$port" {
		util.URL = "http://0.0.0.0:" + util.ServerPort
	}

	if !isCanonicalName(util.HomePage) {
		log.Fatal("Error: you must use a proper name for the homepage")
	}

	if !isCanonicalName(util.UserTree) {
		log.Fatal("Error: you must use a proper name for user tree")
	}

	switch util.AuthMethod {
	case "none":
	case "fixed":
		user.AuthUsed = true
		user.ReadUsersFromFilesystem()
	default:
		log.Fatal("Error: unknown auth method:", util.AuthMethod)
	}
}
