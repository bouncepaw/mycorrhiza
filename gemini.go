package main

// Gemini-related stuff. This is currently a proof-of-concept implementation, no one really uses it.
// Maybe we should deprecate it until we find power to do it properly?
//
// When this stuff gets more serious, a separate module will be needed.

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"

	"github.com/bouncepaw/mycomarkup/doc"
)

func geminiHomeHypha(w *gemini.ResponseWriter, rq *gemini.Request) {
	log.Println(rq.URL)
	_, _ = io.WriteString(w, `# MycorrhizaWiki

You have successfully served the wiki through Gemini. Currently, support is really work-in-progress; you should resort to using Mycorrhiza through the web protocols.

Visit home hypha:
=> /hypha/`+cfg.HomeHypha)
}

func geminiHypha(w *gemini.ResponseWriter, rq *gemini.Request) {
	log.Println(rq.URL)
	var (
		hyphaName = geminiHyphaNameFromRq(rq, "page", "hypha")
		h         = hyphae.ByName(hyphaName)
		hasAmnt   = h.Exists && h.BinaryPath != ""
		contents  string
	)
	if h.Exists {
		fileContentsT, errT := ioutil.ReadFile(h.TextPath)
		if errT == nil {
			md := doc.Doc(hyphaName, string(fileContentsT))
			contents = md.AsGemtext()
		}
	}
	if hasAmnt {
		_, _ = io.WriteString(w, "This hypha has an attachment\n")
	}
	_, _ = io.WriteString(w, contents)
}

func handleGemini() {
	if cfg.GeminiCertificatePath == "" {
		return
	}
	certPath, err := filepath.Abs(cfg.GeminiCertificatePath)
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

// geminiHyphaNameFromRq extracts hypha name from gemini request. You have to also pass the action which is embedded in the url or several actions. For url /hypha/hypha, the action would be "hypha".
func geminiHyphaNameFromRq(rq *gemini.Request, actions ...string) string {
	p := rq.URL.Path
	for _, action := range actions {
		if strings.HasPrefix(p, "/"+action+"/") {
			return util.CanonicalName(strings.TrimPrefix(p, "/"+action+"/"))
		}
	}
	log.Fatal("HyphaNameFromRq: no matching action passed")
	return ""
}
