package parser

import (
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/m4tty/cajun"
)

func CreoleToHtml(creole []byte) string {
	out, _ := cajun.Transform(string(util.NormalizeEOL(creole)))
	return out
}
