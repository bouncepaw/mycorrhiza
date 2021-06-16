package util

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycomarkup/util"
	"github.com/bouncepaw/mycorrhiza/cfg"
)

// PrepareRq strips the trailing / in rq.URL.Path. In the future it might do more stuff for making all request structs uniform.
func PrepareRq(rq *http.Request) {
	rq.URL.Path = strings.TrimSuffix(rq.URL.Path, "/")
}

// ShorterPath is used by handlerList to display shorter path to the files. It simply strips WikiDir.
func ShorterPath(path string) string {
	if strings.HasPrefix(path, cfg.WikiDir) {
		tmp := strings.TrimPrefix(path, cfg.WikiDir)
		if tmp == "" {
			return ""
		}
		return tmp[1:]
	}
	return path
}

// HTTP404Page writes a 404 error in the status, needed when no content is found on the page.
func HTTP404Page(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(page))
}

// HTTP200Page wraps some frequently used things for successful 200 responses.
func HTTP200Page(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page))
}

// RandomString generates a random string of the given length. It is cryptographically secure to some extent.
func RandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// BeautifulName makes the ugly name beautiful by replacing _ with spaces and using title case.
func BeautifulName(uglyName string) string {
	// Why not reuse
	return util.BeautifulName(uglyName)
}

// CanonicalName makes sure the `name` is canonical. A name is canonical if it is lowercase and all spaces are replaced with underscores.
func CanonicalName(name string) string {
	return util.CanonicalName(name)
}

// hyphaPattern is a pattern which all hypha names must match.
var hyphaPattern = regexp.MustCompile(`[^?!:#@><*|"'&%{}]+`)

var usernamePattern = regexp.MustCompile(`[^?!:#@><*|"'&%{}/]+`)

// IsCanonicalName checks if the `name` is canonical.
func IsCanonicalName(name string) bool {
	return hyphaPattern.MatchString(name)
}

// IsPossibleUsername is true if the given username is ok. Same as IsCanonicalName, but cannot have / in it and cannot be equal to "anon" or "wikimind"
func IsPossibleUsername(username string) bool {
	return username != "anon" && username != "wikimind" && usernamePattern.MatchString(strings.TrimSpace(username))
}

// HyphaNameFromRq extracts hypha name from http request. You have to also pass the action which is embedded in the url or several actions. For url /hypha/hypha, the action would be "hypha".
func HyphaNameFromRq(rq *http.Request, actions ...string) string {
	p := rq.URL.Path
	for _, action := range actions {
		if strings.HasPrefix(p, "/"+action+"/") {
			return CanonicalName(strings.TrimPrefix(p, "/"+action+"/"))
		}
	}
	log.Println("HyphaNameFromRq: this request is invalid, fall back to home hypha")
	return cfg.HomeHypha
}
