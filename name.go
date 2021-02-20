package main

import (
	"log"
	"net/http"
	"strings"

	"git.sr.ht/~adnano/go-gemini"

	"github.com/bouncepaw/mycorrhiza/util"
)

// HyphaNameFromRq extracts hypha name from http request. You have to also pass the action which is embedded in the url or several actions. For url /hypha/hypha, the action would be "hypha".
func HyphaNameFromRq(rq *http.Request, actions ...string) string {
	p := rq.URL.Path
	for _, action := range actions {
		if strings.HasPrefix(p, "/"+action+"/") {
			return util.CanonicalName(strings.TrimPrefix(p, "/"+action+"/"))
		}
	}
	panic("HyphaNameFromRq: no matching action passed")
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
