package mimetype

import (
	"log/slog"
	"mime"
	"path/filepath"
	"strings"

	"golang.org/x/exp/collection/set"

	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/util/fp"
)

// ToExtension returns dotted extension for given mime-type.
func ToExtension(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/x-icon":
		return ".ico"
	case "application/ogg":
		return ".ogg"
	case "audio/mp3", "audio/mpeg3", "audio/mpeg":
		return ".mp3"
	case "video/mp4":
		return ".mp4"
	case "audio/vnd.wav", "audio/vnd.wave", "audio/wave", "audio/x-pn-wav", "audio/x-wav":
		return ".wav"
	}
	if exts, err := mime.ExtensionsByType(mimeType); err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ".bin"
}

// FromExtension returns mime-type for given extension. The extension must start with a dot.
func FromExtension(ext string) string {
	switch ext {
	case ".ico":
		return "image/x-icon"
	case ".ogg":
		return "application/ogg"
	case ".flac":
		return "audio/flac"
	case ".wav":
		return "audio/wav"
	}
	if mimeType := mime.TypeByExtension(strings.ToLower(ext)); mimeType != "" {
		return mimeType
	}
	return "application/octet-stream"
}

// Defined by the WHATWG MIME Sniffing Standard:
// => https://mimesniff.spec.whatwg.org/#javascript-mime-type
// With X-Content-Type-Options: nosniff, browsers execute a script only if its response has one of these types.
var javaScriptTypes = set.Of(
	"application/ecmascript",
	"application/javascript",
	"application/x-ecmascript",
	"application/x-javascript",
	"text/ecmascript",
	"text/javascript",
	"text/javascript1.0",
	"text/javascript1.1",
	"text/javascript1.2",
	"text/javascript1.3",
	"text/javascript1.4",
	"text/javascript1.5",
	"text/jscript",
	"text/livescript",
	"text/x-ecmascript",
	"text/x-javascript",
)

func IsJavaScript(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		slog.Warn("Failed to parse media type; considering it as JavaScript. Better safe than sorry",
			"type", contentType)
		return true
	}
	return javaScriptTypes.Contains(mediaType)
}

func ForServing(ext string) string {
	contentType := FromExtension(ext)
	return fp.If(IsJavaScript(contentType), "text/plain; charset=utf-8", contentType)
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
