package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/tree"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	http.HandleFunc("/page/", handlerPage)
	http.HandleFunc("/text/", handlerText)
	http.HandleFunc("/binary/", handlerBinary)
	http.HandleFunc("/rev/", handlerRevision)
}

// handlerRevision displays a specific revision of text part a page
func handlerRevision(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		shorterUrl        = strings.TrimPrefix(rq.URL.Path, "/rev/")
		firstSlashIndex   = strings.IndexRune(shorterUrl, '/')
		revHash           = shorterUrl[:firstSlashIndex]
		hyphaName         = CanonicalName(shorterUrl[firstSlashIndex+1:])
		contents          = fmt.Sprintf(`<p>This hypha had no text at this revision.</p>`)
		textPath          = hyphaName + ".myco"
		textContents, err = history.FileAtRevision(textPath, revHash)
		u                 = user.FromRequest(rq)
	)
	if err == nil {
		contents = markup.Doc(hyphaName, textContents).AsHTML()
	}
	treeHTML, _, _ := tree.Tree(hyphaName, IterateHyphaNamesWith)
	page := templates.RevisionHTML(
		rq,
		hyphaName,
		naviTitle(hyphaName),
		contents,
		treeHTML,
		revHash,
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base(hyphaName, page, u)))
}

// handlerText serves raw source text of the hypha.
func handlerText(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "text")
	if data, ok := HyphaStorage[hyphaName]; ok {
		log.Println("Serving", data.textPath)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, rq, data.textPath)
	}
}

// handlerBinary serves binary part of the hypha.
func handlerBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "binary")
	if data, ok := HyphaStorage[hyphaName]; ok {
		log.Println("Serving", data.binaryPath)
		w.Header().Set("Content-Type", ExtensionToMime(filepath.Ext(data.binaryPath)))
		http.ServeFile(w, rq, data.binaryPath)
	}
}

// handlerPage is the main hypha action that displays the hypha and the binary upload form along with some navigation.
func handlerPage(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName         = HyphaNameFromRq(rq, "page")
		data, hyphaExists = HyphaStorage[hyphaName]
		hasAmnt           = hyphaExists && data.binaryPath != ""
		contents          string
		openGraph         string
		u                 = user.FromRequest(rq)
	)
	if hyphaExists {
		fileContentsT, errT := ioutil.ReadFile(data.textPath)
		_, errB := os.Stat(data.binaryPath)
		if errT == nil {
			md := markup.Doc(hyphaName, string(fileContentsT))
			contents = md.AsHTML()
			openGraph = md.OpenGraphHTML()
		}
		if !os.IsNotExist(errB) {
			contents = binaryHtmlBlock(hyphaName, data) + contents
		}
	}
	treeHTML, prevHypha, nextHypha := tree.Tree(hyphaName, IterateHyphaNamesWith)
	util.HTTP200Page(w,
		templates.BaseHTML(
			hyphaName,
			templates.PageHTML(rq, hyphaName,
				naviTitle(hyphaName),
				contents,
				treeHTML, prevHypha, nextHypha,
				hasAmnt),
			u,
			openGraph))
}
