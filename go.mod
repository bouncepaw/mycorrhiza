module github.com/bouncepaw/mycorrhiza

go 1.14

require (
	github.com/gorilla/mux v1.7.4
	mvdan.cc/gogrep v0.0.0-20200420132841-24e8804e5b3c // indirect
)

require (
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
