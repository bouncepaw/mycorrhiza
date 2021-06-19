// Package files is used to get paths to different files Mycorrhiza uses. Also see cfg.
package files

import (
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

var paths struct {
	gitRepo                     string
	cacheDir                    string
	staticFiles                 string
	configPath                  string
	tokensJSON                  string
	registrationCredentialsJSON string
	fixedCredentialsJSON        string
}

// HyphaeDir returns the path to hyphae storage.
// A separate function is needed to easily know where a general storage path is
// needed rather than a concrete Git or the whole wiki storage path, so that we
// could easily refactor things later if we'll ever support different storages.
func HyphaeDir() string { return paths.gitRepo }

// GitRepo returns the path to the Git repository of the wiki.
func GitRepo() string { return paths.gitRepo }

// StaticFiles returns the path to static files directory
func StaticFiles() string { return paths.staticFiles }

// ConfigPath returns the path to the config file.
func ConfigPath() string { return paths.configPath }

// TokensJSON returns the path to the JSON user tokens storage.
func TokensJSON() string { return paths.tokensJSON }

// RegistrationCredentialsJSON returns the path to the JSON registration
// credentials storage.
func RegistrationCredentialsJSON() string { return paths.registrationCredentialsJSON }

// FixedCredentialsJSON returns the path to the JSON fixed credentials storage.
func FixedCredentialsJSON() string { return paths.fixedCredentialsJSON }

// PrepareWikiRoot ensures all needed directories and files exist and have
// correct permissions.
func PrepareWikiRoot() error {
	if err := os.MkdirAll(cfg.WikiDir, os.ModeDir|0777); err != nil {
		return err
	}

	paths.cacheDir = filepath.Join(cfg.WikiDir, "cache")
	if err := os.MkdirAll(paths.cacheDir, os.ModeDir|0777); err != nil {
		return err
	}

	paths.gitRepo = filepath.Join(cfg.WikiDir, "wiki.git")
	if err := os.MkdirAll(paths.gitRepo, os.ModeDir|0777); err != nil {
		return err
	}

	paths.staticFiles = filepath.Join(cfg.WikiDir, "static")
	if err := os.MkdirAll(paths.staticFiles, os.ModeDir|0777); err != nil {
		return err
	}

	paths.configPath = filepath.Join(cfg.WikiDir, "config.ini")

	paths.tokensJSON = filepath.Join(paths.cacheDir, "tokens.json")
	paths.fixedCredentialsJSON = filepath.Join(cfg.WikiDir, "fixed-users.json")
	paths.registrationCredentialsJSON = filepath.Join(paths.cacheDir, "registered-users.json")

	return nil
}
