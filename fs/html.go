package fs

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bouncepaw/mycorrhiza/plugin"
	"github.com/bouncepaw/mycorrhiza/util"
)

func (h *Hypha) asHtml() (string, error) {
	rev := h.actual
	ret := `<article class="page">
	<h1 class="page__title">` + rev.FullName + `</h1>
`
	// What about using <figure>?
	// TODO: support other things
	if h.hasBinaryData() {
		ret += fmt.Sprintf(`<img src="/:%s?action=binary&rev=%d" class="page__amnt"/>`, util.DisplayToCanonical(rev.FullName), rev.Id)
	}

	contents, err := ioutil.ReadFile(rev.TextPath)
	if err != nil {
		log.Println("Failed to read contents of", rev.FullName, ":", err)
		return "", err
	}

	// TODO: support more markups.
	// TODO: support mycorrhiza extensions like transclusion.
	parser := plugin.ParserForMime(rev.TextMime)
	ret += parser(contents)

	ret += `
</article>`

	return ret, nil
}
