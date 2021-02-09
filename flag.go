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
	flag.StringVar(&util.UserHypha, "user-hypha", "u", "Hypha which is a superhypha of all user pages")
	flag.StringVar(&util.AuthMethod, "auth-method", "none", "What auth method to use. Variants: \"none\", \"fixed\"")
	flag.StringVar(&util.FixedCredentialsPath, "fixed-credentials-path", "mycocredentials.json", "Used when -auth-method=fixed. Path to file with user credentials.")
	flag.StringVar(&util.HeaderLinksHypha, "header-links-hypha", "", "Optional hypha that overrides the header links")
	flag.StringVar(&util.GeminiCertPath, "gemini-cert-path", "", "Directory where you store Gemini certificates. Leave empty if you don't want to use Gemini.")
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

	util.HomePage = CanonicalName(util.HomePage)
	util.UserHypha = CanonicalName(util.UserHypha)
	util.HeaderLinksHypha = CanonicalName(util.HeaderLinksHypha)

	switch util.AuthMethod {
	case "none":
	case "fixed":
		user.AuthUsed = true
		user.ReadUsersFromFilesystem()
	default:
		log.Fatal("Error: unknown auth method:", util.AuthMethod)
	}
}
