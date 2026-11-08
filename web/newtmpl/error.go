package newtmpl

import (
	"log/slog"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/web/viewdata"
)

var pageError *Page

func Init() {
	pageError = NewPage("error.html")
}

type errorData struct {
	BaseData
	Message     string
	ReturnURL   string
	ReturnLabel string
}

func Error(meta viewdata.Meta, status int, message, returnURL, returnLabel string) {
	w := meta.W.(http.ResponseWriter)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(status)
	if err := pageError.RenderDataTo(meta, &errorData{
		Message: message, ReturnURL: returnURL, ReturnLabel: returnLabel,
	}); err != nil {
		slog.Error("Failed to execute page", "page", "error", "err", err)
	}
}

func HttpErr(meta viewdata.Meta, status int, name, errMsg string) {
	Error(meta, status, errMsg+".", "/hypha/"+name, "Go back to the hypha.")
}

func HandlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Page not found"))
}
