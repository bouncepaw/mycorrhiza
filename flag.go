package main

import (
	"bufio"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/internal/version"

	"golang.org/x/term"
)

// CLI options are read and parsed here.

// printHelp prints the help message.
func printHelp() {
	_, _ = fmt.Fprintf(
		flag.CommandLine.Output(),
		"Usage: %s WIKI_PATH\n",
		os.Args[0],
	)
	flag.PrintDefaults()
}

// parseCliArgs parses CLI options and sets several important global variables. Call it early.
func parseCliArgs() error {
	var createAdminName string
	var versionFlag bool

	flag.StringVar(&cfg.ListenAddr, "listen-addr", "", "Address to listen on. For example, 127.0.0.1:1737 or /run/mycorrhiza.sock.")
	flag.StringVar(&createAdminName, "create-admin", "", "Create a new admin. The password will be prompted in the terminal.")
	flag.BoolVar(&versionFlag, "version", false, "Print version information and exit.")
	flag.Usage = printHelp
	flag.Parse()

	if versionFlag {
		slog.Info("Running Mycorrhiza Wiki", "version", version.Long)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		slog.Error("Pass a wiki directory")
		return errors.New("wiki directory not passed")
	}

	wikiDir, err := filepath.Abs(args[0])
	if err != nil {
		slog.Error("Failed to take absolute filepath of wiki directory",
			"path", args[0], "err", err)
		return err
	}

	cfg.WikiDir = wikiDir

	if createAdminName != "" {
		if err := createAdminCommand(createAdminName); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	return nil
}

func createAdminCommand(name string) error {
	if err := files.PrepareWikiRoot(); err != nil {
		slog.Error("Failed to prepare wiki root", "err", err)
		return err
	}
	cfg.UseAuth = true
	cfg.AllowRegistration = true
	user.InitUserDatabase()

	password, err := askPass("Password")
	if err != nil {
		slog.Error("Failed to prompt password", "err", err)
		return err
	}
	if err := user.Register(name, password, "admin", "local", true); err != nil {
		slog.Error("Failed to register admin", "err", err)
		return err
	}
	return nil
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
