package markup

import (
	"fmt"
	"regexp"
	"strings"
)

var imgRe = regexp.MustCompile(`^img\s+{`)

func MatchesImg(line string) bool {
	return imgRe.MatchString(line)
}

type imgEntry struct {
	path  string
	sizeH string
	sizeV string
	desc  string
}

type Img struct {
	entries   []imgEntry
	inDesc    bool
	hyphaName string
}

func (img *Img) Process(line string) (shouldGoBackToNormal bool) {
	if img.inDesc {
		rightBraceIndex := strings.IndexRune(line, '}')
		if cnt := len(img.entries); rightBraceIndex == -1 && cnt != 0 {
			img.entries[cnt-1].desc += "\n" + line
		} else if rightBraceIndex != -1 && cnt != 0 {
			img.entries[cnt-1].desc += "\n" + line[:rightBraceIndex]
			img.inDesc = false
		}
		if strings.Count(line, "}") > 1 {
			return true
		}
	} else if s := strings.TrimSpace(line); s != "" {
		if s[0] == '}' {
			return true
		}
		img.parseStartOfEntry(line)
	}
	return false
}

func ImgFromFirstLine(line, hyphaName string) Img {
	img := Img{
		hyphaName: hyphaName,
		entries:   make([]imgEntry, 0),
	}
	line = line[strings.IndexRune(line, '{'):]
	if len(line) == 1 { // if { only
	} else {
		line = line[1:] // Drop the {
	}
	return img
}

func (img *Img) canonicalPathFor(path string) string {
	path = strings.TrimSpace(path)
	if strings.IndexRune(path, ':') != -1 || strings.IndexRune(path, '/') == 0 {
		return path
	} else {
		return "/binary/" + xclCanonicalName(img.hyphaName, path)
	}
}

func (img *Img) parseStartOfEntry(line string) (entry imgEntry, followedByDesc bool) {
	pipeIndex := strings.IndexRune(line, '|')
	if pipeIndex == -1 { // If no : in string
		entry.path = img.canonicalPathFor(line)
	} else {
		entry.path = img.canonicalPathFor(line[:pipeIndex])
		line = strings.TrimPrefix(line, line[:pipeIndex+1])

		var (
			leftBraceIndex  = strings.IndexRune(line, '{')
			rightBraceIndex = strings.IndexRune(line, '}')
			dimensions      string
		)

		if leftBraceIndex == -1 {
			dimensions = line
		} else {
			dimensions = line[:leftBraceIndex]
		}

		sizeH, sizeV := parseDimensions(dimensions)
		entry.sizeH = sizeH
		entry.sizeV = sizeV

		if leftBraceIndex != -1 && rightBraceIndex == -1 {
			img.inDesc = true
			followedByDesc = true
			entry.desc = strings.TrimPrefix(line, line[:leftBraceIndex+1])
		} else if leftBraceIndex != -1 && rightBraceIndex != -1 {
			entry.desc = line[leftBraceIndex+1 : rightBraceIndex]
		}
	}
	img.entries = append(img.entries, entry)
	return
}

func parseDimensions(dimensions string) (sizeH, sizeV string) {
	xIndex := strings.IndexRune(dimensions, '*')
	if xIndex == -1 { // If no x in dimensions
		sizeH = strings.TrimSpace(dimensions)
	} else {
		sizeH = strings.TrimSpace(dimensions[:xIndex])
		sizeV = strings.TrimSpace(strings.TrimPrefix(dimensions, dimensions[:xIndex+1]))
	}
	return
}

func (img Img) ToHtml() (html string) {
	for _, entry := range img.entries {
		html += fmt.Sprintf(`<figure>
	<img src="%s" width="%s" height="%s">
`, entry.path, entry.sizeH, entry.sizeV)
		if entry.desc != "" {
			html += `	<figcaption>`
			for i, line := range strings.Split(entry.desc, "\n") {
				if line != "" {
					if i > 0 {
						html += `<br>`
					}
					html += ParagraphToHtml(line)
				}
			}
			html += `</figcaption>`
		}
		html += `</figure>`
	}
	return `<section class="img-gallery">
` + html + `
</section>`
}
