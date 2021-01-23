package util

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

var (
	URL                  string
	ServerPort           string
	HomePage             string
	SiteNavIcon          string
	SiteName             string
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
