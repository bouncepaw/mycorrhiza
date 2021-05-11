package util

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

func PrepareRq(rq *http.Request) {
	log.Println(rq.RequestURI)
	rq.URL.Path = strings.TrimSuffix(rq.URL.Path, "/")
}

// LettersNumbersOnly keeps letters and numbers only in the given string.
func LettersNumbersOnly(s string) string {
	var (
		ret            strings.Builder
		usedUnderscore bool
	)
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			ret.WriteRune(r)
			usedUnderscore = false
		} else if !usedUnderscore {
			ret.WriteRune('_')
			usedUnderscore = true
		}
	}
	return strings.Trim(ret.String(), "_")
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
	w.Write([]byte(page))
}

// HTTP200Page wraps some frequently used things for successful 200 responses.
func HTTP200Page(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(page))
}

func RandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Strip hypha name from all ancestor names, replace _ with spaces, title case
func BeautifulName(uglyName string) string {
	if uglyName == "" {
		return uglyName
	}
	return strings.Title(strings.ReplaceAll(uglyName, "_", " "))
}

// CanonicalName makes sure the `name` is canonical. A name is canonical if it is lowercase and all spaces are replaced with underscores.
func CanonicalName(name string) string {
	return strings.ToLower(
		strings.ReplaceAll(
			strings.TrimRight(
				strings.TrimLeft(name, "_"),
				"_",
			), " ", "_"))
}

// HyphaPattern is a pattern which all hyphae must match.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"'&%{}]+`)

var UsernamePattern = regexp.MustCompile(`[^?!:#@><*|"'&%{}/]+`)

// IsCanonicalName checks if the `name` is canonical.
func IsCanonicalName(name string) bool {
	return HyphaPattern.MatchString(name)
}

func IsPossibleUsername(username string) bool {
	return UsernamePattern.MatchString(strings.TrimSpace(username))
}

// HyphaNameFromRq extracts hypha name from http request. You have to also pass the action which is embedded in the url or several actions. For url /hypha/hypha, the action would be "hypha".
func HyphaNameFromRq(rq *http.Request, actions ...string) string {
	p := rq.URL.Path
	for _, action := range actions {
		if strings.HasPrefix(p, "/"+action+"/") {
			return CanonicalName(strings.TrimPrefix(p, "/"+action+"/"))
		}
	}
	log.Println("HyphaNameFromRq: this request is invalid, fallback to home hypha")
	return cfg.HomeHypha
}
