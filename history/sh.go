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
)

type Revision struct {
	Hash     string
	Username string
	Time     time.Time
	Message  string
}

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
func gitsh(args ...string) (out bytes.Buffer, err error) {
	cmd := exec.Command(gitpath, args...)
	cmd.Stdout = &out
	cmd.Run()
	if err != nil {
		log.Println("gitsh:", err)
	}
	return out, err
}

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

func unixTimestampAsTime(ts string) *time.Time {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil
	}
	tm := time.Unix(i, 0)
	return &tm
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

func (rev *Revision) AsHtmlTableRow() string {
	return fmt.Sprintf(`
<tr>
	<td>%s</td>
	<td>%s</td>
	<td><time>%s</time></td>
	<td>%s</td>
</tr>`, rev.Hash, rev.Username, rev.Time.String(), rev.Message)
}
