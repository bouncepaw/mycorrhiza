package web

import (
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/shroom"
	"github.com/bouncepaw/mycorrhiza/mycoopts"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/util/fp"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

func handlerNewHypha(w http.ResponseWriter, rq *http.Request) {
	var (
		meta = viewutil.MetaFrom(w, rq)
		u    = meta.U
	)
	if !u.CanProceed("upload-text") || !u.CanProceed("upload-binary") {
		viewutil.HttpErr(meta, http.StatusForbidden, "", meta.Lc.Get("ui.act_norights_edit"))
		return
	}

	data := map[string]any{
		"IsUnnamed": true,
		"IsNew":     true,
		"HyphaName": rq.URL.Query().Get("name"),
	}
	render := func(status int, problem string) {
		data["Problem"] = problem
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		_ = pageHyphaEdit.RenderTo(meta, data)
	}

	if rq.Method == http.MethodGet {
		render(http.StatusOK, "")
		return
	}

	err := rq.ParseMultipartForm(multipartMemoryThreshold)
	if rq.MultipartForm != nil {
		defer rq.MultipartForm.RemoveAll()
	}
	var (
		name    = util.CanonicalName(rq.PostFormValue("name"))
		text    = strings.TrimSpace(rq.PostFormValue("text"))
		message = rq.PostFormValue("message")
		action  = rq.PostFormValue("action")
	)
	data["HyphaName"], data["Content"], data["Message"] = name, text, message
	if err != nil { // Far from call on purpose.
		render(http.StatusBadRequest, "upload")
		return
	}

	var media = fp.FirstOrEmpty(rq.MultipartForm.File["binary"])
	data["ReselectMedia"] = media != nil

	if action == "preview" {
		ctx, _ := mycocontext.ContextFromStringInput(text, mycoopts.MarkupOptions(name))
		data["Preview"] = template.HTML(mycomarkup.BlocksToHTML(ctx, mycomarkup.BlockTree(ctx)))
		render(http.StatusOK, "")
		return
	}

	var problem string
	switch {
	case !hyphae.IsValidName(name):
		problem = "name"
	case text == "" && media == nil:
		problem = "empty"
	case media != nil && media.Size == 0:
		problem = "upload"
	}
	if problem != "" {
		render(http.StatusBadRequest, problem)
		return
	}

	h := hyphae.ByName(name)
	if media != nil {
		file, err := media.Open()
		if err == nil {
			defer file.Close()
			err = shroom.UploadBinary(h, media.Header.Get("Content-Type"), file, u)
		}
		if err != nil {
			slog.Error("Failed to create media part of a new hypha", "hypha", name, "err", err)
			render(http.StatusInternalServerError, "upload")
			return
		}
		// UploadBinary may replace the hypha in the index.
		h = hyphae.ByName(name)
	}

	if text != "" {
		if err = shroom.UploadText(h, []byte(text), message, u); err != nil {
			slog.Error("Failed to create text part of a new hypha", "hypha", name, "err", err)
			problem = "save"
			if media != nil {
				data["ReselectMedia"] = false
				data["CreatedName"] = name
				problem = "partial"
			}
			render(http.StatusInternalServerError, problem)
			return
		}
	}

	destination := fp.If(action == "save-and-edit", "/edit/", "/hypha/") + name
	http.Redirect(w, rq, destination, http.StatusSeeOther)
}
