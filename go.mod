module github.com/bouncepaw/mycorrhiza

go 1.21

require (
	git.sr.ht/~bouncepaw/mycomarkup/v5 v5.6.0
	github.com/go-ini/ini v1.67.0
	github.com/gorilla/feeds v1.2.0
	github.com/gorilla/mux v1.8.1
	golang.org/x/crypto v0.31.0
	golang.org/x/term v0.27.0
	golang.org/x/text v0.21.0
)

require (
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

// Use this trick to test local Mycomarkup changes, replace the path with yours,
// but do not commit the change to the path:
// replace git.sr.ht/~bouncepaw/mycomarkup/v5 v5.6.0 => "/Users/bouncepaw/src/mycomarkup"
