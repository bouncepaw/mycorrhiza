package markup

import (
	"regexp"
	"strings"

	"github.com/bouncepaw/mycorrhiza/link"
)

// OutLinks returns a channel of names of hyphae this mycodocument links.
// Links include:
// * Regular links
// * Rocketlinks
// * Transclusion
// * Image galleries
// Not needed anymore, I guess.
func (md *MycoDoc) OutLinks() chan string {
	ch := make(chan string)
	if !md.parsedAlready {
		md.Lex(0)
	}
	go func() {
		for _, line := range md.ast {
			switch v := line.contents.(type) {
			case string:
				if strings.HasPrefix(v, "<p") || strings.HasPrefix(v, "<ul class='launchpad'") {
					extractLinks(v, ch)
				}
			case Transclusion:
				ch <- v.name
			case Img:
				extractImageLinks(v, ch)
			}
		}
		close(ch)
	}()
	return ch
}

var reLinks = regexp.MustCompile(`<a href="/hypha/([^"]*)".*?</a>`)

func extractLinks(html string, ch chan string) {
	if results := reLinks.FindAllStringSubmatch(html, -1); results != nil {
		for _, result := range results {
			// result[0] is always present at this point and is not needed, because it is the whole matched substring (which we don't need)
			ch <- result[1]
		}
	}
}

func extractImageLinks(img Img, ch chan string) {
	for _, entry := range img.entries {
		if entry.srclink.Kind == link.LinkLocalHypha {
			ch <- entry.srclink.Address
		}
	}
}
