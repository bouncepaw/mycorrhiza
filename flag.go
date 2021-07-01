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

func init() {
	flag.StringVar(&cfg.HTTPPort, "port", "", "Listen on another port. This option also updates the config file for your convenience.")
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
	if len(args) == 0 {
		log.Fatal("Error: pass a wiki directory")
	}

	wikiDir, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	cfg.WikiDir = wikiDir
}
