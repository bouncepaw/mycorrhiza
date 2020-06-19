package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"
)

func EditHyphaPage(name, textMime, content, tags string) string {
	keys := map[string]string{
		"Title":  fmt.Sprintf(TitleTemplate, "Edit "+name),
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

func HyphaPage(hyphae map[string]*Hypha, rev Revision, content string) string {
	sidebar := DefaultSidebar
	bside, err := ioutil.ReadFile("Hypha/view/sidebar.html")
	if err == nil {
		sidebar = string(bside)
	}
	keys := map[string]string{
		"Title":   fmt.Sprintf(TitleTemplate, rev.FullName),
		"Sidebar": sidebar,
	}
	return renderBase(renderFromString(content, "Hypha/view/index.html"), keys)
}

/*
Collect and render page from base template
Args:
	content: string or pre-rendered template
	keys: map with replaced standart fields
*/
func renderBase(content string, keys map[string]string) string {
	page := map[string]string{
		"Title":      DefaultTitle,
		"Head":       DefaultStyles,
		"Sidebar":    DefaultSidebar,
		"Main":       DefaultContent,
		"BodyBottom": DefaultBodyBottom,
		"Header":     renderFromString(DefaultHeaderText, "header.html"),
		"Footer":     renderFromString(DefaultFooterText, "footer.html"),
	}
	for key, val := range keys {
		page[key] = val
	}
	page["Main"] = content
	return renderFromMap(page, "base.html")
}

func renderFromMap(data map[string]string, templatePath string) string {
	filePath := path.Join("templates", templatePath)
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

func renderFromString(data string, templatePath string) string {
	filePath := path.Join("templates", templatePath)
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
