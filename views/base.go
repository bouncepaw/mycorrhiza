package views

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/viewutil"
)

var (
	//go:embed *.html
	fs   embed.FS
	Base = viewutil.Base
)
