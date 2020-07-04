package plugin

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/plugin/parser"
)

func ParserForMime(mime string) func([]byte) string {
	parsers := map[string]func([]byte) string{
		"text/markdown": parser.MarkdownToHtml,
		"text/creole":   parser.CreoleToHtml,
		"text/gemini":   parser.GeminiToHtml,
	}
	if parserFunc, ok := parsers[mime]; ok {
		return parserFunc
	}
	return func(contents []byte) string {
		return fmt.Sprintf(`<pre><code>%s</code></pre>`, contents)
	}
}
