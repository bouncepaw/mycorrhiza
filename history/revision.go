package history

import (
	"fmt"
	"html"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"
)

// WithRevisions returns an HTML representation of `revs` that is meant to be inserted in a history page.
func WithRevisions(hyphaName string, revs []Revision) string {
	var buf strings.Builder

	for _, grp := range groupRevisionsByMonth(revs) {
		currentYear := grp[0].Time.Year()
		currentMonth := grp[0].Time.Month()
		sectionId := fmt.Sprintf("%04d-%02d", currentYear, currentMonth)

		buf.WriteString(fmt.Sprintf(
			`<section class="history__month">
	<a href="#%s" class="history__month-anchor">
		<h2 id="%s" class="history__month-title">%d %s</h2>
	</a>
	<ul class="history__entries">`,
			sectionId, sectionId, currentYear, currentMonth.String(),
		))

		for _, rev := range grp {
			buf.WriteString(fmt.Sprintf(
				`<li class="history__entry">
	<a class="history-entry" href="/rev/%s/%s">
		<time class="history-entry__time">%s</time>
	</a>
	<span class="history-entry__hash"><a href="/primitive-diff/%s/%s">%s</a></span>
	<span class="history-entry__msg">%s</span>`,
				rev.Hash, hyphaName,
				rev.timeToDisplay(),
				rev.Hash, hyphaName, rev.Hash,
				html.EscapeString(rev.Message),
			))

			if rev.Username != "anon" {
				buf.WriteString(fmt.Sprintf(
					`<span class="history-entry__author">by <a href="/hypha/%s/%s" rel="author">%s</a></span>`,
					cfg.UserHypha, rev.Username, rev.Username,
				))
			}

			buf.WriteString("</li>\n")
		}

		buf.WriteString(`</ul></section>`)
	}

	return buf.String()
}

// Revision represents a revision of a hypha.
type Revision struct {
	// Hash is usually short.
	Hash string
	// Username is extracted from email.
	Username          string
	Time              time.Time
	Message           string
	filesAffectedBuf  []string
	hyphaeAffectedBuf []string
}

// HyphaeDiffsHTML returns a comma-separated list of diffs links of current revision for every affected file as HTML string.
func (rev Revision) HyphaeDiffsHTML() string {
	entries := rev.hyphaeAffected()
	if len(entries) == 1 {
		return fmt.Sprintf(
			`<a href="/primitive-diff/%s/%s">%s</a>`,
			rev.Hash, entries[0], rev.Hash,
		)
	}

	var buf strings.Builder
	for i, hyphaName := range entries {
		if i > 0 {
			buf.WriteString(`<span aria-hidden="true">, </span>`)
		}
		buf.WriteString(`<a href="/primitive-diff/`)
		buf.WriteString(rev.Hash)
		buf.WriteString(`/`)
		buf.WriteString(hyphaName)
		buf.WriteString(`">`)
		if i == 0 {
			buf.WriteString(rev.Hash)
			buf.WriteString("&nbsp;")
		}
		buf.WriteString(hyphaName)
		buf.WriteString(`</a>`)
	}
	return buf.String()
}

// descriptionForFeed generates a good enough HTML contents for a web feed.
func (rev *Revision) descriptionForFeed() string {
	return fmt.Sprintf(
		`<p><b>%s</b> (by %s at %s)</p>
<p>Hyphae affected: %s</p>
<pre><code>%s</code></pre>`,
		rev.Message, rev.Username, rev.TimeString(),
		rev.HyphaeLinksHTML(),
		rev.textDiff(),
	)
}

// HyphaeLinksHTML returns a comma-separated list of hyphae that were affected by this revision as HTML string.
func (rev Revision) HyphaeLinksHTML() string {
	var buf strings.Builder
	for i, hyphaName := range rev.hyphaeAffected() {
		if i > 0 {
			buf.WriteString(`<span aria-hidden="true">, <span>`)
		}

		urlSafeHyphaName := url.PathEscape(hyphaName)
		buf.WriteString(fmt.Sprintf(`<a href="/hypha/%s">%s</a>`, urlSafeHyphaName, hyphaName))
	}
	return buf.String()
}

// gitLog calls `git log` and parses the results.
func gitLog(args ...string) ([]Revision, error) {
	args = append([]string{
		"log", "--abbrev-commit", "--no-merges",
		"--pretty=format:%h\t%ae\t%at\t%s",
	}, args...)
	args = append(args, "--")
	out, err := silentGitsh(args...)
	if strings.Contains(out.String(), "bad revision 'HEAD'") {
		// Then we have no recent changes! It's a hack.
		return nil, nil
	}
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
		// TODO: return error
		slog.Error("Failed to git log", "err", err)
		os.Exit(1)
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
	slog.Info("Found recent changes", "n", len(revs))
	return revs
}

// Revisions returns slice of revisions for the given hypha name, ordered most recent first.
func Revisions(hyphaName string) ([]Revision, error) {
	revs, err := gitLog("--", hyphaName+".*")
	slog.Info("Found revisions", "hyphaName", hyphaName, "n", len(revs), "err", err)
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
