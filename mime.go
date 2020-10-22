package main

import (
	"path/filepath"
	"strings"
)

// TextType is content type of text part of a hypha.
type TextType int

const (
	// TextPlain is default text content type.
	TextPlain TextType = iota
	// TextGemini is content type for MycorrhizaWiki's dialect of gemtext.
	TextGemini
)

// Mime returns mime type representation of `t`.
func (t TextType) Mime() string {
	return [...]string{"text/plain", "text/gemini"}[t]
}

// Extension returns extension (with dot) to be used for files with content type `t`.
func (t TextType) Extension() string {
	return [...]string{".txt", ".gmi"}[t]
}

// BinaryType is content type of binary part of a hypha.
type BinaryType int

// Supported binary content types
const (
	// BinaryOctet is default binary content type.
	BinaryOctet BinaryType = iota
	BinaryJpeg
	BinaryGif
	BinaryPng
	BinaryWebp
	BinarySvg
	BinaryIco
	BinaryOgg
	BinaryWebm
	BinaryMp3
	BinaryMp4
)

var binaryMimes = [...]string{}

// Mime returns mime type representation of `t`.
func (t BinaryType) Mime() string {
	return binaryMimes[t]
}

var binaryExtensions = [...]string{
	".bin", ".jpg", ".gif", ".png", ".webp", ".svg", ".ico",
	".ogg", ".webm", ".mp3", ".mp4",
}

// Extension returns extension (with dot) to be used for files with content type `t`.
func (t BinaryType) Extension() string {
	return binaryExtensions[t]
}

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

// MimeToBinaryType converts mime type to BinaryType. If the mime type is not supported, BinaryOctet is returned as a fallback type.
func MimeToBinaryType(mime string) BinaryType {
	for i, binaryMime := range binaryMimes {
		if binaryMime == mime {
			return BinaryType(i)
		}
	}
	return BinaryOctet
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

// mimeData determines what content type file has judging by its `ext`ension. `itText` and `mimeId` are the same as in DataFromFilename.
func mimeData(ext string) (isText bool, mimeId int) {
	switch ext {
	case ".txt":
		return true, int(TextPlain)
	case ".gmi":
		return true, int(TextGemini)
	}
	for i, binExt := range binaryExtensions {
		if ext == binExt {
			return false, i
		}
	}
	return false, 0
}
