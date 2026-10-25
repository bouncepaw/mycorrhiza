module github.com/bouncepaw/mycorrhiza

go 1.27.3

require (
	git.sr.ht/~bouncepaw/mycomarkup/v5 v5.6.1
	github.com/gorilla/feeds v1.2.0
	github.com/gorilla/mux v1.8.1
	golang.org/x/crypto v0.58.0
	golang.org/x/term v0.47.0
	golang.org/x/text v0.43.0
	gopkg.in/ini.v1 v1.67.3
)

require golang.org/x/sys v0.49.0 // indirect

// Use this trick to test local Mycomarkup changes, replace the path with yours,
// but do not commit the change to the path:
// replace git.sr.ht/~bouncepaw/mycomarkup/v5 v5.6.1 => "/Users/bouncepaw/src/mycomarkup"
