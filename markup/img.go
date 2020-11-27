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
	trimmedPath string
	path        strings.Builder
	sizeW       strings.Builder
	sizeH       strings.Builder
	desc        strings.Builder
}

func (entry *imgEntry) descriptionAsHtml(hyphaName string) (html string) {
	if entry.desc.Len() == 0 {
		return ""
	}
	lines := strings.Split(entry.desc.String(), "\n")
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			if html != "" {
				html += `<br>`
			}
			html += ParagraphToHtml(hyphaName, line)
		}
	}
	return `<figcaption>` + html + `</figcaption>`
}

func (entry *imgEntry) sizeWAsAttr() string {
	if entry.sizeW.Len() == 0 {
		return ""
	}
	return ` width="` + entry.sizeW.String() + `"`
}

func (entry *imgEntry) sizeHAsAttr() string {
	if entry.sizeH.Len() == 0 {
		return ""
	}
	return ` height="` + entry.sizeH.String() + `"`
}

type imgState int

const (
	inRoot imgState = iota
	inName
	inDimensionsW
	inDimensionsH
	inDescription
)

type Img struct {
	entries   []imgEntry
	currEntry imgEntry
	hyphaName string
	state     imgState
}

func (img *Img) pushEntry() {
	if strings.TrimSpace(img.currEntry.path.String()) != "" {
		img.entries = append(img.entries, img.currEntry)
		img.currEntry = imgEntry{}
		img.currEntry.path.Reset()
	}
}

func (img *Img) Process(line string) (shouldGoBackToNormal bool) {
	stateToProcessor := map[imgState]func(rune) bool{
		inRoot:        img.processInRoot,
		inName:        img.processInName,
		inDimensionsW: img.processInDimensionsW,
		inDimensionsH: img.processInDimensionsH,
		inDescription: img.processInDescription,
	}
	for _, r := range line {
		if shouldReturnTrue := stateToProcessor[img.state](r); shouldReturnTrue {
			return true
		}
	}
	return false
}

func (img *Img) processInDescription(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.state = inName
	default:
		img.currEntry.desc.WriteRune(r)
	}
	return false
}

func (img *Img) processInRoot(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case '\n', '\r':
		img.pushEntry()
	case ' ', '\t':
	default:
		img.state = inName
		img.currEntry = imgEntry{}
		img.currEntry.path.Reset()
		img.currEntry.path.WriteRune(r)
	}
	return false
}

func (img *Img) processInName(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case '|':
		img.state = inDimensionsW
	case '{':
		img.state = inDescription
	case '\n', '\r':
		img.pushEntry()
		img.state = inRoot
	default:
		img.currEntry.path.WriteRune(r)
	}
	return false
}

func (img *Img) processInDimensionsW(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case '*':
		img.state = inDimensionsH
	case ' ', '\t', '\n':
	case '{':
		img.state = inDescription
	default:
		img.currEntry.sizeW.WriteRune(r)
	}
	return false
}

func (img *Img) processInDimensionsH(r rune) (shouldGoBackToNormal bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case ' ', '\t', '\n':
	case '{':
		img.state = inDescription
	default:
		img.currEntry.sizeH.WriteRune(r)
	}
	return false
}

func ImgFromFirstLine(line, hyphaName string) (img *Img, shouldGoBackToNormal bool) {
	img = &Img{
		hyphaName: hyphaName,
		entries:   make([]imgEntry, 0),
	}
	line = line[strings.IndexRune(line, '{')+1:]
	return img, img.Process(line)
}

func (img *Img) binaryPathFor(path string) string {
	path = strings.TrimSpace(path)
	if strings.IndexRune(path, ':') != -1 || strings.IndexRune(path, '/') == 0 {
		return path
	} else {
		return "/binary/" + xclCanonicalName(img.hyphaName, path)
	}
}

func (img *Img) pagePathFor(path string) string {
	path = strings.TrimSpace(path)
	if strings.IndexRune(path, ':') != -1 || strings.IndexRune(path, '/') == 0 {
		return path
	} else {
		return "/page/" + xclCanonicalName(img.hyphaName, path)
	}
}

func parseDimensions(dimensions string) (sizeW, sizeH string) {
	xIndex := strings.IndexRune(dimensions, '*')
	if xIndex == -1 { // If no x in dimensions
		sizeW = strings.TrimSpace(dimensions)
	} else {
		sizeW = strings.TrimSpace(dimensions[:xIndex])
		sizeH = strings.TrimSpace(strings.TrimPrefix(dimensions, dimensions[:xIndex+1]))
	}
	return
}

func (img *Img) checkLinks() map[string]bool {
	m := make(map[string]bool)
	for i, entry := range img.entries {
		// Also trim them for later use
		entry.trimmedPath = strings.TrimSpace(entry.path.String())
		isAbsoluteUrl := strings.ContainsRune(entry.trimmedPath, ':')
		if !isAbsoluteUrl {
			entry.trimmedPath = canonicalName(entry.trimmedPath)
		}
		img.entries[i] = entry
		m[entry.trimmedPath] = isAbsoluteUrl
	}
	HyphaIterate(func(hyphaName string) {
		for _, entry := range img.entries {
			if hyphaName == entry.trimmedPath {
				m[entry.trimmedPath] = true
			}
		}
	})
	return m
}

func (img *Img) ToHtml() (html string) {
	linkAvailabilityMap := img.checkLinks()
	isOneImageOnly := len(img.entries) == 1 && img.entries[0].desc.Len() == 0
	if isOneImageOnly {
		html += `<section class="img-gallery img-gallery_one-image">`
	} else {
		html += `<section class="img-gallery img-gallery_many-images">`
	}

	for _, entry := range img.entries {
		html += `<figure>`
		// If is existing hypha or an external path
		if linkAvailabilityMap[entry.trimmedPath] {
			html += fmt.Sprintf(
				`<a href="%s"><img src="%s" %s %s></a>`,
				img.pagePathFor(entry.trimmedPath),
				img.binaryPathFor(entry.trimmedPath),
				entry.sizeWAsAttr(), entry.sizeHAsAttr())
		} else { // If is a non-existent hypha
			html += fmt.Sprintf(`<a class="wikilink_new" href="%s">Hypha <em>%s</em> does not exist</a>`, img.pagePathFor(entry.trimmedPath), entry.trimmedPath)
		}
		html += entry.descriptionAsHtml(img.hyphaName)
		html += `</figure>`
	}
	return html + `</section>`
}
