package util

import (
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
