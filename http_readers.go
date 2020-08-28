package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bouncepaw/mycorrhiza/gemtext"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/tree"
)

func init() {
	http.HandleFunc("/page/", handlerPage)
	http.HandleFunc("/text/", handlerText)
	http.HandleFunc("/binary/", handlerBinary)
	http.HandleFunc("/history/", handlerHistory)
	http.HandleFunc("/rev/", handlerRevision)
}

// handlerRevision displays a specific revision of text part a page
func handlerRevision(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		shorterUrl        = strings.TrimPrefix(rq.URL.Path, "/rev/")
		revHash           = path.Dir(shorterUrl)
		hyphaName         = CanonicalName(strings.TrimPrefix(shorterUrl, revHash+"/"))
		contents          = fmt.Sprintf(`<p>This hypha had no text at this revision.</p>`)
		textPath          = hyphaName + "&.gmi"
		textContents, err = history.FileAtRevision(textPath, revHash)
	)
	log.Println(revHash, hyphaName, textPath, textContents, err)
	if err == nil {
		contents = gemtext.ToHtml(hyphaName, textContents)
	}
	form := fmt.Sprintf(`
		<main>
			<nav>
				<ul>
					<li><a href="/page/%[1]s">See the latest revision</a></li>
					<li><a href="/history/%[1]s">History</a></li>
				</ul>
			</nav>
			<article>
				<p>Please note that viewing binary parts of hyphae is not supported in history for now.</p>
				%[2]s
				%[3]s
			</article>
			<hr/>
			<aside>
				%[4]s
			</aside>
		</main>
`, hyphaName,
		naviTitle(hyphaName),
		contents,
		tree.TreeAsHtml(hyphaName, IterateHyphaNamesWith))
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base(hyphaName, form)))
}

// handlerHistory lists all revisions of a hypha
func handlerHistory(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "history")
	var tbody string
	if _, ok := HyphaStorage[hyphaName]; ok {
		revs, err := history.Revisions(hyphaName)
		if err == nil {
			for _, rev := range revs {
				tbody += rev.AsHtmlTableRow(hyphaName)
			}
		}
		log.Println(revs)
	}

	table := fmt.Sprintf(`
		<main>
			<table>
				<thead>
					<tr>
						<th>Hash</th>
						<th>Username</th>
						<th>Time</th>
						<th>Message</th>
					</tr>
				</thead>
				<tbody>
					%s
				</tbody>
			</table>
		</main>`, tbody)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base(hyphaName, table)))
}

// handlerText serves raw source text of the hypha.
func handlerText(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "text")
	if data, ok := HyphaStorage[hyphaName]; ok {
		log.Println("Serving", data.textPath)
		w.Header().Set("Content-Type", data.textType.Mime())
		http.ServeFile(w, rq, data.textPath)
	}
}

// handlerBinary serves binary part of the hypha.
func handlerBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "binary")
	if data, ok := HyphaStorage[hyphaName]; ok {
		log.Println("Serving", data.binaryPath)
		w.Header().Set("Content-Type", data.binaryType.Mime())
		http.ServeFile(w, rq, data.binaryPath)
	}
}

// handlerPage is the main hypha action that displays the hypha and the binary upload form along with some navigation.
func handlerPage(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName         = HyphaNameFromRq(rq, "page")
		contents          = fmt.Sprintf(`<p>This hypha has no text. Why not <a href="/edit/%s">create it</a>?</p>`, hyphaName)
		data, hyphaExists = HyphaStorage[hyphaName]
	)
	if hyphaExists {
		fileContentsT, errT := ioutil.ReadFile(data.textPath)
		_, errB := os.Stat(data.binaryPath)
		if errT == nil {
			contents = gemtext.ToHtml(hyphaName, string(fileContentsT))
		}
		if !os.IsNotExist(errB) {
			contents = binaryHtmlBlock(hyphaName, data) + contents
		}
	}
	form := fmt.Sprintf(`
		<main>
			<nav>
				<ul>
					<li><a href="/edit/%[1]s">Edit</a></li>
					<li><a href="/text/%[1]s">Raw text</a></li>
					<li><a href="/binary/%[1]s">Binary part</a></li>
					<li><a href="/history/%[1]s">History</a></li>
				</ul>
			</nav>
			<article>
				%[2]s
				%[3]s
			</article>
			<hr/>
			<form action="/upload-binary/%[1]s"
			      method="post" enctype="multipart/form-data">
				<label for="upload-binary__input">Upload new binary part</label>
				<br>
				<input type="file" id="upload-binary__input" name="binary"/>
				<input type="submit"/>
			</form>
			<hr/>
			<aside>
				%[4]s
			</aside>
		</main>
`, hyphaName,
		naviTitle(hyphaName),
		contents,
		tree.TreeAsHtml(hyphaName, IterateHyphaNamesWith))
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base(hyphaName, form)))
}
