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
