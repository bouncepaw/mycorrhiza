package main

import (
	"testing"
)

func TestMimeData(t *testing.T) {
	check := func(ext string, expectedIsText bool, expectedMimeId int) {
		isText, mimeId := mimeData(ext)
		if isText != expectedIsText || mimeId != expectedMimeId {
			t.Error(ext, isText, mimeId)
		}
	}
	check(".txt", true, int(TextPlain))
	check(".gmi", true, int(TextGemini))
	check(".bin", false, int(BinaryOctet))
	check(".jpg", false, int(BinaryJpeg))
	check(".bin", false, int(BinaryOctet))
}
