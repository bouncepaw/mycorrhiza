package markup

import ()

const maxRecursionLevel = 3

type GemParserState struct {
	recursionLevel int
}

func Parse(ast []Line, from, to int, state GemParserState) (html string) {
	if state.recursionLevel > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		if line.id >= from && (line.id <= to || to == 0) || line.id == -1 {
			switch v := line.contents.(type) {
			case Transclusion:
				html += Transclude(v, state)
			case Img:
				html += v.ToHtml()
			case string:
				html += v
			default:
				html += "Unknown"
			}
		}
	}
	return html
}

func ToHtml(name, text string) string {
	state := GemParserState{}
	return Parse(lex(name, text), 0, 0, state)
}
