package util

import (
	"net/http"
	"strings"
)

var WikiDir string

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