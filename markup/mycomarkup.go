// This is not done yet
package markup

import (
	"html"
	"strings"
)

// A Mycomarkup-formatted document
type MycoDoc struct {
	// data
	hyphaName string
	contents  string

	// state
	recursionDepth int

	// results
}

// Constructor
func Doc(hyphaName, contents string) *MycoDoc {
	return &MycoDoc{
		hyphaName: hyphaName,
		contents:  contents,
	}
}

// AsHtml returns an html representation of the document
func (md *MycoDoc) AsHtml() string {
	return ""
}

type BlockType int

const (
	BlockH1 = iota
	BlockH2
	BlockH3
	BlockH4
	BlockH5
	BlockH6
	BlockRocket
	BlockPre
	BlockQuote
	BlockPara
)

type CrawlWhere int

const (
	inSomewhere = iota
	inPre
	inEnd
)

func crawl(name, content string) []string {
	stateStack := []CrawlWhere{inSomewhere}

	startsWith := func(token string) bool {
		return strings.HasPrefix(content, token)
	}

	pop := func() {
		stateStack = stateStack[:len(stateStack)-1]
	}

	push := func(s CrawlWhere) {
		stateStack = append(stateStack, s)
	}

	readln := func(c string) (string, string) {
		parts := strings.SplitN(c, "\n", 1)
		return parts[0], parts[1]
	}

	preAcc := ""
	line := ""

	for {
		switch stateStack[0] {
		case inSomewhere:
			switch {
			case startsWith("```"):
				push(inPre)
				_, content = readln(content)
			default:
			}
		case inPre:
			switch {
			case startsWith("```"):
				pop()
				_, content = readln(content)
			default:
				line, content = readln(content)
				preAcc += html.EscapeString(line)
			}
		}
	}

	return []string{}
}
