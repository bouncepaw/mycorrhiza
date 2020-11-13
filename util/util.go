package util

import (
	"net/http"
	"strings"
)

var (
	ServerPort           string
	HomePage             string
	SiteTitle            string
	WikiDir              string
	UserTree             string
	AuthMethod           string
	FixedCredentialsPath string
)

// ShorterPath is used by handlerList to display shorter path to the files. It simply strips WikiDir.
func ShorterPath(path string) string {
	if strings.HasPrefix(path, WikiDir) {
		tmp := strings.TrimPrefix(path, WikiDir)
		if tmp == "" {
			return ""
		}
		return tmp[1:]
	}
	return path
}

// HTTP200Page wraps some frequently used things for successful 200 responses.
func HTTP200Page(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(page))
}

// FindSubhyphae finds names of existing hyphae given the `hyphaIterator`.
func FindSubhyphae(hyphaName string, hyphaIterator func(func(string))) []string {
	subhyphae := make([]string, 0)
	hyphaIterator(func(otherHyphaName string) {
		if strings.HasPrefix(otherHyphaName, hyphaName+"/") {
			subhyphae = append(subhyphae, otherHyphaName)
		}
	})
	return subhyphae
}
