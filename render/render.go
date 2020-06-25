package render

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/fs"
)

// EditHyphaPage returns HTML page of hypha editor.
func EditHyphaPage(name, textMime, content, tags string) string {
	page := map[string]string{
		"Text":     content,
		"TextMime": textMime,
		"Name":     name,
		"Tags":     tags,
	}
	keys := map[string]string{
		"Title":   fmt.Sprintf(cfg.TitleEditTemplate, name),
		"Header":  renderFromString(name, "Hypha/edit/header.html"),
		"Sidebar": renderFromMap(page, "Hypha/edit/sidebar.html"),
	}
	return renderBase(renderFromMap(page, "Hypha/edit/index.html"), keys)
}

// Hypha404 returns 404 page for nonexistent page.
func Hypha404(name, _ string) string {
	return HyphaGeneric(name, name, "Hypha/view/404.html")
}

// HyphaPage returns HTML page of hypha viewer.
func HyphaPage(name, content string) string {
	return HyphaGeneric(name, content, "Hypha/view/index.html")
}

// HyphaGeneric is used when building renderers for all types of hypha pages
func HyphaGeneric(name string, content, templatePath string) string {
	var sidebar string
	if aside := renderFromMap(map[string]string{
		"Tree": fs.Hs.GetTree(name, true).AsHtml(),
	}, "Hypha/view/sidebar.html"); aside != "" {
		sidebar = aside
	}
	keys := map[string]string{
		"Title":   fmt.Sprintf(cfg.TitleTemplate, name),
		"Sidebar": sidebar,
	}
	return renderBase(renderFromString(content, templatePath), keys)
}

// renderBase collects and renders page from base template
// Args:
//   content: string or pre-rendered template
//   keys: map with replaced standart fields
func renderBase(content string, keys map[string]string) string {
	page := map[string]string{
		"Title":     cfg.SiteTitle,
		"Main":      "",
		"SiteTitle": cfg.SiteTitle,
	}
	for key, val := range keys {
		page[key] = val
	}
	page["Main"] = content
	return renderFromMap(page, "base.html")
}

// renderFromMap applies `data` map to template in `templatePath` and returns the result.
func renderFromMap(data map[string]string, templatePath string) string {
	hyphPath := path.Join(cfg.TemplatesDir, cfg.Theme, templatePath)
	h, _ := fs.Hs.Open(hyphPath)
	h.OnRevision("0")
	tmpl, err := template.ParseFiles(h.TextPath())
	if err != nil {
		return err.Error()
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return err.Error()
	}
	return buf.String()
}

// renderFromMap applies `data` string to template in `templatePath` and returns the result.
func renderFromString(data string, templatePath string) string {
	hyphPath := path.Join(cfg.TemplatesDir, cfg.Theme, templatePath)
	h, _ := fs.Hs.Open(hyphPath)
	h.OnRevision("0")
	tmpl, err := template.ParseFiles(h.TextPath())
	if err != nil {
		return err.Error()
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return err.Error()
	}
	return buf.String()
}
