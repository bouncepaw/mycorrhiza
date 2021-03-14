package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/assets"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

var printExampleConfig bool

func init() {
	// flag.StringVar(&util.URL, "url", "http://0.0.0.0:$port", "URL at which your wiki can be found. Used to generate feeds and social media previews")
	// flag.StringVar(&util.ServerPort, "port", "1737", "Port to serve the wiki at using HTTP")
	// flag.StringVar(&util.HomePage, "home", "home", "The home page name")
	// flag.StringVar(&util.SiteNavIcon, "icon", "🍄", "What to show in the navititle in the beginning, before the colon")
	// flag.StringVar(&util.SiteName, "name", "wiki", "What is the name of your wiki")
	// flag.StringVar(&util.UserHypha, "user-hypha", "u", "Hypha which is a superhypha of all user pages")
	// flag.StringVar(&util.AuthMethod, "auth-method", "none", "What auth method to use. Variants: \"none\", \"fixed\"")
	// flag.StringVar(&util.FixedCredentialsPath, "fixed-credentials-path", "mycocredentials.json", "Used when -auth-method=fixed. Path to file with user credentials.")
	// flag.StringVar(&util.HeaderLinksHypha, "header-links-hypha", "", "Optional hypha that overrides the header links")
	// flag.StringVar(&util.GeminiCertPath, "gemini-cert-path", "", "Directory where you store Gemini certificates. Leave empty if you don't want to use Gemini.")
	flag.StringVar(&util.ConfigFilePath, "config-path", "", "Path to a configuration file. Leave empty if you don't want to use it.")
	flag.BoolVar(&printExampleConfig, "print-example-config", false, "If true, print an example configuration file contents and exit. You can save the output to a file and base your own configuration on it.")
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			assets.HelpMessage(),
			os.Args[0],
		)
		flag.PrintDefaults()
	}
}

// Do the things related to cli args and die maybe
func parseCliArgs() {
	flag.Parse()

	args := flag.Args()
	if printExampleConfig {
		fmt.Printf(assets.ExampleConfig())
		os.Exit(0)
	}

	if len(args) == 0 {
		log.Fatal("Error: pass a wiki directory")
	}

	// It is ok if the path is ""
	util.ReadConfigFile(util.ConfigFilePath)

	var err error
	WikiDir, err = filepath.Abs(args[0])
	util.WikiDir = WikiDir
	if err != nil {
		log.Fatal(err)
	}

	if util.URL == "" {
		util.URL = "http://0.0.0.0:" + util.ServerPort
	}

	util.HomePage = util.CanonicalName(util.HomePage)
	util.UserHypha = util.CanonicalName(util.UserHypha)
	util.HeaderLinksHypha = util.CanonicalName(util.HeaderLinksHypha)
	user.AuthUsed = util.UseFixedAuth
	if user.AuthUsed && util.FixedCredentialsPath != "" {
		user.ReadUsersFromFilesystem()
	}
}
