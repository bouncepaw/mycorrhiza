package markup

import (
	"strings"

	"github.com/bouncepaw/mycorrhiza/link"
)

type imgEntry struct {
	srclink *link.Link
	path    strings.Builder
	sizeW   strings.Builder
	sizeH   strings.Builder
	desc    strings.Builder
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
