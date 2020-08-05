package gemtext

import (
	"fmt"
	"html"
	"path"
	"strings"
)

// HyphaExists holds function that checks that a hypha is present.
var HyphaExists func(string) bool

// HyphaAccess holds function that accesses a hypha by its name.
var HyphaAccess func(string) (rawText, binaryHtml string, err error)

// GemLexerState is used by gemtext parser to remember what is going on.
type GemLexerState struct {
	// Name of hypha being parsed
	name  string
	where string // "", "list", "pre"
	// Line id
	id  int
	buf string
}

// GeminiToHtml converts gemtext `content` of hypha `name` to html string.
func GeminiToHtml(name, content string) string {
	return "TODO: do"
}

type Line struct {
	id int
	// interface{} may be bad. What I need is a sum of string and Transclusion
	contents interface{}
}

// Parse gemtext line starting with "=>" according to wikilink rules.
// See http://localhost:1737/page/wikilink
func wikilink(src string, state *GemLexerState) (href, text, class string) {
	src = strings.TrimSpace(remover("=>")(src))
	if src == "" {
		return
	}
	// Href is text after => till first whitespace
	href = strings.Fields(src)[0]
	// Text is everything after whitespace.
	// If there's no text, make it same as href
	if text = strings.TrimPrefix(src, href); text == "" {
		text = href
	}

	class = "wikilink_internal"

	switch {
	case strings.HasPrefix(href, "./"):
		hyphaName := canonicalName(path.Join(
			state.name, strings.TrimPrefix(href, "./")))
		if !HyphaExists(hyphaName) {
			class = "wikilink_new"
		}
		href = path.Join("/page", hyphaName)
	case strings.HasPrefix(href, "../"):
		hyphaName := canonicalName(path.Join(
			path.Dir(state.name), strings.TrimPrefix(href, "../")))
		if !HyphaExists(hyphaName) {
			class = "wikilink_new"
		}
		href = path.Join("/page", hyphaName)
	case strings.HasPrefix(href, "/"):
	case strings.ContainsRune(href, ':'):
		class = "wikilink_external"
	default:
		href = path.Join("/page", href)
	}
	return href, strings.TrimSpace(text), class
}

func lex(name, content string) (ast []Line) {
	var state = GemLexerState{name: name}

	for _, line := range strings.Split(content, "\n") {
		geminiLineToAST(line, &state, &ast)
	}
	return ast
}

// Lex `line` in gemtext and save it to `ast` using `state`.
func geminiLineToAST(line string, state *GemLexerState, ast *[]Line) {
	if "" == strings.TrimSpace(line) {
		return
	}

	startsWith := func(token string) bool {
		return strings.HasPrefix(line, token)
	}
	addLine := func(text interface{}) {
		*ast = append(*ast, Line{id: state.id, contents: text})
	}

	// Beware! Usage of goto. Some may say it is considered evil but in this case it helped to make a better-structured code.
	switch state.where {
	case "pre":
		goto preformattedState
	case "list":
		goto listState
	default:
		goto normalState
	}

preformattedState:
	switch {
	case startsWith("```"):
		state.where = ""
		state.buf = strings.TrimSuffix(state.buf, "\n")
		addLine(state.buf + "</code></pre>")
		state.buf = ""
	default:
		state.buf += html.EscapeString(line) + "\n"
	}
	return

listState:
	switch {
	case startsWith("*"):
		state.buf += fmt.Sprintf("\t<li>%s</li>\n", remover("*")(line))
	case startsWith("```"):
		state.where = "pre"
		addLine(state.buf + "</ul>")
		state.id++
		state.buf = fmt.Sprintf("<pre id='%d' alt='%s' class='codeblock'><code>", state.id, strings.TrimPrefix(line, "```"))
	default:
		state.where = ""
		addLine(state.buf + "</ul>")
		goto normalState
	}
	return

normalState:
	state.id++
	switch {

	case startsWith("```"):
		state.where = "pre"
		state.buf = fmt.Sprintf("<pre id='%d' alt='%s' class='codeblock'><code>", state.id, strings.TrimPrefix(line, "```"))
	case startsWith("*"):
		state.where = "list"
		state.buf = fmt.Sprintf("<ul id='%d'>\n", state.id)
		goto listState

	case startsWith("###"):
		addLine(fmt.Sprintf(
			"<h3 id='%d'>%s</h3>", state.id, removeHeadingOctothorps(line)))
	case startsWith("##"):
		addLine(fmt.Sprintf(
			"<h2 id='%d'>%s</h2>", state.id, removeHeadingOctothorps(line)))
	case startsWith("#"):
		addLine(fmt.Sprintf(
			"<h1 id='%d'>%s</h1>", state.id, removeHeadingOctothorps(line)))

	case startsWith(">"):
		addLine(fmt.Sprintf(
			"<blockquote id='%d'>%s</blockquote>", state.id, remover(">")(line)))
	case startsWith("=>"):
		source, content, class := wikilink(line, state)
		addLine(fmt.Sprintf(
			`<p><a id='%d' class='%s' href="%s">%s</a></p>`, state.id, class, source, content))

	case startsWith("<="):
		addLine(parseTransclusion(line, state.name))
	default:
		addLine(fmt.Sprintf("<p id='%d'>%s</p>", state.id, line))
	}
}
