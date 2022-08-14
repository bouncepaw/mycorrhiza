package viewutil

import (
	"fmt"
	"mime"
	"net/http"
)

// HttpErr is used by many handlers to signal errors in a compact way.
// TODO: get rid of this abomination
func HttpErr(meta Meta, status int, name, errMsg string) {
	meta.W.(http.ResponseWriter).Header().Set("Content-Type", mime.TypeByExtension(".html"))
	meta.W.(http.ResponseWriter).WriteHeader(status)
	fmt.Fprint(
		meta.W,
		Base(
			meta,
			"Error",
			fmt.Sprintf(
				`<main class="main-width"><p>%s. <a href="/hypha/%s">%s<a></p></main>`,
				errMsg,
				name,
				meta.Lc.Get("ui.error_go_back"),
			),
			map[string]string{},
		),
	)
}
