package main

import (
	"fmt"
	"log"
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

// HyphaNameFromRq extracts hypha name from http request. You have to also pass the action which is embedded in the url or several actions. For url /hypha/hypha, the action would be "hypha".
func HyphaNameFromRq(rq *http.Request, actions ...string) string {
	p := rq.URL.Path
	for _, action := range actions {
		if strings.HasPrefix(p, "/"+action+"/") {
			return util.CanonicalName(strings.TrimPrefix(p, "/"+action+"/"))
		}
	}
	log.Fatal("HyphaNameFromRq: no matching action passed")
	return ""
}
