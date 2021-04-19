package util

import (
	"log"
	"strconv"

	"github.com/go-ini/ini"
)

// See https://mycorrhiza.lesarbr.es/hypha/configuration/fields
type Config struct {
	WikiName      string
	NaviTitleIcon string
	Hyphae
	Network
	Authorization
}

type Hyphae struct {
	HomeHypha        string
	UserHypha        string
	HeaderLinksHypha string
}

type Network struct {
	HTTPPort              uint64
	URL                   string
	GeminiCertificatePath string
}

type Authorization struct {
	UseFixedAuth             bool
	FixedAuthCredentialsPath string

	UseRegistration             bool
	RegistrationCredentialsPath string
	LimitRegistration           uint64
}

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
		log.Println("Loading config at", path)
		err := ini.MapTo(cfg, path)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Map the struct to the global variables
	SiteName = cfg.WikiName
	SiteNavIcon = cfg.NaviTitleIcon
	HomePage = cfg.HomeHypha
	UserHypha = cfg.UserHypha
	HeaderLinksHypha = cfg.HeaderLinksHypha
	ServerPort = strconv.FormatUint(cfg.HTTPPort, 10)
	URL = cfg.URL
	GeminiCertPath = cfg.GeminiCertificatePath
	UseFixedAuth = cfg.UseFixedAuth
	FixedCredentialsPath = cfg.FixedAuthCredentialsPath
	UseRegistration = cfg.UseRegistration
	RegistrationCredentialsPath = cfg.RegistrationCredentialsPath
	LimitRegistration = int(cfg.LimitRegistration)
}
