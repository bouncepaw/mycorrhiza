package render

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/fs"
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
			"Title":   fmt.Sprintf(cfg.TitleEditTemplate, h.FullName),
			"Header":  layout("edit/header").withString(h.FullName).String(),
			"Sidebar": layout("edit/sidebar").withMap(hyphaData).String(),
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
	return hyphaGeneric(name, name, "view/404")
}

// HyphaPage renders hypha viewer.
func HyphaPage(name, content string) []byte {
	return hyphaGeneric(name, content, "view/index")
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
		"Main":      "",
		"SiteTitle": cfg.SiteTitle,
	}
	for key, val := range keys {
		page[key] = val
	}
	page["Main"] = lyt.String()
	return layout("base").withMap(page).Bytes()
}

func hyphaTree(name string) string {
	return layout("view/sidebar").
		withMap(map[string]string{"Tree": fs.Hs.GetTree(name, true).AsHtml()}).
		String()
}

type Layout struct {
	tmpl    *template.Template
	buf     *bytes.Buffer
	invalid bool
	err     error
}

func layout(name string) *Layout {
	h := fs.Hs.Open(path.Join(cfg.TemplatesDir, cfg.Theme, name+".html")).OnRevision("0")
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
