package markup

import (
	"fmt"
	"testing"
)

func TestParseStartOfEntry(t *testing.T) {
	img := ImgFromFirstLine("img {", "h")
	tests := []struct {
		line           string
		entry          imgEntry
		followedByDesc bool
	}{
		{"apple", imgEntry{"/binary/apple", "", "", ""}, false},
		{"pear|", imgEntry{"/binary/pear", "", "", ""}, false},
		{"яблоко| 30*60", imgEntry{"/binary/яблоко", "30", "60", ""}, false},
		{"груша   | 65  ", imgEntry{"/binary/груша", "65", "", ""}, false},
		{"жеронимо | 30 { full desc }", imgEntry{"/binary/жеронимо", "30", "", " full desc "}, false},
		{"жорно жованна |   *5555 {partial description", imgEntry{"/binary/жорно_жованна", "", "5555", "partial description"}, true},
		{"иноске | {full}", imgEntry{"/binary/иноске", "", "", "full"}, false},
		{"j|{partial", imgEntry{"/binary/j", "", "", "partial"}, true},
	}
	for _, triplet := range tests {
		entry, followedByDesc := img.parseStartOfEntry(triplet.line)
		if entry != triplet.entry || followedByDesc != triplet.followedByDesc {
			t.Error(fmt.Sprintf("%q:%q != %q; %v != %v", triplet.line, entry, triplet.entry, followedByDesc, triplet.followedByDesc))
		}
	}
}

func TestParseDimensions(t *testing.T) {
	tests := [][]string{
		{"500", "500", ""},
		{"3em", "3em", ""},
		{"500*", "500", ""},
		{"*500", "", "500"},
		{"800*520", "800", "520"},
		{"17%*5rem", "17%", "5rem"},
	}
	for _, triplet := range tests {
		sizeH, sizeV := parseDimensions(triplet[0])
		if sizeH != triplet[1] || sizeV != triplet[2] {
			t.Error(sizeH, "*", sizeV, " != ", triplet[1], "*", triplet[2])
		}
	}
}
