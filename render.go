package main

import (
	"fmt"
)

func Layout(f map[string]string) string {
	template := `
<!doctype html>
<html>
	<head>
		<title>%s</title>
		%s
	</head>
	<body>
		<header class="header">%s</header>
		<main class="main">%s</main>
		<aside class="sidebar">%s</aside>
		<footer class="footer">%s</footer>
		%s
	</body>
</html>
`
	return fmt.Sprintf(template, f["title"], f["head"], f["header"], f["main"], f["sidebar"], FooterText, f["bodyBottom"])
}

func EditHyphaPage(name, text_mime, binary_mime, content, tags string) string {
	template := `
<div class="naviwrapper">
	<form class="naviwrapper__edit edit-box">
		<div class="naviwrapper__buttons">
			<input type="submit" name="action" value="update"/>
		</div>

		<div class="edit-box__left">
			<h4>Edit box</h4>
			<textarea class="edit-box__text" name="text" cols="80" rows="25">
%s
			</textarea>

			<h4>Upload file</h4>
			<p>If this hypha has a file like that, the text above is meant to be a description of it</p>
			<input type="file" name="binary"/>
		</div>

		<div class="edit-box__right">
			<h4>Text MIME-type</h4>
			<p>Good types are <code>text/markdown</code> and <code>text/plain</code></p>
			<input type="text" name="text_mime" value="%s"/>

			<h4>Media MIME-type</h4>
			<p>For now, only image formats are supported. Choose any, but <code>image/jpeg</code> and <code>image/png</code> are recommended</p>
			<input type="text" name="binary_mime" value="%s"/>

			<h4>Revision comment</h4>
			<p>Please make your comment helpful</p>
			<input type="text" name="comment" value="%s"/>

			<h4>Edit tags</h4>
			<p>Tags are separated by commas, whitespace is ignored</p>
			<input type="text" name="comment" value="%s"/>
		</div>
	</form>
</div>
`
	args := map[string]string{
		"title":   fmt.Sprintf(TitleTemplate, "Edit "+name),
		"head":    DefaultStyles,
		"header":  `<h1 class="header__edit-title">Edit ` + name + `</h1>`,
		"main":    fmt.Sprintf(template, content, text_mime, binary_mime, "Update "+name, tags),
		"sidebar": "",
		"footer":  FooterText,
	}

	return Layout(args)
}

func HyphaPage(hyphae map[string]*Hypha, rev Revision, content string) string {
	template := `
<div class="naviwrapper">
	<form class="naviwrapper__buttons">
		<input type="submit" name="action" value="edit"/>
		<input type="submit" name="action" value="getBinary"/>
		<input type="submit" name="action" value="zen"/>
		<input type="submit" name="action" value="raw"/>
	</form>
	%s
</div>
`
	args := map[string]string{
		"title":   fmt.Sprintf(TitleTemplate, rev.FullName),
		"head":    DefaultStyles,
		"header":  DefaultHeader,
		"main":    fmt.Sprintf(template, content),
		"sidebar": "",
		"footer":  FooterText,
	}

	return Layout(args)
}
