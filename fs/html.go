package fs

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func markdownToHtml(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.ToHTML([]byte(md), p, renderer))
}

func (h *Hypha) asHtml() (string, error) {
	rev := h.actual
	ret := `<article class="page">
	<h1 class="page__title">` + rev.FullName + `</h1>
`
	// What about using <figure>?
	if h.hasBinaryData() {
		ret += fmt.Sprintf(`<img src="/%s?action=binary&rev=%d" class="page__amnt"/>`, rev.FullName, rev.Id)
	}

	contents, err := ioutil.ReadFile(rev.TextPath)
	if err != nil {
		log.Println("Failed to render", rev.FullName, ":", err)
		return "", err
	}

	// TODO: support more markups.
	// TODO: support mycorrhiza extensions like transclusion.
	switch rev.TextMime {
	case "text/markdown":
		html := markdown.ToHTML(contents, nil, nil)
		ret += string(html)
	default:
		ret += fmt.Sprintf(`<pre>%s</pre>`, contents)
	}

	ret += `
</article>`

	return ret, nil
}
