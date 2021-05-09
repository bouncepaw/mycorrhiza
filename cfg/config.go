package cfg

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/go-ini/ini"
)

var (
	WikiName      string
	NaviTitleIcon string

	HomeHypha        string
	UserHypha        string
	HeaderLinksHypha string

	HTTPPort              string
	URL                   string
	GeminiCertificatePath string

	WikiDir        string
	ConfigFilePath string

	UseFixedAuth                bool
	FixedAuthCredentialsPath    string
	UseRegistration             bool
	RegistrationCredentialsPath string
	LimitRegistration           int
)

// Config represents a Mycorrhiza wiki configuration file.
//
// See https://mycorrhiza.lesarbr.es/hypha/configuration/fields for fields' docs.
type Config struct {
	WikiName      string
	NaviTitleIcon string
	Hyphae
	Network
	Authorization
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

// Authorization is a section of Config that has fields related to authorization and authentication.
type Authorization struct {
	UseFixedAuth             bool
	FixedAuthCredentialsPath string

	UseRegistration             bool
	RegistrationCredentialsPath string
	LimitRegistration           uint64
}

// ReadConfigFile reads a config on the given path and stores the configuration.
func ReadConfigFile(path string) {
	cfg := &Config{
		WikiName:      "MycorrhizaWiki",
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
			FixedAuthCredentialsPath: "",

			UseRegistration:             false,
			RegistrationCredentialsPath: "",
			LimitRegistration:           0,
		},
	}

	if path != "" {
		path, err := filepath.Abs(path)
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
	FixedAuthCredentialsPath = cfg.FixedAuthCredentialsPath
	UseRegistration = cfg.UseRegistration
	RegistrationCredentialsPath = cfg.RegistrationCredentialsPath
	LimitRegistration = int(cfg.LimitRegistration)
}
