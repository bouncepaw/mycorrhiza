package history

// information.go
// 	Things related to gathering existing information.
import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"

	"github.com/gorilla/feeds"
)

func recentChangesFeed() *feeds.Feed {
	feed := &feeds.Feed{
		Title:       "Recent changes",
		Link:        &feeds.Link{Href: cfg.URL},
		Description: "List of 30 recent changes on the wiki",
		Author:      &feeds.Author{Name: "Wikimind", Email: "wikimind@mycorrhiza"},
		Updated:     time.Now(),
	}
	var (
		out, err = silentGitsh(
			"log", "--oneline", "--no-merges",
			"--pretty=format:\"%h\t%ae\t%at\t%s\"",
			"--max-count=30",
		)
		revs []Revision
	)
	if err == nil {
		for _, line := range strings.Split(out.String(), "\n") {
			revs = append(revs, parseRevisionLine(line))
		}
	}
	log.Printf("Found %d recent changes", len(revs))
	for _, rev := range revs {
		feed.Add(&feeds.Item{
			Title:       rev.Message,
			Author:      &feeds.Author{Name: rev.Username},
			Id:          rev.Hash,
			Description: rev.descriptionForFeed(),
			Created:     rev.Time,
			Updated:     rev.Time,
			Link:        &feeds.Link{Href: cfg.URL + rev.bestLink()},
		})
	}
	return feed
}

// RecentChangesRSS creates recent changes feed in RSS format.
func RecentChangesRSS() (string, error) {
	return recentChangesFeed().ToRss()
}

// RecentChangesAtom creates recent changes feed in Atom format.
func RecentChangesAtom() (string, error) {
	return recentChangesFeed().ToAtom()
}

// RecentChangesJSON creates recent changes feed in JSON format.
func RecentChangesJSON() (string, error) {
	return recentChangesFeed().ToJSON()
}

// RecentChanges gathers an arbitrary number of latest changes in form of revisions slice.
func RecentChanges(n int) []Revision {
	var (
		out, err = silentGitsh(
			"log", "--oneline", "--no-merges",
			"--pretty=format:\"%h\t%ae\t%at\t%s\"",
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

// FileChanged tells you if the file has been changed.
func FileChanged(path string) bool {
	_, err := gitsh("diff", "--exit-code", path)
	return err != nil
}

// Revisions returns slice of revisions for the given hypha name.
func Revisions(hyphaName string) ([]Revision, error) {
	var (
		out, err = silentGitsh(
			"log", "--oneline", "--no-merges",
			// Hash, author email, author time, commit msg separated by tab
			"--pretty=format:\"%h\t%ae\t%at\t%s\"",
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

// WithRevisions returns an html representation of `revs` that is meant to be inserted in a history page.
func WithRevisions(hyphaName string, revs []Revision) (html string) {
	var (
		currentYear  int
		currentMonth time.Month
	)
	for i, rev := range revs {
		if rev.Time.Month() != currentMonth || rev.Time.Year() != currentYear {
			currentYear = rev.Time.Year()
			currentMonth = rev.Time.Month()
			if i != 0 {
				html += `
	</ul>
</section>`
			}
			html += fmt.Sprintf(`
<section class="history__month">
	<a href="#%[1]d-%[2]d" class="history__month-anchor">
		<h2 id="%[1]d-%[2]d" class="history__month-title">%[3]s</h2>
	</a>
	<ul class="history__entries">`,
				currentYear, currentMonth,
				strconv.Itoa(currentYear)+" "+rev.Time.Month().String())
		}
		html += rev.asHistoryEntry(hyphaName)
	}
	return html
}

func (rev *Revision) asHistoryEntry(hyphaName string) (html string) {
	author := ""
	if rev.Username != "anon" {
		author = fmt.Sprintf(`
		<span class="history-entry__author">by <a href="/hypha/%[1]s/%[2]s" rel="author">%[2]s</span>`, cfg.UserHypha, rev.Username)
	}
	return fmt.Sprintf(`
<li class="history__entry">
	<a class="history-entry" href="/rev/%[3]s/%[1]s">
		<time class="history-entry__time">%[2]s</time>
		<span class="history-entry__hash"><a href="/primitive-diff/%[3]s/%[1]s">%[3]s</a></span>
		<span class="history-entry__msg">%[4]s</span>
	</a>%[5]s
</li>
`, hyphaName, rev.timeToDisplay(), rev.Hash, rev.Message, author)
}

// Return time like mm-dd 13:42
func (rev *Revision) timeToDisplay() string {
	D := rev.Time.Day()
	h, m, _ := rev.Time.Clock()
	return fmt.Sprintf("%02d — %02d:%02d", D, h, m)
}

// This regex is wrapped in "". For some reason, these quotes appear at some time and we have to get rid of them.
var revisionLinePattern = regexp.MustCompile("\"(.*)\t(.*)@.*\t(.*)\t(.*)\"")

func parseRevisionLine(line string) Revision {
	results := revisionLinePattern.FindStringSubmatch(line)
	return Revision{
		Hash:     results[1],
		Username: results[2],
		Time:     *unixTimestampAsTime(results[3]),
		Message:  results[4],
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
