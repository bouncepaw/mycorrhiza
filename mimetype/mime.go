package mimetype

import (
	"path/filepath"
	"strings"

	"github.com/bouncepaw/mycorrhiza/util"
)

// ToExtension returns dotted extension for given mime-type.
func ToExtension(mime string) string {
	if ext, ok := mapMime2Ext[mime]; ok {
		return "." + ext
	}
	return ".bin"
}

// FromExtension returns mime-type for given extension. The extension must start with a dot.
func FromExtension(ext string) string {
	if mime, ok := mapExt2Mime[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

// DataFromFilename fetches all meta information from hypha content file with path `fullPath`. If it is not a content file, `skip` is true, and you are expected to ignore this file when indexing hyphae. `name` is name of the hypha to which this file relates. `isText` is true when the content file is text, false when is binary.
func DataFromFilename(fullPath string) (name string, isText bool, skip bool) {
	shortPath := util.ShorterPath(fullPath)
	ext := filepath.Ext(shortPath)
	name = util.CanonicalName(strings.TrimSuffix(shortPath, ext))
	switch ext {
	case ".myco":
		isText = true
	case "", shortPath:
		skip = true
	}

	return
}

var mapMime2Ext = map[string]string{
	"application/octet-stream": "bin",
	"image/jpeg":               "jpg",
	"image/gif":                "gif",
	"image/png":                "png",
	"image/webp":               "webp",
	"image/svg+xml":            "svg",
	"image/x-icon":             "ico",
	"application/ogg":          "ogg",
	"video/webm":               "webm",
	"audio/mp3":                "mp3",
	"video/mp4":                "mp4",
}

var mapExt2Mime = map[string]string{
	".bin":  "application/octet-stream",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".png":  "image/png",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",
	".ogg":  "application/ogg",
	".webm": "video/webm",
	".mp3":  "audio/mp3",
	".mp4":  "video/mp4",
}
