// Package history provides a git wrapper.
package history

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// Revision represents a revision, duh. Hash is usually short. Username is extracted from email.
type Revision struct {
	Hash              string
	Username          string
	Time              time.Time
	Message           string
	filesAffectedBuf  []string
	hyphaeAffectedBuf []string
}

// filesAffected tells what files have been affected by the revision.
func (rev *Revision) filesAffected() (filenames []string) {
	if nil != rev.filesAffectedBuf {
		return rev.filesAffectedBuf
	}
	// List of files affected by this revision, one per line.
	out, err := silentGitsh("diff-tree", "--no-commit-id", "--name-only", "-r", rev.Hash)
	// There's an error? Well, whatever, let's just assign an empty slice, who cares.
	if err != nil {
		rev.filesAffectedBuf = []string{}
	} else {
		rev.filesAffectedBuf = strings.Split(out.String(), "\n")
	}
	return rev.filesAffectedBuf
}

// determine what hyphae were affected by this revision
func (rev *Revision) hyphaeAffected() (hyphae []string) {
	if nil != rev.hyphaeAffectedBuf {
		return rev.hyphaeAffectedBuf
	}
	hyphae = make([]string, 0)
	var (
		// set is used to determine if a certain hypha has been already noted (hyphae are stored in 2 files at most currently).
		set       = make(map[string]bool)
		isNewName = func(hyphaName string) bool {
			if _, present := set[hyphaName]; present {
				return false
			}
			set[hyphaName] = true
			return true
		}
		filesAffected = rev.filesAffected()
	)
	for _, filename := range filesAffected {
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

// HyphaeLinksHTML returns a comma-separated list of hyphae that were affected by this revision as HTML string.
func (rev Revision) HyphaeLinksHTML() (html string) {
	hyphae := rev.hyphaeAffected()
	for i, hyphaName := range hyphae {
		if i > 0 {
			html += `<span aria-hidden="true">, </span>`
		}
		html += fmt.Sprintf(`<a href="/hypha/%[1]s">%[1]s</a>`, hyphaName)
	}
	return html
}

// descriptionForFeed generates a good enough HTML contents for a web feed.
func (rev *Revision) descriptionForFeed() (htmlDesc string) {
	return fmt.Sprintf(
		`<p>%s</p>
<p><b>Hyphae affected:</b> %s</p>
<pre><code>%s</code></pre>`, rev.Message, rev.HyphaeLinksHTML(), html.EscapeString(rev.textDiff()))
}

// textDiff generates a good enough diff to display in a web feed. It is not html-escaped.
func (rev *Revision) textDiff() (diff string) {
	filenames, ok := rev.mycoFiles()
	if !ok {
		return "No text changes"
	}
	for _, filename := range filenames {
		text, err := PrimitiveDiffAtRevision(filename, rev.Hash)
		if err != nil {
			diff += "\nAn error has occured with " + filename + "\n"
		}
		diff += text + "\n"
	}
	return diff
}

// mycoFiles returns filenames of .myco file. It is not ok if there are no myco files.
func (rev *Revision) mycoFiles() (filenames []string, ok bool) {
	filenames = []string{}
	for _, filename := range rev.filesAffected() {
		if strings.HasSuffix(filename, ".myco") {
			filenames = append(filenames, filename)
		}
	}
	return filenames, len(filenames) > 0
}

// Try and guess what link is the most important by looking at the message.
func (rev *Revision) bestLink() string {
	var (
		revs      = rev.hyphaeAffected()
		renameRes = renameMsgPattern.FindStringSubmatch(rev.Message)
	)
	switch {
	case renameRes != nil:
		return "/hypha/" + renameRes[1]
	case len(revs) == 0:
		return ""
	default:
		return "/hypha/" + revs[0]
	}
}

// I pronounce it as [gɪt͡ʃ].
// gitsh is async-safe, therefore all other git-related functions in this module are too.
func gitsh(args ...string) (out bytes.Buffer, err error) {
	fmt.Printf("$ %v\n", args)
	cmd := exec.Command(gitpath, args...)
	cmd.Dir = files.HyphaeDir()
	cmd.Env = gitEnv

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
	cmd.Env = gitEnv

	b, err := cmd.CombinedOutput()
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
