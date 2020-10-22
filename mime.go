package main

import (
	"path/filepath"
	"strings"
)

func MimeToExtension(mime string) string {
	mm := map[string]string{
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
	if ext, ok := mm[mime]; ok {
		return "." + ext
	}
	return ".bin"
}

func ExtensionToMime(ext string) string {
	mm := map[string]string{
		"bin":  "application/octet-stream",
		"jpg":  "image/jpeg",
		"gif":  "image/gif",
		"png":  "image/png",
		"webp": "image/webp",
		"svg":  "image/svg+xml",
		"ico":  "image/x-icon",
		"ogg":  "application/ogg",
		"webm": "video/webm",
		"mp3":  "audio/mp3",
		"mp4":  "video/mp4",
	}
	if mime, ok := mm[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

// DataFromFilename fetches all meta information from hypha content file with path `fullPath`. If it is not a content file, `skip` is true, and you are expected to ignore this file when indexing hyphae. `name` is name of the hypha to which this file relates. `isText` is true when the content file is text, false when is binary. `mimeId` is an integer representation of content type. Cast it to TextType if `isText == true`, cast it to BinaryType if `isText == false`.
func DataFromFilename(fullPath string) (name string, isText bool, skip bool) {
	shortPath := strings.TrimPrefix(fullPath, WikiDir)[1:]
	ext := filepath.Ext(shortPath)
	name = CanonicalName(strings.TrimSuffix(shortPath, ext))
	switch ext {
	case ".myco":
		isText = true
	case "", shortPath:
		skip = true
	}

	return
}
