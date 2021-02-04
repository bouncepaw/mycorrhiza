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
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/tree"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	http.HandleFunc("/page/", handlerHypha)
	http.HandleFunc("/hypha/", handlerHypha)
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
		TextPath          = hyphaName + ".myco"
		textContents, err = history.FileAtRevision(TextPath, revHash)
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
	w.Write([]byte(base(util.BeautifulName(hyphaName), page, u)))
}

// handlerText serves raw source text of the hypha.
func handlerText(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "text")
	if data, ok := HyphaStorage[hyphaName]; ok {
		log.Println("Serving", data.TextPath)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, rq, data.TextPath)
	}
}

// handlerBinary serves binary part of the hypha.
func handlerBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "binary")
	if data, ok := HyphaStorage[hyphaName]; ok {
		log.Println("Serving", data.BinaryPath)
		w.Header().Set("Content-Type", mimetype.FromExtension(filepath.Ext(data.BinaryPath)))
		http.ServeFile(w, rq, data.BinaryPath)
	}
}

// handlerHypha is the main hypha action that displays the hypha and the binary upload form along with some navigation.
func handlerHypha(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		hyphaName         = HyphaNameFromRq(rq, "page", "hypha")
		data, hyphaExists = HyphaStorage[hyphaName]
		hasAmnt           = hyphaExists && data.BinaryPath != ""
		contents          string
		openGraph         string
		u                 = user.FromRequest(rq)
	)
	if hyphaExists {
		fileContentsT, errT := ioutil.ReadFile(data.TextPath)
		_, errB := os.Stat(data.BinaryPath)
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
			util.BeautifulName(hyphaName),
			templates.PageHTML(rq, hyphaName,
				naviTitle(hyphaName),
				contents,
				treeHTML,
				prevHypha, nextHypha,
				hasAmnt),
			u,
			openGraph))
}
