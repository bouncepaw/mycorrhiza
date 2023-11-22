// Package history provides a git wrapper.
package history

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Path to git executable. Set at init()
var gitpath string

var renameMsgPattern = regexp.MustCompile(`^Rename ‘(.*)’ to ‘.*’`)

var gitEnv = []string{"GIT_COMMITTER_NAME=wikimind", "GIT_COMMITTER_EMAIL=wikimind@mycorrhiza"}

// Start finds git and initializes git credentials.
func Start() {
	path, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Could not find the git executable. Check your $PATH.")
	}
	gitpath = path
}

// InitGitRepo checks a Git repository and initializes it if necessary.
func InitGitRepo() {
	// Detect if the Git repo directory is a Git repository
	isGitRepo := true
	buf, err := silentGitsh("rev-parse", "--git-dir")
	if err != nil {
		isGitRepo = false
	}
	if isGitRepo {
		gitDir := buf.String()
		if filepath.IsAbs(gitDir) && !filepath.HasPrefix(gitDir, files.HyphaeDir()) {
			isGitRepo = false
		}
	}
	if !isGitRepo {
		log.Println("Initializing Git repo at", files.HyphaeDir())
		gitsh("init")
		gitsh("config", "core.quotePath", "false")
	}
}

// I pronounce it as [gɪt͡ʃ].
// gitsh is async-safe, therefore all other git-related functions in this module are too.
func gitsh(args ...string) (out bytes.Buffer, err error) {
	fmt.Printf("$ %v\n", args)
	cmd := exec.Command(gitpath, args...)
	cmd.Dir = files.HyphaeDir()
	cmd.Env = append(cmd.Environ(), gitEnv...)

	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("gitsh:", err)
	}
	return *bytes.NewBuffer(b), err
}

// silentGitsh is like gitsh, except it writes less to the stdout.
func silentGitsh(args ...string) (out bytes.Buffer, err error) {
	cmd := exec.Command(gitpath, args...)
	cmd.Dir = files.HyphaeDir()
	cmd.Env = append(cmd.Environ(), gitEnv...)

	b, err := cmd.CombinedOutput()
	return *bytes.NewBuffer(b), err
}

// Rename renames from `from` to `to` using `git mv`.
func Rename(from, to string) error {
	log.Println(util.ShorterPath(from), util.ShorterPath(to))
	_, err := gitsh("mv", "--force", from, to)
	return err
}
