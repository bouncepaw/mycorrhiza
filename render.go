package main

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

// EditHyphaPage returns HTML page of hypha editor.
func EditHyphaPage(name, textMime, content, tags string) string {
	keys := map[string]string{
		"Title":  fmt.Sprintf(cfg.TitleEditTemplate, name),
		"Header": renderFromString(name, "Hypha/edit/header.html"),
	}
	page := map[string]string{
		"Text":     content,
		"TextMime": textMime,
		"Name":     name,
		"Tags":     tags,
	}
	return renderBase(renderFromMap(page, "Hypha/edit/index.html"), keys)
}

// Hypha404 returns 404 page for nonexistent page.
func Hypha404(name string) string {
	return HyphaGeneric(name, name, "Hypha/view/404.html")
}

// HyphaPage returns HTML page of hypha viewer.
func HyphaPage(rev Revision, content string) string {
	return HyphaGeneric(rev.FullName, content, "Hypha/view/index.html")
}

// HyphaGeneric is used when building renderers for all types of hypha pages
func HyphaGeneric(name string, content, templatePath string) string {
	sidebar := cfg.DefaultSidebar
	if bside := renderFromMap(map[string]string{
		"Tree": GetTree(name, true).AsHtml(),
	}, "Hypha/view/sidebar.html"); bside != "" {
		sidebar = bside
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
		"Title":      cfg.DefaultTitle,
		"Head":       cfg.DefaultStyles,
		"Sidebar":    cfg.DefaultSidebar,
		"Main":       cfg.DefaultContent,
		"BodyBottom": cfg.DefaultBodyBottom,
		"Header":     renderFromString(cfg.DefaultHeaderText, "header.html"),
		"Footer":     renderFromString(cfg.DefaultFooterText, "footer.html"),
	}
	for key, val := range keys {
		page[key] = val
	}
	page["Main"] = content
	return renderFromMap(page, "base.html")
}

// renderFromMap applies `data` map to template in `templatePath` and returns the result.
func renderFromMap(data map[string]string, templatePath string) string {
	filePath := path.Join(cfg.TemplatesDir, templatePath)
	tmpl, err := template.ParseFiles(filePath)
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
	filePath := path.Join(cfg.TemplatesDir, templatePath)
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return err.Error()
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return err.Error()
	}
	return buf.String()
}
