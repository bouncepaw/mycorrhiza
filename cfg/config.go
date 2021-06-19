// Package cfg contains global variables that represent the current wiki configuration, including CLI options, configuration file values and header links.
package cfg

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/go-ini/ini"
)

// These variables represent the configuration. You are not meant to modify them after they were set.
//
// See https://mycorrhiza.lesarbr.es/hypha/configuration/fields for their docs.
var (
	WikiName      string
	NaviTitleIcon string

	HomeHypha        string
	UserHypha        string
	HeaderLinksHypha string

	HTTPPort              string
	URL                   string
	GeminiCertificatePath string

	UseFixedAuth                bool
	UseRegistration             bool
	LimitRegistration           int

	OmnipresentScripts []string
	ViewScripts        []string
	EditScripts        []string
)

// These variables are set before reading the config file, they are set in main.parseCliArgs.
var (
	// WikiDir is a full path to the wiki storage directory, which also must be a git repo.
	WikiDir string
	// ConfigFilePath is a path to the config file. Its value is used when calling ReadConfigFile.
	ConfigFilePath string
)

// Config represents a Mycorrhiza wiki configuration file. This type is used only when reading configs.
type Config struct {
	WikiName      string
	NaviTitleIcon string
	Hyphae
	Network
	Authorization
	CustomScripts
}

// Hyphae is a section of Config which has fields related to special hyphae.
type Hyphae struct {
	HomeHypha        string
	UserHypha        string
	HeaderLinksHypha string
}

// Network is a section of Config that has fields related to network stuff: HTTP and Gemini.
type Network struct {
	HTTPPort              uint64
	URL                   string
	GeminiCertificatePath string
}

// CustomScripts is a section with paths to JavaScript files that are loaded on specified pages.
type CustomScripts struct {
	// OmnipresentScripts: everywhere...
	OmnipresentScripts []string `delim:","`
	// ViewScripts: /hypha, /rev
	ViewScripts []string `delim:","`
	// Edit: /edit
	EditScripts []string `delim:","`
}

// Authorization is a section of Config that has fields related to authorization and authentication.
type Authorization struct {
	UseFixedAuth             bool

	UseRegistration             bool
	LimitRegistration           uint64
}

// ReadConfigFile reads a config on the given path and stores the configuration. Call it sometime during the initialization.
//
// Note that it may log.Fatal.
func ReadConfigFile() {
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
			UseFixedAuth:             false,

			UseRegistration:             false,
			LimitRegistration:           0,
		},
		CustomScripts: CustomScripts{
			OmnipresentScripts: []string{},
			ViewScripts:        []string{},
			EditScripts:        []string{},
		},
	}

	if ConfigFilePath != "" {
		path, err := filepath.Abs(ConfigFilePath)
		if err != nil {
			log.Fatalf("cannot expand config file path: %s", err)
		}

		log.Println("Loading config at", path)
		err = ini.MapTo(cfg, path)
		if err != nil {
			log.Fatal(err)
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
}
