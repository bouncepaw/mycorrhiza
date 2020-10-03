// information.go
// 	Things related to gathering existing information.
package history

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bouncepaw/mycorrhiza/templates"
)

func RecentChanges(n int) string {
	var (
		out, err = gitsh(
			"log", "--oneline", "--no-merges",
			"--pretty=format:\"%h\t%ce\t%ct\t%s\"",
			"--max-count="+strconv.Itoa(n),
		)
		revs []Revision
	)
	if err == nil {
		for _, line := range strings.Split(out.String(), "\n") {
			revs = append(revs, parseRevisionLine(line))
		}
	}
	entries := make([]string, len(revs))
	for i, rev := range revs {
		entries[i] = rev.RecentChangesEntry()
	}
	return templates.RecentChangesHTML(entries, n)
}

// Revisions returns slice of revisions for the given hypha name.
func Revisions(hyphaName string) ([]Revision, error) {
	var (
		out, err = gitsh(
			"log", "--oneline", "--no-merges",
			// Hash, Commiter email, Commiter time, Commit msg separated by tab
			"--pretty=format:\"%h\t%ce\t%ct\t%s\"",
			"--", hyphaName+"&.*",
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

// Represent revision as a table row.
func (rev *Revision) AsHtmlTableRow(hyphaName string) string {
	return fmt.Sprintf(`
<tr>
	<td><time>%s</time></td>
	<td><a href="/rev/%s/%s">%s</a></td>
	<td>%s</td>
</tr>`, rev.TimeString(), rev.Hash, hyphaName, rev.Hash, rev.Message)
}

// See how the file with `filepath` looked at commit with `hash`.
func FileAtRevision(filepath, hash string) (string, error) {
	out, err := gitsh("show", hash+":"+filepath)
	return out.String(), err
}
