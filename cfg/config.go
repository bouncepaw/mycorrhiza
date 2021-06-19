// Package cfg contains global variables that represent the current wiki configuration, including CLI options, configuration file values and header links.
package cfg

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-ini/ini"
)

// These variables represent the configuration. You are not meant to modify
// them after they were set.
// See https://mycorrhiza.lesarbr.es/hypha/configuration/fields for the
// documentation.
var (
	WikiName      string
	NaviTitleIcon string

	HomeHypha        string
	UserHypha        string
	HeaderLinksHypha string

	HTTPPort              string
	URL                   string
	GeminiCertificatePath string

	UseFixedAuth      bool
	UseRegistration   bool
	LimitRegistration int

	OmnipresentScripts []string
	ViewScripts        []string
	EditScripts        []string
)

// WikiDir is a full path to the wiki storage directory, which also must be a
// git repo. This variable is set in parseCliArgs().
var WikiDir string

// Config represents a Mycorrhiza wiki configuration file. This type is used
// only when reading configs.
type Config struct {
	WikiName      string `comment:"This name appears in the header and on various pages."`
	NaviTitleIcon string `comment:"This icon is used in the breadcrumbs bar."`
	Hyphae
	Network
	Authorization `comment:""`
	CustomScripts `comment:"You can specify additional scripts to load on different kinds of pages, delimited by a comma ',' sign."`
}

// Hyphae is a section of Config which has fields related to special hyphae.
type Hyphae struct {
	HomeHypha        string `comment:"This hypha will be the main (index) page of your wiki, served on /."`
	UserHypha        string `comment:"This hypha is used as a prefix for user hyphae."`
	HeaderLinksHypha string `comment:"You can also specify a hypha to populate your own custom header links from."`
}

// Network is a section of Config that has fields related to network stuff:
// HTTP and Gemini.
type Network struct {
	HTTPPort              uint64
	URL                   string `comment:"Set your wiki's public URL here. It's used for OpenGraph generation and syndication feeds."`
	GeminiCertificatePath string `comment:"Gemini requires servers to use TLS for client connections. Specify your certificate's path here."`
}

// CustomScripts is a section with paths to JavaScript files that are loaded on
// specified pages.
type CustomScripts struct {
	// OmnipresentScripts: everywhere...
	OmnipresentScripts []string `delim:"," comment:"These scripts are loaded from anywhere."`
	// ViewScripts: /hypha, /rev
	ViewScripts []string `delim:"," comment:"These scripts are only loaded on view pages."`
	// Edit: /edit
	EditScripts []string `delim:"," comment:"These scripts are only loaded on the edit page."`
}

// Authorization is a section of Config that has fields related to
// authorization and authentication.
type Authorization struct {
	UseFixedAuth      bool
	UseRegistration   bool
	LimitRegistration uint64 `comment:"This field controls the maximum amount of allowed registrations."`
}

// ReadConfigFile reads a config on the given path and stores the
// configuration. Call it sometime during the initialization.
func ReadConfigFile(path string) error {
	cfg := &Config{
		WikiName:      "Mycorrhiza Wiki",
		NaviTitleIcon: "🍄",
		Hyphae: Hyphae{
			HomeHypha:        "home",
			UserHypha:        "u",
			HeaderLinksHypha: "",
		},
		Network: Network{
			HTTPPort:              1737,
			URL:                   "",
			GeminiCertificatePath: "",
		},
		Authorization: Authorization{
			UseFixedAuth:      false,
			UseRegistration:   false,
			LimitRegistration: 0,
		},
		CustomScripts: CustomScripts{
			OmnipresentScripts: []string{},
			ViewScripts:        []string{},
			EditScripts:        []string{},
		},
	}

	dirty := false

	f, err := ini.Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f = ini.Empty()
			dirty = true
		} else {
			return fmt.Errorf("Failed to open the config file: %w", err)
		}
	}

	// Map the config file to the config struct. It'll do nothing if the file
	// doesn't exist or is empty.
	f.MapTo(cfg)

	// Update the port if it's set externally and is different from what's in
	// the config file
	if HTTPPort != "" && HTTPPort != strconv.FormatUint(cfg.Network.HTTPPort, 10) {
		port, err := strconv.ParseUint(HTTPPort, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse the port from command-line arguments: %w", err)
		}

		cfg.Network.HTTPPort = port
		dirty = true
	}

	// Save changes, if there are any
	if dirty {
		err = f.ReflectFrom(cfg)
		if err != nil {
			return fmt.Errorf("Failed to serialize the config: %w", err)
		}

		// Disable key-value auto-aligning, but retain spaces around '=' sign
		ini.PrettyFormat = false
		ini.PrettyEqual = true
		if err = f.SaveTo(path); err != nil {
			return fmt.Errorf("Failed to save the config file: %w", err)
		}
	}

	// Map the struct to the global variables
	WikiName = cfg.WikiName
	NaviTitleIcon = cfg.NaviTitleIcon
	HomeHypha = cfg.HomeHypha
	UserHypha = cfg.UserHypha
	HeaderLinksHypha = cfg.HeaderLinksHypha
	HTTPPort = strconv.FormatUint(cfg.HTTPPort, 10)
	URL = cfg.URL
	GeminiCertificatePath = cfg.GeminiCertificatePath
	UseFixedAuth = cfg.UseFixedAuth
	UseRegistration = cfg.UseRegistration
	LimitRegistration = int(cfg.LimitRegistration)
	OmnipresentScripts = cfg.OmnipresentScripts
	ViewScripts = cfg.ViewScripts
	EditScripts = cfg.EditScripts

	// This URL makes much more sense.
	if URL == "" {
		URL = "http://0.0.0.0:" + HTTPPort
	}

	return nil
}
