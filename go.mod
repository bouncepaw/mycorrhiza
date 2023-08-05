module github.com/bouncepaw/mycorrhiza

go 1.19

require (
	git.sr.ht/~bouncepaw/mycomarkup/v5 v5.5.0
	github.com/go-ini/ini v1.63.2
	github.com/gorilla/feeds v1.1.1
	github.com/gorilla/mux v1.8.0
	github.com/valyala/quicktemplate v1.7.0
	golang.org/x/crypto v0.1.0
	golang.org/x/exp v0.0.0-20220414153411-bcd21879b8fd
	golang.org/x/term v0.1.0
	golang.org/x/text v0.4.0
)

require (
	github.com/kr/pretty v0.2.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

// Use this trick to test local Mycomarkup changes, replace the path with yours,
// but do not commit the change to the path:
// replace git.sr.ht/~bouncepaw/mycomarkup/v5 v5.5.0 => "/Users/bouncepaw/src/mycomarkup"
