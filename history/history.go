package history

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/bouncepaw/mycorrhiza/util"
)

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
	Hash     string
	Username string
	Time     time.Time
	Message  string
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
	_, err := gitsh("mv", from, to)
	return err
}
