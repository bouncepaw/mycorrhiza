package main

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"

	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/util"
)

func geminiHomeHypha(w *gemini.ResponseWriter, rq *gemini.Request) {
	log.Println(rq.URL)
	w.Write([]byte(`# MycorrhizaWiki

You have successfully served the wiki through Gemini. Currently, support is really work-in-progress; you should resort to using Mycorrhiza through the web protocols.

Visit home hypha:
=> /hypha/` + util.HomePage))
}

func geminiHypha(w *gemini.ResponseWriter, rq *gemini.Request) {
	log.Println(rq.URL)
	var (
		hyphaName         = geminiHyphaNameFromRq(rq, "page", "hypha")
		data, hyphaExists = HyphaStorage[hyphaName]
		hasAmnt           = hyphaExists && data.BinaryPath != ""
		contents          string
	)
	if hyphaExists {
		fileContentsT, errT := ioutil.ReadFile(data.TextPath)
		if errT == nil {
			md := markup.Doc(hyphaName, string(fileContentsT))
			contents = md.AsGemtext()
		}
	}
	if hasAmnt {
		w.Write([]byte("This hypha has an attachment\n"))
	}
	w.Write([]byte(contents))
}

func handleGemini() {
	if util.GeminiCertPath == "" {
		return
	}
	certPath, err := filepath.Abs(util.GeminiCertPath)
	if err != nil {
		log.Fatal(err)
	}

	var server gemini.Server
	server.ReadTimeout = 30 * time.Second
	server.WriteTimeout = 1 * time.Minute
	if err := server.Certificates.Load(certPath); err != nil {
		log.Fatal(err)
	}
	server.CreateCertificate = func(hostname string) (tls.Certificate, error) {
		return certificate.Create(certificate.CreateOptions{
			Subject: pkix.Name{
				CommonName: hostname,
			},
			DNSNames: []string{hostname},
			Duration: 365 * 24 * time.Hour,
		})
	}

	var mux gemini.ServeMux
	mux.HandleFunc("/", geminiHomeHypha)
	mux.HandleFunc("/hypha/", geminiHypha)
	mux.HandleFunc("/page/", geminiHypha)

	server.Handle("localhost", &mux)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
