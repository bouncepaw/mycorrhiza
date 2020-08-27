package history

import (
	"fmt"
	"log"

	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/go-git/go-git/v5"
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
	gitsh("config", "user.name", "wikimind")
	gitsh("config", "user.email", "wikimind@mycorrhiza")
	log.Println("Wiki repository found")
}

type HistoryOp struct {
	Errs      []error
	opType    OpType
	userMsg   string
	signature *object.Signature
	name      string
	email     string
}

func Operation(opType OpType) *HistoryOp {
	hop := &HistoryOp{
		Errs:   []error{},
		opType: opType,
	}
	return hop
}

func (hop *HistoryOp) gitop(args ...string) *HistoryOp {
	out, err := gitsh(args...)
	fmt.Println("out:", out.String())
	if err != nil {
		hop.Errs = append(hop.Errs, err)
	}
	return hop
}

// WithFiles stages all passed `paths`. Paths can be rooted or not.
func (hop *HistoryOp) WithFiles(paths ...string) *HistoryOp {
	for i, path := range paths {
		paths[i] = util.ShorterPath(path)
	}
	return hop.gitop(append([]string{"add"}, paths...)...)
}

// WithMsg sets what message will be used for the future commit. If user message exceeds one line, it is stripped down.
func (hop *HistoryOp) WithMsg(userMsg string) *HistoryOp {
	for _, ch := range userMsg {
		if ch == '\r' || ch == '\n' {
			break
		}
		hop.userMsg += string(ch)
	}
	return hop
}

// WithSignature sets a signature for the future commit. You need to pass a username only, the rest is upon us (including email and time).
func (hop *HistoryOp) WithSignature(username string) *HistoryOp {
	hop.name = username
	hop.email = username + "@mycorrhiza" // A fake email, why not
	return hop
}

// Apply applies history operation by doing the commit.
func (hop *HistoryOp) Apply() *HistoryOp {
	hop.gitop(
		"commit",
		"--author='"+hop.name+" <"+hop.email+">'",
		"--message="+hop.userMsg,
	)
	return hop
}

// Rename renames from `from` to `to`. NB. It uses os.Rename internally rather than git.Move because git.Move works wrong for some reason.
func Rename(from, to string) error {
	log.Println(util.ShorterPath(from), util.ShorterPath(to))
	_, err := gitsh("mv", from, to)
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
