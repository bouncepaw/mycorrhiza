package history

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/files"
)

// Revision represents a revision, duh. Hash is usually short. Username is extracted from email.
type Revision struct {
	Hash              string
	Username          string
	Time              time.Time
	Message           string
	filesAffectedBuf  []string
	hyphaeAffectedBuf []string
}

// RecentChanges gathers an arbitrary number of latest changes in form of revisions slice, ordered most recent first.
func RecentChanges(n int) []Revision {
	var (
		out, err = silentGitsh(
			"log", "--oneline", "--no-merges",
			"--pretty=format:%h\t%ae\t%at\t%s",
			"--max-count="+strconv.Itoa(n),
		)
		revs []Revision
	)
	if err == nil {
		for _, line := range strings.Split(out.String(), "\n") {
			revs = append(revs, parseRevisionLine(line))
		}
	}
	log.Printf("Found %d recent changes", len(revs))
	return revs
}

// FileChanged tells you if the file has been changed since the last commit.
func FileChanged(path string) bool {
	_, err := gitsh("diff", "--exit-code", path)
	return err != nil
}

// Revisions returns slice of revisions for the given hypha name, ordered most recent first.
func Revisions(hyphaName string) ([]Revision, error) {
	var (
		out, err = silentGitsh(
			"log", "--oneline", "--no-merges",
			// Hash, author email, author time, commit msg separated by tab
			"--pretty=format:%h\t%ae\t%at\t%s",
			"--", hyphaName+".*",
		)
		revs []Revision
	)
	if err == nil {
		for _, line := range strings.Split(out.String(), "\n") {
			if line != "" {
				revs = append(revs, parseRevisionLine(line))
			}
		}
	}
	log.Printf("Found %d revisions for ‘%s’\n", len(revs), hyphaName)
	return revs, err
}

// Return time like dd — 13:42
func (rev *Revision) timeToDisplay() string {
	D := rev.Time.Day()
	h, m, _ := rev.Time.Clock()
	return fmt.Sprintf("%02d — %02d:%02d", D, h, m)
}

var revisionLinePattern = regexp.MustCompile("(.*)\t(.*)@.*\t(.*)\t(.*)")

// Convert a UNIX timestamp as string into a time. If nil is returned, it means that the timestamp could not be converted.
func unixTimestampAsTime(ts string) *time.Time {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil
	}
	tm := time.Unix(i, 0)
	return &tm
}

func parseRevisionLine(line string) Revision {
	results := revisionLinePattern.FindStringSubmatch(line)
	return Revision{
		Hash:     results[1],
		Username: results[2],
		Time:     *unixTimestampAsTime(results[3]),
		Message:  results[4],
	}
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

// textDiff generates a good enough diff to display in a web feed. It is not html-escaped.
func (rev *Revision) textDiff() (diff string) {
	filenames, ok := rev.mycoFiles()
	if !ok {
		return "No text changes"
	}
	for _, filename := range filenames {
		text, err := PrimitiveDiffAtRevision(filename, rev.Hash)
		if err != nil {
			diff += "\nAn error has occurred with " + filename + "\n"
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

// FileAtRevision shows how the file with the given file path looked at the commit with the hash. It may return an error if git fails.
func FileAtRevision(filepath, hash string) (string, error) {
	out, err := gitsh("show", hash+":"+strings.TrimPrefix(filepath, files.HyphaeDir()+"/"))
	if err != nil {
		return "", err
	}
	return out.String(), err
}

// PrimitiveDiffAtRevision generates a plain-text diff for the given filepath at the commit with the given hash. It may return an error if git fails.
func PrimitiveDiffAtRevision(filepath, hash string) (string, error) {
	out, err := silentGitsh("diff", "--unified=1", "--no-color", hash+"~", hash, "--", filepath)
	if err != nil {
		return "", err
	}
	return out.String(), err
}
