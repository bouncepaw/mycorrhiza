package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/plugin/lang"
)

type MyceliumConfig struct {
	Names []string `json:"names"`
	Type  string   `json:"type"`
}

const (
	HyphaPattern    = `[^\s\d:/?&\\][^:?&\\]*`
	HyphaUrl        = `/{hypha:` + HyphaPattern + `}`
	RevisionPattern = `[\d]+`
	RevQuery        = `{rev:` + RevisionPattern + `}`
	MyceliumPattern = `[^\s\d:/?&\\][^:?&\\/]*`
	MyceliumUrl     = `/:{mycelium:` + MyceliumPattern + `}`
)

var (
	Locale         map[string]string
	WikiDir        string
	configJsonPath string

	// Default values that can be overriden in config.json
	Address           = "0.0.0.0:80"
	LocaleName        = "en"
	SiteTitle         = `MycorrhizaWiki`
	Theme             = `default-light`
	HomePage          = `/Home`
	BinaryLimit int64 = 10 << 20
	Mycelia           = []MyceliumConfig{
		{[]string{"main"}, "main"},
		{[]string{"sys", "system"}, "system"},
	}
)

func InitConfig(wd string) bool {
	log.Println("WikiDir is", wd)
	WikiDir = wd
	configJsonPath = filepath.Join(WikiDir, "config.json")

	if _, err := os.Stat(configJsonPath); os.IsNotExist(err) {
		log.Println("config.json not found, using default values")
		return false
	}
	log.Println("config.json found, overriding default values...")
	return readConfig()
}

func readConfig() bool {
	configJsonContents, err := ioutil.ReadFile(configJsonPath)
	if err != nil {
		log.Fatal("Error when reading config.json:", err)
		return false
	}

	cfg := struct {
		Address       string           `json:"address"`
		Theme         string           `json:"theme"`
		SiteTitle     string           `json:"site-title"`
		HomePage      string           `json:"home-page"`
		BinaryLimitMB int64            `json:"binary-limit-mb"`
		LocaleName    string           `json:"locale"`
		Mycelia       []MyceliumConfig `json:"mycelia"`
	}{}

	err = json.Unmarshal(configJsonContents, &cfg)
	if err != nil {
		log.Fatal("Error when parsing config.json:", err)
		return false
	}

	Address = cfg.Address
	Theme = cfg.Theme
	SiteTitle = cfg.SiteTitle
	HomePage = "/" + cfg.HomePage
	BinaryLimit = 1024 * cfg.BinaryLimitMB
	Mycelia = cfg.Mycelia

	switch cfg.LocaleName {
	case "en":
		Locale = lang.EnglishMap
	default:
		Locale = lang.EnglishMap
	}

	return true
}
