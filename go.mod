module github.com/bouncepaw/mycorrhiza

go 1.21
toolchain go1.24.1

require (
	git.sr.ht/~bouncepaw/mycomarkup/v5 v5.6.0
	github.com/go-ini/ini v1.67.0
	github.com/gorilla/feeds v1.2.0
	github.com/gorilla/mux v1.8.1
	github.com/valyala/quicktemplate v1.7.0
	golang.org/x/crypto v0.35.0
	golang.org/x/term v0.29.0
	golang.org/x/text v0.22.0
)

require (
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

// Use this trick to test local Mycomarkup changes, replace the path with yours,
// but do not commit the change to the path:
// replace git.sr.ht/~bouncepaw/mycomarkup/v5 v5.6.0 => "/Users/bouncepaw/src/mycomarkup"
