package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/util"
)

// WikiDir is a rooted path to the wiki storage directory.
var WikiDir string

// HyphaPattern is a pattern which all hyphae must match. Not used currently.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"\'&%]+`)

// HyphaStorage is a mapping between canonical hypha names and their meta information.
var HyphaStorage = make(map[string]*HyphaData)

// IterateHyphaNamesWith is a closure to be passed to subpackages to let them iterate all hypha names read-only.
func IterateHyphaNamesWith(f func(string)) {
	for hyphaName, _ := range HyphaStorage {
		f(hyphaName)
	}
}

// HttpErr is used by many handlers to signal errors in a compact way.
func HttpErr(w http.ResponseWriter, status int, name, title, errMsg string) {
	log.Println(errMsg, "for", name)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, base(title, fmt.Sprintf(
		`<p>%s. <a href="/page/%s">Go back to the hypha.<a></p>`,
		errMsg, name)))
}

// shorterPath is used by handlerList to display shorter path to the files. It simply strips WikiDir. It was moved to util package, this is an alias. TODO: demolish.
var shorterPath = util.ShorterPath

// Show all hyphae
func handlerList(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf := `
	<h1>List of pages</h1>
	<table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Text path</th>
				<th>Text type</th>
				<th>Binary path</th>
				<th>Binary type</th>
			</tr>
		</thead>
		<tbody>`
	for name, data := range HyphaStorage {
		buf += fmt.Sprintf(`
			<tr>
				<td><a href="/page/%s">%s</a></td>
				<td>%s</td>
				<td>%d</td>
				<td>%s</td>
				<td>%d</td>
			</tr>`,
			name, name,
			shorterPath(data.textPath), data.textType,
			shorterPath(data.binaryPath), data.binaryType,
		)
	}
	buf += `
		</tbody>
	</table>
`
	w.Write([]byte(base("List of pages", buf)))
}

// This part is present in all html documents.
func base(title, body string) string {
	return fmt.Sprintf(`
<!doctype html>
<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" type="text/css" href="/static/common.css">
		<title>%s</title>
	</head>
	<body>
		%s
	</body>
</html>
`, title, body)
}

// Reindex all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	HyphaStorage = make(map[string]*HyphaData)
	log.Println("Wiki storage directory is", WikiDir)
	log.Println("Start indexing hyphae...")
	Index(WikiDir)
	log.Println("Indexed", len(HyphaStorage), "hyphae")
}

func handlerListCommits(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(history.CommitsTable()))
}

func handlerStatus(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(history.StatusTable()))
}

func handlerCommitTest(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	history.CommitTest()
	w.Write([]byte("if you are here, a commit has been done"))
}

func main() {
	log.Println("Running MycorrhizaWiki β")

	var err error
	WikiDir, err = filepath.Abs(os.Args[1])
	util.WikiDir = WikiDir
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Wiki storage directory is", WikiDir)
	log.Println("Start indexing hyphae...")
	Index(WikiDir)
	log.Println("Indexed", len(HyphaStorage), "hyphae")

	history.Start(WikiDir)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(WikiDir+"/static"))))
	// See http_readers.go for /page/, /text/, /binary/.
	// See http_mutators.go for /upload-binary/, /upload-text/, /edit/.
	http.HandleFunc("/list", handlerList)
	http.HandleFunc("/reindex", handlerReindex)
	http.HandleFunc("/git/list", handlerListCommits)
	http.HandleFunc("/git/status", handlerStatus)
	http.HandleFunc("/git/commit", handlerCommitTest)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, WikiDir+"/static/favicon.ico")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/page/home", http.StatusSeeOther)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:1737", nil))
}
