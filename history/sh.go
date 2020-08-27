package history

import (
	"fmt"
	"regexp"
	"strings"
)

func Revisions(filepath string) ([]Revision, error) {
	if filepath == "" {
		return []Revision{}, nil
	}
	var (
		out, err = gitsh(
			"log", "--oneline", "--no-merges",
			// Hash, Commiter email, Commiter time, Commit msg separated by tab
			"--pretty=format:\"%h\t%ce\t%ct\t%s\"",
			"--", filepath,
		)
		revs []Revision
	)
	if err == nil {
		for _, line := range strings.Split(out.String(), "\n") {
			if rev := parseRevisionLine(line); rev != nil {
				revs = append(revs, *rev)
			}
		}
	}
	return revs, err
}

// This regex is wrapped in "". For some reason, these quotes appear at some time and we have to get rid of them.
var revisionLinePattern = regexp.MustCompile("\"(.*)\t(.*)@.*\t(.*)\t(.*)\"")

func parseRevisionLine(line string) *Revision {
	var (
		results = revisionLinePattern.FindStringSubmatch(line)
		rev     = Revision{
			Hash:     results[1],
			Username: results[2],
			Time:     *unixTimestampAsTime(results[3]),
			Message:  results[4],
		}
	)
	return &rev
}

func (rev *Revision) AsHtmlTableRow(hyphaName string) string {
	return fmt.Sprintf(`
<tr>
	<td><a href="/rev/%s/%s">%s</a></td>
	<td>%s</td>
	<td><time>%s</time></td>
	<td>%s</td>
</tr>`, rev.Hash, hyphaName, rev.Hash, rev.Username, rev.Time.String(), rev.Message)
}

func FileAtRevision(filepath, hash string) (string, error) {
	out, err := gitsh("show", hash+":"+filepath)
	return out.String(), err
}
