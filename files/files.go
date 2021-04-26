package files

import (
	"errors"
	"github.com/adrg/xdg"
	"github.com/bouncepaw/mycorrhiza/util"
	"log"
	"path/filepath"
	"strings"
)

var paths struct {
	tokensJSON                  string
	registrationCredentialsJSON string
	fixedCredentialsJSON        string
	configINI                   string
}

func TokensJSON() string                  { return paths.tokensJSON }
func RegistrationCredentialsJSON() string { return paths.registrationCredentialsJSON }
func FixedCredentialsJSON() string        { return paths.fixedCredentialsJSON }
func ConfigINI() string                   { return paths.configINI }

// CalculatePaths looks for all external paths and stores them. Tries its best to find any errors. It is safe it to call it multiple times in order to save new paths.
func CalculatePaths() error {
	if dir, err := registrationCredentialsPath(); err != nil {
		return err
	} else {
		paths.registrationCredentialsJSON = dir
	}

	if dir, err := tokenStoragePath(); err != nil {
		return err
	} else {
		paths.tokensJSON = dir
	}

	if dir, err := fixedCredentialsPath(); err != nil {
		return err
	} else {
		paths.tokensJSON = dir
	}

	if dir, err := configPath(); err != nil {
		return err
	} else {
		paths.configINI = dir
	}

	return nil
}

func tokenStoragePath() (string, error) {
	dir, err := xdg.DataFile("mycorrhiza/tokens.json")
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(dir, util.WikiDir) {
		return "", errors.New("wiki storage directory includes private config files")
	}
	return dir, nil
}

func registrationCredentialsPath() (string, error) {
	path, err := filepath.Abs(util.RegistrationCredentialsPath)
	if err != nil {
		return "", nil
	}

	if path == "" {
		dir, err := xdg.DataFile("mycorrhiza/registration.json")
		if err != nil {
			log.Println("Error: cannot get a file to registration credentials, so no registered users will be saved:", err)
			return "", err
		}
		path = dir
	}
	return path, nil
}

func fixedCredentialsPath() (string, error) {
	path, err := filepath.Abs(util.FixedCredentialsPath)
	return path, err
}

func configPath() (string, error) {
	path, err := filepath.Abs(util.ConfigFilePath)
	return path, err
}
