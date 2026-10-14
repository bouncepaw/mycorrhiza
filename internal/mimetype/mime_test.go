package mimetype

import (
	"fmt"
	"testing"
)

func TestToExtension(t *testing.T) {
	var tests = []struct {
		mimeType string
		want     string
	}{
		{"application/octet-stream", ".bin"},
		{"image/jpeg", ".jpg"},
		{"image/gif", ".gif"},
		{"image/png", ".png"},
		{"image/webp", ".webp"},
		{"image/svg+xml", ".svg"},
		{"image/x-icon", ".ico"},
		{"application/ogg", ".ogg"},
		{"video/webm", ".webm"},
		{"audio/mp3", ".mp3"},
		{"audio/mpeg", ".mp3"},
		{"audio/mpeg3", ".mp3"},
		{"video/mp4", ".mp4"},
		{"audio/flac", ".flac"},
		{"audio/wav", ".wav"},
		{"audio/vnd.wav", ".wav"},
		{"audio/vnd.wave", ".wav"},
		{"audio/wave", ".wav"},
		{"audio/x-pn-wav", ".wav"},
		{"audio/x-wav", ".wav"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.mimeType)

		t.Run(testname, func(t *testing.T) {
			ans := ToExtension(tt.mimeType)

			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

func TestFromExtension(t *testing.T) {
	var tests = []struct {
		ext  string
		want string
	}{
		{".bin", "application/octet-stream"},
		{".jpg", "image/jpeg"},
		{".jpeg", "image/jpeg"},
		{".gif", "image/gif"},
		{".png", "image/png"},
		{".webp", "image/webp"},
		{".svg", "image/svg+xml"},
		{".ico", "image/x-icon"},
		{".ogg", "application/ogg"},
		{".webm", "video/webm"},
		{".mp3", "audio/mpeg"},
		{".mp4", "video/mp4"},
		{".flac", "audio/flac"},
		{".wav", "audio/wav"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.ext)

		t.Run(testname, func(t *testing.T) {
			ans := FromExtension(tt.ext)

			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
