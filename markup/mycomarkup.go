// This is not done yet
package markup

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycorrhiza/cfg"

	"github.com/bouncepaw/mycomarkup/links"
)

// A Mycomarkup-formatted document
type MycoDoc struct {
	// data
	hyphaName string
	contents  string
	// indicators
	parsedAlready bool
	// results
	ast           []Line
	html          string
	firstImageURL string
	description   string
}

// Constructor
func Doc(hyphaName, contents string) *MycoDoc {
	md := &MycoDoc{
		hyphaName: hyphaName,
		contents:  contents,
	}
	return md
}

func (md *MycoDoc) Lex(recursionLevel int) *MycoDoc {
	if !md.parsedAlready {
		md.ast = md.lex()
	}
	md.parsedAlready = true
	return md
}

// AsHTML returns an html representation of the document
func (md *MycoDoc) AsHTML() string {
	md.html = Parse(md.Lex(0).ast, 0, 0, 0)
	return md.html
}

// AsGemtext returns a gemtext representation of the document. Currently really limited, just returns source text
func (md *MycoDoc) AsGemtext() string {
	return md.contents
}

// Used to clear opengraph description from html tags. This method is usually bad because of dangers of malformed HTML, but I'm going to use it only for Mycorrhiza-generated HTML, so it's okay. The question mark is required; without it the whole string is eaten away.
var htmlTagRe = regexp.MustCompile(`<.*?>`)

// OpenGraphHTML returns an html representation of og: meta tags.
func (md *MycoDoc) OpenGraphHTML() string {
	md.ogFillVars()
	return strings.Join([]string{
		ogTag("title", md.hyphaName),
		ogTag("type", "article"),
		ogTag("image", md.firstImageURL),
		ogTag("url", cfg.URL+"/hypha/"+md.hyphaName),
		ogTag("determiner", ""),
		ogTag("description", htmlTagRe.ReplaceAllString(md.description, "")),
	}, "\n")
}

func (md *MycoDoc) ogFillVars() *MycoDoc {
	md.firstImageURL = cfg.URL + "/favicon.ico"
	foundDesc := false
	foundImg := false
	for _, line := range md.ast {
		switch v := line.contents.(type) {
		case string:
			if !foundDesc {
				md.description = v
				foundDesc = true
			}
		case Img:
			if !foundImg && len(v.entries) > 0 {
				md.firstImageURL = v.entries[0].srclink.ImgSrc()
				if !v.entries[0].srclink.OfKind(links.LinkExternal) {
					md.firstImageURL = cfg.URL + md.firstImageURL
				}
				foundImg = true
			}
		}
	}
	return md
}

func ogTag(property, content string) string {
	return fmt.Sprintf(`<meta property="og:%s" content="%s"/>`, property, content)
}
