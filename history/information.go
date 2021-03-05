// information.go
// 	Things related to gathering existing information.
package history

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/gorilla/feeds"
)

func recentChangesFeed() *feeds.Feed {
	feed := &feeds.Feed{
		Title:       "Recent changes",
		Link:        &feeds.Link{Href: util.URL},
		Description: "List of 30 recent changes on the wiki",
		Author:      &feeds.Author{Name: "Wikimind", Email: "wikimind@mycorrhiza"},
		Updated:     time.Now(),
	}
	var (
		out, err = gitsh(
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
	for _, rev := range revs {
		feed.Add(&feeds.Item{
			Title:       rev.Message,
			Author:      &feeds.Author{Name: rev.Username},
			Id:          rev.Hash,
			Description: rev.descriptionForFeed(),
			Created:     rev.Time,
			Updated:     rev.Time,
			Link:        &feeds.Link{Href: util.URL + rev.bestLink()},
		})
	}
	return feed
}

func RecentChangesRSS() (string, error) {
	return recentChangesFeed().ToRss()
}

func RecentChangesAtom() (string, error) {
	return recentChangesFeed().ToAtom()
}

func RecentChangesJSON() (string, error) {
	return recentChangesFeed().ToJSON()
}

func RecentChanges(n int) []Revision {
	var (
		out, err = gitsh(
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
		out, err = gitsh(
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
	return revs, err
}

// HistoryWithRevisions returns an html representation of `revs` that is meant to be inserted in a history page.
func HistoryWithRevisions(hyphaName string, revs []Revision) (html string) {
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
		<span class="history-entry__author">by <a href="/page/%[1]s/%[2]s" rel="author">%[2]s</span>`, util.UserHypha, rev.Username)
	}
	return fmt.Sprintf(`
<li class="history__entry">
	<a class="history-entry" href="/rev/%[3]s/%[1]s">
		<time class="history-entry__time">%[2]s</time>
		<span class="history-entry__hash">%[3]s</span>
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

// See how the file with `filepath` looked at commit with `hash`.
func FileAtRevision(filepath, hash string) (string, error) {
	out, err := gitsh("show", hash+":"+strings.TrimPrefix(filepath, util.WikiDir+"/"))
	return out.String(), err
}
