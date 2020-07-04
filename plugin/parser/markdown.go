package parser

import (
	"github.com/bouncepaw/mycorrhiza/util"
	"gopkg.in/russross/blackfriday.v2"
)

func MarkdownToHtml(md []byte) string {
	return string(blackfriday.Run(util.NormalizeEOL(md)))
}
