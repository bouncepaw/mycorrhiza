package render

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/fs"
	"github.com/bouncepaw/mycorrhiza/mycelium"
)

// HyphaEdit renders hypha editor.
func HyphaEdit(h *fs.Hypha) []byte { //
	hyphaData := map[string]string{
		"Name":     h.FullName,
		"Tags":     h.TagsJoined(),
		"TextMime": h.TextMime(),
		"Text":     h.TextContent(),
	}
	return layout("edit/index").
		withMap(hyphaData).
		wrapInBase(map[string]string{
			"Title": fmt.Sprintf(cfg.TitleEditTemplate, h.FullName),
		})
}

// HyphaUpdateOk is used to inform that update was successful.
func HyphaUpdateOk(h *fs.Hypha) []byte { //
	return layout("update_ok").
		withMap(map[string]string{"Name": h.FullName}).
		Bytes()
}

// Hypha404 renders 404 page for nonexistent page.
func Hypha404(name, _ string) []byte {
	return layout("view/404").
		withMap(map[string]string{
			"PageTitle": name,
			"Tree":      hyphaTree(name),
		}).
		wrapInBase(map[string]string{
			"Title": fmt.Sprintf(cfg.TitleTemplate, name),
		})
}

// HyphaPage renders hypha viewer.
func HyphaPage(name, content string) []byte {
	return layout("view/index").
		withMap(map[string]string{
			"Content": content,
			"Tree":    hyphaTree(name),
		}).
		wrapInBase(map[string]string{
			"Title": fmt.Sprintf(cfg.TitleTemplate, name),
		})
}

// hyphaGeneric is used when building renderers for all types of hypha pages
func hyphaGeneric(name, content, templateName string) []byte {
	return layout(templateName).
		withString(content).
		wrapInBase(map[string]string{
			"Title":   fmt.Sprintf(cfg.TitleTemplate, name),
			"Sidebar": hyphaTree(name),
		})
}

// wrapInBase is used to wrap layouts in things that are present on all pages.
func (lyt *Layout) wrapInBase(keys map[string]string) []byte {
	if lyt.invalid {
		return lyt.Bytes()
	}
	page := map[string]string{
		"Title":     cfg.SiteTitle,
		"Content":   lyt.String(),
		"SiteTitle": cfg.SiteTitle,
	}
	for key, val := range keys {
		page[key] = val
	}
	return layout("base").withMap(page).Bytes()
}

func hyphaTree(name string) string {
	return fs.Hs.GetTree(name, true).AsHtml()
}

type Layout struct {
	tmpl    *template.Template
	buf     *bytes.Buffer
	invalid bool
	err     error
}

func layout(name string) *Layout {
	lytName := path.Join("theme", cfg.Theme, name+".html")
	h := fs.Hs.OpenFromMap(map[string]string{
		"mycelium": mycelium.SystemMycelium,
		"hypha":    lytName,
		"rev":      "0",
	})
	if h.Invalid {
		return &Layout{nil, nil, true, h.Err}
	}
	tmpl, err := template.ParseFiles(h.TextPath())
	if err != nil {
		return &Layout{nil, nil, true, err}
	}
	return &Layout{tmpl, new(bytes.Buffer), false, nil}
}

func (lyt *Layout) withString(data string) *Layout {
	if lyt.invalid {
		return lyt
	}
	if err := lyt.tmpl.Execute(lyt.buf, data); err != nil {
		lyt.invalid = true
		lyt.err = err
	}
	return lyt
}

func (lyt *Layout) withMap(data map[string]string) *Layout {
	if lyt.invalid {
		return lyt
	}
	if err := lyt.tmpl.Execute(lyt.buf, data); err != nil {
		lyt.invalid = true
		lyt.err = err
	}
	return lyt
}

func (lyt *Layout) Bytes() []byte {
	if lyt.invalid {
		return []byte(lyt.err.Error())
	}
	return lyt.buf.Bytes()
}

func (lyt *Layout) String() string {
	if lyt.invalid {
		return lyt.err.Error()
	}
	return lyt.buf.String()
}
