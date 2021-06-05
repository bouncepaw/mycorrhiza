package main

import (
	"flag"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"log"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/assets"
)

// CLI options are read and parsed here.

var printExampleConfig bool

func init() {
	flag.StringVar(&cfg.ConfigFilePath, "config-path", "", "Path to a configuration file. Leave empty if you don't want to use it.")
	flag.BoolVar(&printExampleConfig, "print-example-config", false, "If true, print an example configuration file contents and exit. You can save the output to a file and base your own configuration on it.")
	flag.Usage = printHelp
}

// printHelp prints the help message. The help message is stored in assets.
func printHelp() {
	_, err := fmt.Fprintf(
		flag.CommandLine.Output(),
		assets.HelpMessage(),
		os.Args[0],
	)
	if err != nil {
		log.Fatal(err)
	}
	flag.PrintDefaults()
}

// parseCliArgs parses CLI options and sets several important global variables. Call it early.
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

	wikiDir, err := filepath.Abs(args[0])
	cfg.WikiDir = wikiDir
	if err != nil {
		log.Fatal(err)
	}
}
