// Package files is used to get paths to different files Mycorrhiza uses. Also see cfg.
package files

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/mitchellh/go-homedir"
)

var paths struct {
	tokensJSON                  string
	registrationCredentialsJSON string
	fixedCredentialsJSON        string
}

// TokensJSON returns a path to the JSON file where users' tokens are stored.
//
// Default path: $XDG_DATA_HOME/mycorrhiza/tokens.json
func TokensJSON() string { return paths.tokensJSON }

// RegistrationCredentialsJSON returns a path to the JSON file where registration credentials are stored.
//
// Default path: $XDG_DATA_HOME/mycorrhiza/registration.json
func RegistrationCredentialsJSON() string { return paths.registrationCredentialsJSON }

// FixedCredentialsJSON returns a path to the JSON file where fixed credentials are stored.
//
// There is no default path.
func FixedCredentialsJSON() string { return paths.fixedCredentialsJSON }

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
		paths.fixedCredentialsJSON = dir
	}

	return nil
}

func tokenStoragePath() (string, error) {
	dir, err := xdg.DataFile("mycorrhiza/tokens.json")
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(dir, cfg.WikiDir) {
		return "", errors.New("wiki storage directory includes private config files")
	}
	return dir, nil
}

func registrationCredentialsPath() (string, error) {
	var err error
	path := cfg.RegistrationCredentialsPath

	if len(path) == 0 {
		path, err = xdg.DataFile("mycorrhiza/registration.json")
		if err != nil {
			return "", fmt.Errorf("cannot get a file to registration credentials, so no registered users will be saved: %w", err)
		}
	} else {
		path, err = homedir.Expand(path)
		if err != nil {
			return "", fmt.Errorf("cannot expand RegistrationCredentialsPath: %w", err)
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("cannot expand RegistrationCredentialsPath: %w", err)
		}
	}

	return path, nil
}

func fixedCredentialsPath() (string, error) {
	var err error
	path := cfg.FixedAuthCredentialsPath

	if len(path) > 0 {
		path, err = homedir.Expand(path)
		if err != nil {
			return "", fmt.Errorf("cannot expand FixedAuthCredentialsPath: %w", err)
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("cannot expand FixedAuthCredentialsPath: %w", err)
		}
	}
	return path, nil
}
