package mimetype

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/bouncepaw/mycorrhiza/util"
)

// ToExtension returns dotted extension for given mime-type.
func ToExtension(mimeType string) string {
	if exts, err := mime.ExtensionsByType(mimeType); err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ".bin"
}

// FromExtension returns mime-type for given extension. The extension must start with a dot.
func FromExtension(ext string) string {
	if mimeType := mime.TypeByExtension(strings.ToLower(ext)); mimeType != "" {
		return mimeType
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
