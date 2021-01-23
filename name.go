package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bouncepaw/mycorrhiza/util"
)

// isCanonicalName checks if the `name` is canonical.
func isCanonicalName(name string) bool {
	return HyphaPattern.MatchString(name)
}

// CanonicalName makes sure the `name` is canonical. A name is canonical if it is lowercase and all spaces are replaced with underscores.
func CanonicalName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

// naviTitle turns `canonicalName` into html string with each hypha path parts higlighted as links.
// TODO: rework as a template
func naviTitle(canonicalName string) string {
	var (
		html = fmt.Sprintf(`<h1 class="navi-title" id="navi-title">
	<a href="/page/%s">%s</a><span aria-hidden="true" class="navi-title__colon">:</span>`, util.HomePage, util.SiteNavIcon)
		prevAcc = `/page/`
		parts   = strings.Split(canonicalName, "/")
		rel     = "up"
	)
	for i, part := range parts {
		if i > 0 {
			html += `<span aria-hidden="true" class="navi-title__separator">/</span>`
		}
		if i == len(parts)-1 {
			rel = "bookmark"
		}
		html += fmt.Sprintf(
			`<a href="%s" rel="%s">%s</a>`,
			prevAcc+part,
			rel,
			util.BeautifulName(part),
		)
		prevAcc += part + "/"
	}
	return html + "</h1>"
}

// HyphaNameFromRq extracts hypha name from http request. You have to also pass the action which is embedded in the url. For url /page/hypha, the action would be "page".
func HyphaNameFromRq(rq *http.Request, action string) string {
	return CanonicalName(strings.TrimPrefix(rq.URL.Path, "/"+action+"/"))
}
