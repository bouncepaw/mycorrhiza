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

// gitLog calls `git log` and parses the results.
func gitLog(args ...string) ([]Revision, error) {
	args = append([]string{
		"log", "--abbrev-commit", "--no-merges",
		"--pretty=format:%h\t%ae\t%at\t%s",
	}, args...)
	out, err := silentGitsh(args...)
	if err != nil {
		return nil, err
	}

	outStr := out.String()
	if outStr == "" {
		// if there are no commits to return
		return nil, nil
	}

	var revs []Revision
	for _, line := range strings.Split(outStr, "\n") {
		revs = append(revs, parseRevisionLine(line))
	}
	return revs, nil
}

type recentChangesStream struct {
	currHash string
}

func newRecentChangesStream() recentChangesStream {
	// next returns the next n revisions from the stream, ordered most recent first.
	// If there are less than n revisions remaining, it will return only those.
	return recentChangesStream{currHash: ""}
}

func (stream *recentChangesStream) next(n int) []Revision {
	args := []string{"--max-count=" + strconv.Itoa(n)}
	if stream.currHash == "" {
		args = append(args, "HEAD")
	} else {
		// currHash is the last revision from the last call, so skip it
		args = append(args, "--skip=1", stream.currHash)
	}
	res, err := gitLog(args...)
	if err != nil {
		log.Fatal(err)
	}
	if len(res) != 0 {
		stream.currHash = res[len(res)-1].Hash
	}
	return res
}

// recentChangesIterator returns a function that returns successive revisions from the stream.
// It buffers revisions to avoid calling git every time.
func (stream recentChangesStream) iterator() func() (Revision, bool) {
	var buf []Revision
	return func() (Revision, bool) {
		if len(buf) == 0 {
			// no real reason to choose 30, just needs some large number
			buf = stream.next(30)
			if len(buf) == 0 {
				// revs has no revisions left
				return Revision{}, true
			}
		}
		rev := buf[0]
		buf = buf[1:]
		return rev, false
	}
}

// RecentChanges gathers an arbitrary number of latest changes in form of revisions slice, ordered most recent first.
func RecentChanges(n int) []Revision {
	stream := newRecentChangesStream()
	revs := stream.next(n)
	log.Printf("Found %d recent changes", len(revs))
	return revs
}

// Revisions returns slice of revisions for the given hypha name, ordered most recent first.
func Revisions(hyphaName string) ([]Revision, error) {
	revs, err := gitLog("--", hyphaName+".*")
	log.Printf("Found %d revisions for ‘%s’\n", len(revs), hyphaName)
	return revs, err
}

// FileChanged tells you if the file has been changed since the last commit.
func FileChanged(path string) bool {
	_, err := gitsh("diff", "--exit-code", path)
	return err != nil
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
		if strings.ContainsRune(filename, '.') {
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

// SplitPrimitiveDiff splits a primitive diff of a single file into hunks.
func SplitPrimitiveDiff(text string) (result []string) {
	idx := strings.Index(text, "@@ -")
	if idx < 0 {
		return
	}
	text = text[idx:]
	for {
		idx = strings.Index(text, "\n@@ -")
		if idx < 0 {
			result = append(result, text)
			return
		}
		result = append(result, text[:idx+1])
		text = text[idx+1:]
	}
}
