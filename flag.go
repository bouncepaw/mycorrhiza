package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/user"
	"golang.org/x/term"
)

// CLI options are read and parsed here.

// printHelp prints the help message.
func printHelp() {
	fmt.Fprintf(
		flag.CommandLine.Output(),
		"Usage: %s WIKI_PATH\n",
		os.Args[0],
	)
	flag.PrintDefaults()
}

type CreateUserCommand struct {
	name string
}

// parseCliArgs parses CLI options and sets several important global variables. Call it early.
func parseCliArgs() {
	var createAdminName string

	flag.StringVar(&cfg.ListenAddr, "listen-addr", "", "Address to listen on. For example, 127.0.0.1:1737 or /run/mycorrhiza.sock.")
	flag.StringVar(&createAdminName, "create-admin", "", "Create a new admin. The password will be prompted in the terminal.")
	flag.Usage = printHelp
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("error: pass a wiki directory")
	}

	wikiDir, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	cfg.WikiDir = wikiDir

	if createAdminName != "" {
		createAdminCommand(createAdminName)
		os.Exit(0)
	}
}

func createAdminCommand(name string) {
	wr := log.Writer()
	log.SetFlags(0)

	if err := files.PrepareWikiRoot(); err != nil {
		log.Fatal("error: ", err)
	}
	cfg.UseAuth = true
	cfg.AllowRegistration = true

	log.SetOutput(io.Discard)
	user.InitUserDatabase()
	log.SetOutput(wr)

	handle := int(syscall.Stdin)
	if !term.IsTerminal(handle) {
		log.Fatal("error: not a terminal")
	}

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(handle)
	fmt.Print("\n")
	if err != nil {
		log.Fatal("error: ", err)
	}

	password := string(passwordBytes)

	log.SetOutput(io.Discard)
	err = user.Register(name, password, "admin", "local", true)
	log.SetOutput(wr)

	if err != nil {
		log.Fatal("error: ", err)
	}
}
