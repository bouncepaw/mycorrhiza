package history

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

var renameMsgPattern = regexp.MustCompile(`^Rename ‘(.*)’ to ‘.*’`)

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
	Hash              string
	Username          string
	Time              time.Time
	Message           string
	hyphaeAffectedBuf []string
}

// determine what hyphae were affected by this revision
func (rev *Revision) hyphaeAffected() (hyphae []string) {
	if nil != rev.hyphaeAffectedBuf {
		return rev.hyphaeAffectedBuf
	}
	hyphae = make([]string, 0)
	var (
		// List of files affected by this revision, one per line.
		out, err = gitsh("diff-tree", "--no-commit-id", "--name-only", "-r", rev.Hash)
		// set is used to determine if a certain hypha has been already noted (hyphae are stored in 2 files at most currently).
		set       = make(map[string]bool)
		isNewName = func(hyphaName string) bool {
			if _, present := set[hyphaName]; present {
				return false
			}
			set[hyphaName] = true
			return true
		}
	)
	if err != nil {
		return hyphae
	}
	for _, filename := range strings.Split(out.String(), "\n") {
		if strings.IndexRune(filename, '.') >= 0 {
			dotPos := strings.LastIndexByte(filename, '.')
			hyphaName := string([]byte(filename)[0:dotPos]) // is it safe?
			if isNewName(hyphaName) {
				hyphae = append(hyphae, hyphaName)
			}
		}
	}
	rev.hyphaeAffectedBuf = hyphae
	return hyphae
}

// TimeString returns a human readable time representation.
func (rev Revision) TimeString() string {
	return rev.Time.Format(time.RFC822)
}

// HyphaeLinks returns a comma-separated list of hyphae that were affected by this revision as HTML string.
func (rev Revision) HyphaeLinks() (html string) {
	hyphae := rev.hyphaeAffected()
	for i, hyphaName := range hyphae {
		if i > 0 {
			html += `<span aria-hidden="true">, </span>`
		}
		html += fmt.Sprintf(`<a href="/page/%[1]s">%[1]s</a>`, hyphaName)
	}
	return html
}

func (rev *Revision) descriptionForFeed() (html string) {
	return fmt.Sprintf(
		`<p>%s</p>
<p><b>Hyphae affected:</b> %s</p>`, rev.Message, rev.HyphaeLinks())
}

// Try and guess what link is the most important by looking at the message.
func (rev *Revision) bestLink() string {
	var (
		revs      = rev.hyphaeAffected()
		renameRes = renameMsgPattern.FindStringSubmatch(rev.Message)
	)
	switch {
	case renameRes != nil:
		return "/page/" + renameRes[1]
	case len(revs) == 0:
		return ""
	default:
		return "/page/" + revs[0]
	}
}

func (rev Revision) RecentChangesEntry() (html string) {
	if user.AuthUsed && rev.Username != "anon" {
		return fmt.Sprintf(`
<li class="rc-entry__time"><time>%[1]s</time></li>
<li class="rc-entry__hash">%[2]s</li>
<li class="rc-entry__links">%[5]s</li>
<li class="rc-entry__msg">%[6]s <span class="rc-entry__author">by <a href="/page/%[3]s/%[4]s" rel="author">%[4]s</a></span></li>
`, rev.TimeString(), rev.Hash, util.UserHypha, rev.Username, rev.HyphaeLinks(), rev.Message)
	}
	return fmt.Sprintf(`
<li class="rc-entry__time"><time>%[1]s</time></li>
<li class="rc-entry__hash">%[2]s</li>
<li class="rc-entry__links">%[3]s</li>
<li class="rc-entry__msg">%[4]s</li>
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
