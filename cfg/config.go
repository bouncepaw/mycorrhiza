package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	WikiDir        string
	TemplatesDir   string
	configJsonPath string

	// Default values that can be overriden in config.json
	Address           = "127.0.0.1:80"
	TitleEditTemplate = `Edit %s`
	TitleTemplate     = `%s`
	GenericErrorMsg   = `<b>Sorry, something went wrong</b>`
	SiteTitle         = `MycorrhizaWiki`
)

func InitConfig(wd string) bool {
	log.Println("WikiDir is", wd)
	WikiDir = wd
	TemplatesDir = filepath.Join(filepath.Dir(WikiDir), "templates")
	configJsonPath = filepath.Join(filepath.Dir(WikiDir), "config.json")

	if _, err := os.Stat(configJsonPath); os.IsNotExist(err) {
		log.Println("config.json not found, using default values")
	} else {
		log.Println("config.json found, overriding default values...")
		return readConfig()
	}

	return true
}

func readConfig() bool {
	configJsonContents, err := ioutil.ReadFile(configJsonPath)
	if err != nil {
		log.Fatal("Error when reading config.json:", err)
		return false
	}

	cfg := struct {
		Address        string `json:"address"`
		SiteTitle      string `json:"site-title"`
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
	SiteTitle = cfg.SiteTitle
	TitleEditTemplate = cfg.TitleTemplates.EditHypha
	TitleTemplate = cfg.TitleTemplates.ViewHypha

	return true
}
