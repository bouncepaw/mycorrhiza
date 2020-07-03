package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	DescribeHyphaHerePattern = "Describe %s here"
	WikiDir                  string
	configJsonPath           string

	// Default values that can be overriden in config.json
	Address           = "127.0.0.1:80"
	TitleEditTemplate = `Edit %s`
	TitleTemplate     = `%s`
	GenericErrorMsg   = `<b>Sorry, something went wrong</b>`
	SiteTitle         = `MycorrhizaWiki`
	Theme             = `default-light`
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
		Address        string           `json:"address"`
		Theme          string           `json:"theme"`
		SiteTitle      string           `json:"site-title"`
		Mycelia        []MyceliumConfig `json:"mycelia"`
		TitleTemplates struct {
			EditHypha string `json:"edit-hypha"`
			ViewHypha string `json:"view-hypha"`
		} `json:"title-templates"`
	}{}

	err = json.Unmarshal(configJsonContents, &cfg)
	if err != nil {
		log.Fatal("Error when parsing config.json:", err)
		return false
	}

	Address = cfg.Address
	Theme = cfg.Theme
	SiteTitle = cfg.SiteTitle
	TitleEditTemplate = cfg.TitleTemplates.EditHypha
	TitleTemplate = cfg.TitleTemplates.ViewHypha
	Mycelia = cfg.Mycelia

	return true
}
