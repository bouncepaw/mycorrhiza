package markup

import (
	"fmt"
	"path"
	"strconv"
	"strings"
)

const xclError = -9

// Transclusion is used by markup parser to remember what hyphae shall be transcluded.
type Transclusion struct {
	name string
	from int // inclusive
	to   int // inclusive
}

// Transclude transcludes `xcl` and returns html representation.
func Transclude(xcl Transclusion, recursionLevel int) (html string) {
	recursionLevel++
	tmptOk := `<section class="transclusion transclusion_ok">
	<a class="transclusion__link" href="/page/%s">%s</a>
	<div class="transclusion__content">%s</div>
</section>`
	tmptFailed := `<section class="transclusion transclusion_failed">
	<p class="error">Hypha <a class="wikilink_new" href="/page/%s">%s</a> does not exist</p>
</section>`
	if xcl.from == xclError || xcl.to == xclError || xcl.from > xcl.to {
		return fmt.Sprintf(tmptFailed, xcl.name, xcl.name)
	}

	rawText, binaryHtml, err := HyphaAccess(xcl.name)
	if err != nil {
		return fmt.Sprintf(tmptFailed, xcl.name, xcl.name)
	}
	md := Doc(xcl.name, rawText)
	xclText := Parse(md.lex(), xcl.from, xcl.to, recursionLevel)
	return fmt.Sprintf(tmptOk, xcl.name, xcl.name, binaryHtml+xclText)
}

/* Grammar from hypha ‘transclusion’:
transclusion_line  ::= transclusion_token hypha_name LWS* [":" LWS* range LWS*]
transclusion_token ::= "<=" LWS+
hypha_name         ::= canonical_name | noncanonical_name
range              ::= id | (from_id two_dots to_id) | (from_id two_dots) | (two_dots to_id)
two_dots           ::= ".."
*/

func parseTransclusion(line, hyphaName string) (xclusion Transclusion) {
	line = strings.TrimSpace(remover("<=")(line))
	if line == "" {
		return Transclusion{"", xclError, xclError}
	}

	if strings.ContainsRune(line, ':') {
		parts := strings.SplitN(line, ":", 2)
		xclusion.name = xclCanonicalName(hyphaName, strings.TrimSpace(parts[0]))
		selector := strings.TrimSpace(parts[1])
		xclusion.from, xclusion.to = parseSelector(selector)
	} else {
		xclusion.name = xclCanonicalName(hyphaName, strings.TrimSpace(line))
	}
	return xclusion
}

func xclCanonicalName(hyphaName, xclName string) string {
	switch {
	case strings.HasPrefix(xclName, "./"):
		return canonicalName(path.Join(hyphaName, strings.TrimPrefix(xclName, "./")))
	case strings.HasPrefix(xclName, "../"):
		return canonicalName(path.Join(path.Dir(hyphaName), strings.TrimPrefix(xclName, "../")))
	default:
		return canonicalName(xclName)
	}
}

// At this point:
// selector ::= id
//            | from ".."
//            | from ".." to
//            |      ".." to
// If it is not, return (xclError, xclError).
func parseSelector(selector string) (from, to int) {
	if selector == "" {
		return 0, 0
	}
	if strings.Contains(selector, "..") {
		parts := strings.Split(selector, "..")

		var (
			fromStr       = strings.TrimSpace(parts[0])
			from, fromErr = strconv.Atoi(fromStr)
			toStr         = strings.TrimSpace(parts[1])
			to, toErr     = strconv.Atoi(toStr)
		)
		if fromStr == "" && toStr == "" {
			return 0, 0
		}
		if fromErr == nil || toErr == nil {
			return from, to
		}
	} else if id, err := strconv.Atoi(selector); err == nil {
		return id, id
	}
	return xclError, xclError
}
