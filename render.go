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
        <div class="shroom">
            <button class="shroom__button" id="shroomBtn"><span>🍄</span> Open mycelium</button>
        </div>
		<main class="main">%s</main>
        <div class="left-panel" id="shroomburgerMenu">
            <div class="left-panel__in">
                <div class="shroom mushroom">
                    <button class="shroom__button" id="mushroomBtn"><span>🍄</span> Close mycelium</button>
                </div>
                <div class="left-panel__contents">
                    <header class="header">%s</header>
                    <aside class="sidebar">%s</aside>
                    <footer class="footer">%s</footer>
                </div>
            </div>
        </div>
		%s
	</body>
</html>
`
	return fmt.Sprintf(template, f["title"], f["head"], f["main"], f["header"], f["sidebar"], FooterText, f["bodyBottom"])
}

func EditHyphaPage(name, text_mime, content, tags string) string {
	template := `
<div class="naviwrapper">
	<form class="naviwrapper__edit edit-box"
	      method="POST"
	      enctype="multipart/form-data"
	      action="?action=update">
		<div class="naviwrapper__buttons">
			<input type="submit" value="update"/>
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

			<h4>Revision comment</h4>
			<p>Please make your comment helpful</p>
			<input type="text" name="comment" value="%s"/>

			<h4>Edit tags</h4>
			<p>Tags are separated by commas, whitespace is ignored</p>
			<input type="text" name="tags" value="%s"/>
		</div>
	</form>
</div>
`
	args := map[string]string{
		"title":   fmt.Sprintf(TitleTemplate, "Edit "+name),
		"head":    DefaultStyles,
		"header":  `<h1 class="header__edit-title">Edit ` + name + `</h1>`,
		"main":    fmt.Sprintf(template, content, text_mime, "Update "+name, tags),
		"sidebar": "",
		"footer":  FooterText,
	}

	return Layout(args)
}

func HyphaPage(hyphae map[string]*Hypha, rev Revision, content string) string {
	sidebar := `
<div class="naviwrapper">
    <div class="hypha-actions">
        <ul>
            <li><a href="?action=edit">Edit</a>
            <li><a href="?action=getBinary">Download</a>
            <li><a href="?action=zen">Zen mode</a>
            <li><a href="?action=raw">View raw</a>
        </ul>
    </div>
</div>
`

	bodyBottom := `
<script type="text/javascript">
    var menu = document.getElementById('shroomburgerMenu');
    document.getElementById('shroomBtn').addEventListener('click', function() {
        menu.classList.add('active');
    });
    document.getElementById('mushroomBtn').addEventListener('click', function() {
        menu.classList.remove('active');
    });
</script>
`

	args := map[string]string{
		"title":      fmt.Sprintf(TitleTemplate, rev.FullName),
		"head":       DefaultStyles,
		"header":     DefaultHeader,
		"main":       content,
		"sidebar":    sidebar,
		"footer":     FooterText,
		"bodyBottom": bodyBottom,
	}

	return Layout(args)
}
