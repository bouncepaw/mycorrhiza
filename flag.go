package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/term"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/version"
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

// parseCliArgs parses CLI options and sets several important global variables. Call it early.
func parseCliArgs() {
	var createAdminName string
	var versionFlag bool

	flag.StringVar(&cfg.ListenAddr, "listen-addr", "", "Address to listen on. For example, 127.0.0.1:1737 or /run/mycorrhiza.sock.")
	flag.StringVar(&createAdminName, "create-admin", "", "Create a new admin. The password will be prompted in the terminal.")
	flag.BoolVar(&versionFlag, "version", false, "Print version information and exit.")
	flag.Usage = printHelp
	flag.Parse()

	if versionFlag {
		fmt.Println("Mycorrhiza Wiki", version.Long)
		os.Exit(0)
	}

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
	if err := files.PrepareWikiRoot(); err != nil {
		log.Fatal(err)
	}
	cfg.UseAuth = true
	cfg.AllowRegistration = true
	user.InitUserDatabase()

	password, err := askPass("Password")
	if err != nil {
		log.Fatal(err)
	}
	if err := user.Register(name, password, "admin", "local", true); err != nil {
		log.Fatal(err)
	}
}

func askPass(prompt string) (string, error) {
	var password []byte
	var err error
	fd := int(os.Stdin.Fd())

	if term.IsTerminal(fd) {
		fmt.Printf("%s: ", prompt)
		password, err = term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		fmt.Println()
	} else {
		fmt.Fprintf(os.Stderr, "Warning: Reading password from stdin.\n")
		// TODO: the buffering messes up repeated calls to readPassword
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", io.ErrUnexpectedEOF
		}
		password = scanner.Bytes()

		if len(password) == 0 {
			return "", fmt.Errorf("zero length password")
		}
	}

	return string(password), nil
}
