package history

import (
	"fmt"
	"log"
	"time"

	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type OpType int

const (
	TypeNone OpType = iota
	TypeEditText
	TypeEditBinary
)

var WikiRepo *git.Repository
var Worktree *git.Worktree

func Start(wikiDir string) {
	ry, err := git.PlainOpen(wikiDir)
	if err != nil {
		log.Fatal(err)
	}
	WikiRepo = ry
	Worktree, err = WikiRepo.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki repository found")
}

type HistoryOp struct {
	Errs      []error
	opType    OpType
	userMsg   string
	signature *object.Signature
	isDone    bool
}

func Operation(opType OpType) *HistoryOp {
	hop := &HistoryOp{
		Errs:   []error{},
		opType: opType,
	}
	return hop
}

// WithFiles stages all passed `paths`. Paths can be rooted or not.
func (hop *HistoryOp) WithFiles(paths ...string) *HistoryOp {
	for _, path := range paths {
		if _, err := Worktree.Add(util.ShorterPath(path)); err != nil {
			log.Println(err)
			hop.Errs = append(hop.Errs, err)
		}
	}
	return hop
}

// WithMsg sets what message will be used for the future commit. If user message exceeds one line, it is stripped down.
func (hop *HistoryOp) WithMsg(userMsg string) *HistoryOp {
	// Isn't it too imperative?
	var firstLine string
	for _, ch := range userMsg {
		if ch == '\r' || ch == '\n' {
			break
		}
		firstLine += string(ch)
	}
	hop.userMsg = firstLine
	return hop
}

// WithSignature sets a signature for the future commit. You need to pass a username only, the rest is upon us (including email and time).
func (hop *HistoryOp) WithSignature(username string) *HistoryOp {
	hop.signature = &object.Signature{
		Name: username,
		// A fake email, why not
		Email: username + "@mycorrhiza",
		When:  time.Now(),
	}
	return hop
}

// Apply applies history operation by doing the commit. You can't apply the same operation more than once.
func (hop *HistoryOp) Apply() *HistoryOp {
	if !hop.isDone {
		opts := &git.CommitOptions{
			All:    false,
			Author: hop.signature,
		}
		err := opts.Validate(WikiRepo)
		if err != nil {
			hop.Errs = append(hop.Errs, err)
		}
		// TODO: work on this section:
		_, err = Worktree.Commit(hop.userMsg, opts)
		if err != nil {
			hop.Errs = append(hop.Errs, err)
		}
	}
	return hop
}

func Rename(from, to string) error {
	log.Println(util.ShorterPath(from), util.ShorterPath(to))
	_, err := Worktree.Move(util.ShorterPath(from), util.ShorterPath(to))
	return err
}

func StatusTable() (html string) {
	status, err := Worktree.Status()
	if err != nil {
		log.Fatal(err)
	}
	for path, stat := range status {
		html += fmt.Sprintf(`
	<tr>
		<td>%s</td>
		<td>%v</td>
	</tr>`, path, stat)
	}
	return `<table><tbody>` + html + `</tbody></table>`
}

func CommitsTable() (html string) {
	commitIter, err := WikiRepo.CommitObjects()
	if err != nil {
		log.Fatal(err)
	}
	err = commitIter.ForEach(func(commit *object.Commit) error {
		html += fmt.Sprintf(`
	<tr>
		<td>%v</td>
		<td>%v</td>
		<td>%s</td>
	</tr>`, commit.Hash, commit.Author, commit.Message)
		return nil
	})
	return `<table><tbody>` + html + `</tbody></table>`
}
