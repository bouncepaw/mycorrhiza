package main

import (
	"bytes"
	"fmt"
	"path"
	"text/template"
)

func EditHyphaPage(name, text_mime, binary_mime, content, tags string) string {
	keys := map[string]string{
		"Title": fmt.Sprintf(TitleTemplate, name),
	}
	page := map[string]string{
		"Text":     content,
		"TextMime": text_mime,
		"BinMime":  binary_mime,
		"Name":     name,
		"Tags":     tags,
	}
	return renderBase(renderFromMap(page, "Hypha/edit.html"), keys)
}

func HyphaPage(hyphae map[string]*Hypha, rev Revision, content string) string {
	keys := map[string]string{
		"Title": fmt.Sprintf(TitleTemplate, rev.FullName),
	}
	return renderBase(renderFromString(content, "Hypha/index.html"), keys)
}

/*
Collect and render page from base template
Args:
	content: string or pre-rendered template
	keys: map with replaced standart fields
*/
func renderBase(content string, keys map[string]string) string {
	page := map[string]string{
		"Title":   DefaultTitle,
		"Header":  renderFromString(DefaultHeaderText, "header.html"),
		"Footer":  renderFromString(DefaultFooterText, "footer.html"),
		"Sidebar": DefaultSidebar,
		"Main":    DefaultContent,
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
