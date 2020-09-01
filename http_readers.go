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
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/tree"
	"github.com/bouncepaw/mycorrhiza/util"
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
	if err == nil {
		contents = gemtext.ToHtml(hyphaName, textContents)
	}
	page := templates.RevisionHTML(
		hyphaName,
		naviTitle(hyphaName),
		contents,
		tree.TreeAsHtml(hyphaName, IterateHyphaNamesWith),
		revHash,
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base(hyphaName, page)))
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

	util.HTTP200Page(w,
		base(hyphaName, templates.HistoryHTML(hyphaName, tbody)))
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
		data, hyphaExists = HyphaStorage[hyphaName]
		contents          string
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
	util.HTTP200Page(w, base(hyphaName, templates.PageHTML(hyphaName,
		naviTitle(hyphaName),
		contents,
		tree.TreeAsHtml(hyphaName, IterateHyphaNamesWith))))
}
