package history

import (
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

/*
// 10 should be enough
const ShortHashLength = 10

type EditType int

const (
	TypeRename EditType = iota
	TypeDelete
	TypeEditText
	TypeEditBinary
)

type Revision struct {
	ShortHash [ShortHashLength]byte
	Type      EditType
}*/

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

func Stage(path string) error {
	var err error
	_, err = Worktree.Add(path)
	if err != nil {
		log.Println(err)
	}
	return err
}

func CommitTest() {
	opts := &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			Name:  "wikimind",
			Email: "wikimind@thiswiki",
			When:  time.Now(),
		},
	}
	err := opts.Validate(WikiRepo)
	if err != nil {
		log.Fatal(err)
	}
	_, err = Worktree.Commit("This is a test commit", opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Made a test commit")
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
