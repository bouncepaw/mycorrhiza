module github.com/bouncepaw/mycorrhiza

go 1.16

require (
	github.com/bouncepaw/mycomarkup/v2 v2.0.0
	github.com/go-ini/ini v1.62.0
	github.com/gorilla/feeds v1.1.1
	github.com/gorilla/mux v1.8.0
	github.com/kr/pretty v0.2.1 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/valyala/quicktemplate v1.6.3
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	gopkg.in/ini.v1 v1.62.0 // indirect
)

// Use this trick to test mycomarkup:
replace github.com/bouncepaw/mycomarkup/v2 v2.0.0 => "/Users/bouncepaw/GolandProjects/mycomarkup"
