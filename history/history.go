package history

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/util"
)

// Start initializes git credentials.
func Start(wikiDir string) {
	_, err := gitsh("config", "user.name", "wikimind")
	if err != nil {
		log.Fatal(err)
	}
	_, err = gitsh("config", "user.email", "wikimind@mycorrhiza")
	if err != nil {
		log.Fatal(err)
	}
}

// Revision represents a revision, duh. Hash is usually short. Username is extracted from email.
type Revision struct {
	Hash     string
	Username string
	Time     time.Time
	Message  string
}

// TimeString returns a human readable time representation.
func (rev Revision) TimeString() string {
	return rev.Time.Format(time.RFC822)
}

// HyphaeLinks returns a comma-separated list of hyphae that were affected by this revision as HTML string.
func (rev Revision) HyphaeLinks() (html string) {
	// diff-tree --no-commit-id --name-only -r
	var (
		// List of files affected by this revision, one per line.
		out, err = gitsh("diff-tree", "--no-commit-id", "--name-only", "-r", rev.Hash)
		// set is used to determine if a certain hypha has been already noted (hyphae are stored in 2 files at most).
		set       = make(map[string]bool)
		isNewName = func(hyphaName string) bool {
			if _, present := set[hyphaName]; present {
				return false
			} else {
				set[hyphaName] = true
				return true
			}
		}
	)
	if err != nil {
		return ""
	}
	for _, filename := range strings.Split(out.String(), "\n") {
		// If filename has an ampersand:
		if strings.IndexRune(filename, '.') >= 0 {
			// Remove ampersanded suffix from filename:
			ampersandPos := strings.LastIndexByte(filename, '.')
			hyphaName := string([]byte(filename)[0:ampersandPos]) // is it safe?
			if isNewName(hyphaName) {
				// Entries are separated by commas
				if len(set) > 1 {
					html += `<span aria-hidden="true">, </span>`
				}
				html += fmt.Sprintf(`<a href="/rev/%[1]s/%[2]s">%[2]s</a>`, rev.Hash, hyphaName)
			}
		}
	}
	return html
}

func (rev Revision) RecentChangesEntry() (html string) {
	return fmt.Sprintf(`
<li><time>%s</time></li>
<li>%s</li>
<li>%s</li>
<li>%s</li>
`, rev.TimeString(), rev.Hash, rev.HyphaeLinks(), rev.Message)
}

// Path to git executable. Set at init()
var gitpath string

func init() {
	path, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Cound not find the git executable. Check your $PATH.")
	} else {
		log.Println("Git path is", path)
	}
	gitpath = path

}

// I pronounce it as [gɪt͡ʃ].
// gitsh is async-safe, therefore all other git-related functions in this module are too.
func gitsh(args ...string) (out bytes.Buffer, err error) {
	fmt.Printf("$ %v\n", args)
	cmd := exec.Command(gitpath, args...)
	cmd.Dir = util.WikiDir

	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("gitsh:", err)
	}
	return *bytes.NewBuffer(b), err
}

// Convert a UNIX timestamp as string into a time. If nil is returned, it means that the timestamp could not be converted.
func unixTimestampAsTime(ts string) *time.Time {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil
	}
	tm := time.Unix(i, 0)
	return &tm
}

// Rename renames from `from` to `to` using `git mv`.
func Rename(from, to string) error {
	log.Println(util.ShorterPath(from), util.ShorterPath(to))
	_, err := gitsh("mv", "--force", from, to)
	return err
}
