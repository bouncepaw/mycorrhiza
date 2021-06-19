package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

// CLI options are read and parsed here.

//go:embed assets/config.ini
var defaultConfig []byte

var printExampleConfig bool

func init() {
	flag.StringVar(&cfg.ConfigFilePath, "config-path", "", "Path to a configuration file. Leave empty if you don't want to use it.")
	flag.BoolVar(&printExampleConfig, "print-example-config", false, "If true, print an example configuration file contents and exit. You can save the output to a file and base your own configuration on it.")
	flag.Usage = printHelp
}

// printHelp prints the help message.
func printHelp() {
	fmt.Fprintf(
		flag.CommandLine.Output(),
		"Usage: %s WIKI_PATH\n",
		os.Args[0],
	)
	flag.PrintDefaults()
}

// parseCliArgs parses CLI options and sets several important global variables. Call it early.
func parseCliArgs() {
	flag.Parse()

	args := flag.Args()
	if printExampleConfig {
		os.Stdout.Write(defaultConfig)
		os.Exit(0)
	}

	if len(args) == 0 {
		log.Fatal("Error: pass a wiki directory")
	}

	wikiDir, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	cfg.WikiDir = wikiDir
}
