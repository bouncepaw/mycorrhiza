package markup

const maxRecursionLevel = 3

func Parse(ast []Line, from, to int, recursionLevel int) (html string) {
	if recursionLevel > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		if line.id >= from && (line.id <= to || to == 0) || line.id == -1 {
			switch v := line.contents.(type) {
			case Transclusion:
				html += Transclude(v, recursionLevel)
			case Img:
				html += v.ToHtml()
			case Table:
				html += v.asHtml()
			case *List:
				html += v.RenderAsHtml()
			case string:
				html += v
			default:
				html += "<b class='error'>Unknown element.</b>"
			}
		}
	}
	return html
}
