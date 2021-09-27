// Code generated by qtc from "stuff.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/stuff.qtpl:1
package views

//line views/stuff.qtpl:1
import "path/filepath"

//line views/stuff.qtpl:2
import "github.com/bouncepaw/mycorrhiza/cfg"

//line views/stuff.qtpl:3
import "github.com/bouncepaw/mycorrhiza/hyphae"

//line views/stuff.qtpl:4
import "github.com/bouncepaw/mycorrhiza/user"

//line views/stuff.qtpl:5
import "github.com/bouncepaw/mycorrhiza/util"

//line views/stuff.qtpl:6
import "github.com/bouncepaw/mycorrhiza/l18n"

//line views/stuff.qtpl:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/stuff.qtpl:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/stuff.qtpl:8
func StreamBaseHTML(qw422016 *qt422016.Writer, title, body string, lc *l18n.Localizer, u *user.User, headElements ...string) {
//line views/stuff.qtpl:8
	qw422016.N().S(`
<!doctype html>
<html lang="`)
//line views/stuff.qtpl:10
	qw422016.E().S(lc.Locale)
//line views/stuff.qtpl:10
	qw422016.N().S(`">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>`)
//line views/stuff.qtpl:14
	qw422016.E().S(title)
//line views/stuff.qtpl:14
	qw422016.N().S(`</title>
		<link rel="shortcut icon" href="/static/favicon.ico">
		<link rel="stylesheet" href="/static/style.css">
		<script src="/static/shortcuts.js"></script>
		`)
//line views/stuff.qtpl:18
	for _, el := range headElements {
//line views/stuff.qtpl:18
		qw422016.N().S(el)
//line views/stuff.qtpl:18
	}
//line views/stuff.qtpl:18
	qw422016.N().S(`
	</head>
	<body>
		<header>
			<nav class="main-width top-bar">
				<ul class="top-bar__wrapper">
					<li class="top-bar__section top-bar__section_home">
						<div class="top-bar__home-link-wrapper">
							<a class="top-bar__home-link" href="/">`)
//line views/stuff.qtpl:26
	qw422016.E().S(cfg.WikiName)
//line views/stuff.qtpl:26
	qw422016.N().S(`</a>
						</div>
					</li>
					<li class="top-bar__section top-bar__section_search">
						<form class="top-bar__search" method="GET" action="/title-search">
							<input type="text" name="q" placeholder="`)
//line views/stuff.qtpl:31
	qw422016.E().S(lc.Get("ui.title_search"))
//line views/stuff.qtpl:31
	qw422016.N().S(`" class="top-bar__search-bar">
						</form>
					</li>
					<li class="top-bar__section top-bar__section_auth">
					`)
//line views/stuff.qtpl:35
	if cfg.UseAuth {
//line views/stuff.qtpl:35
		qw422016.N().S(`
						<ul class="top-bar__auth auth-links">
							<li class="auth-links__box auth-links__user-box">
								`)
//line views/stuff.qtpl:38
		if u.Group == "anon" {
//line views/stuff.qtpl:38
			qw422016.N().S(`
								<a href="/login" class="auth-links__link auth-links__login-link">`)
//line views/stuff.qtpl:39
			qw422016.E().S(lc.Get("ui.login"))
//line views/stuff.qtpl:39
			qw422016.N().S(`</a>
								`)
//line views/stuff.qtpl:40
		} else {
//line views/stuff.qtpl:40
			qw422016.N().S(`
								<a href="/hypha/`)
//line views/stuff.qtpl:41
			qw422016.E().S(cfg.UserHypha)
//line views/stuff.qtpl:41
			qw422016.N().S(`/`)
//line views/stuff.qtpl:41
			qw422016.E().S(u.Name)
//line views/stuff.qtpl:41
			qw422016.N().S(`" class="auth-links__link auth-links__user-link">`)
//line views/stuff.qtpl:41
			qw422016.E().S(util.BeautifulName(u.Name))
//line views/stuff.qtpl:41
			qw422016.N().S(`</a>
								`)
//line views/stuff.qtpl:42
		}
//line views/stuff.qtpl:42
		qw422016.N().S(`
							</li>
							`)
//line views/stuff.qtpl:44
		if cfg.AllowRegistration && u.Group == "anon" {
//line views/stuff.qtpl:44
			qw422016.N().S(`
							<li class="auth-links__box auth-links__register-box">
								<a href="/register" class="auth-links__link auth-links__register-link">`)
//line views/stuff.qtpl:46
			qw422016.E().S(lc.Get("ui.register"))
//line views/stuff.qtpl:46
			qw422016.N().S(`</a>
							</li>
							`)
//line views/stuff.qtpl:48
		}
//line views/stuff.qtpl:48
		qw422016.N().S(`
							`)
//line views/stuff.qtpl:49
		if u.Group == "admin" {
//line views/stuff.qtpl:49
			qw422016.N().S(`
							<li class="auth-links__box auth-links__admin-box">
								<a href="/admin" class="auth-links__link auth-links__admin-link">`)
//line views/stuff.qtpl:51
			qw422016.E().S(lc.Get("ui.admin_panel"))
//line views/stuff.qtpl:51
			qw422016.N().S(`</a>
							</li>
							`)
//line views/stuff.qtpl:53
		}
//line views/stuff.qtpl:53
		qw422016.N().S(`
						</ul>
					`)
//line views/stuff.qtpl:55
	}
//line views/stuff.qtpl:55
	qw422016.N().S(`
					</li>
					<li class="top-bar__section top-bar__section_highlights">
						<ul class="top-bar__highlights">
`)
//line views/stuff.qtpl:59
	for _, link := range cfg.HeaderLinks {
//line views/stuff.qtpl:59
		qw422016.N().S(`						`)
//line views/stuff.qtpl:60
		if link.Href != "/" {
//line views/stuff.qtpl:60
			qw422016.N().S(`
							<li class="top-bar__highlight">
								<a class="top-bar__highlight-link" href="`)
//line views/stuff.qtpl:62
			qw422016.E().S(link.Href)
//line views/stuff.qtpl:62
			qw422016.N().S(`">`)
//line views/stuff.qtpl:62
			qw422016.E().S(link.Display)
//line views/stuff.qtpl:62
			qw422016.N().S(`</a>
							</li>
						`)
//line views/stuff.qtpl:64
		}
//line views/stuff.qtpl:64
		qw422016.N().S(`
`)
//line views/stuff.qtpl:65
	}
//line views/stuff.qtpl:65
	qw422016.N().S(`						</ul>
					</li>
				</ul>
			</nav>
		</header>
		`)
//line views/stuff.qtpl:71
	qw422016.N().S(body)
//line views/stuff.qtpl:71
	qw422016.N().S(`
		<template id="dialog-template">
			<div class="dialog-backdrop"></div>
			<div class="dialog" tabindex="0">
				<div class="dialog__header">
					<h1 class="dialog__title"></h1>
					<button class="dialog__close-button" aria-label="`)
//line views/stuff.qtpl:77
	qw422016.E().S(lc.Get("ui.close_dialog"))
//line views/stuff.qtpl:77
	qw422016.N().S(`"></button>
				</div>

				<div class="dialog__content"></div>
			</div>
		</template>
		`)
//line views/stuff.qtpl:83
	StreamCommonScripts(qw422016)
//line views/stuff.qtpl:83
	qw422016.N().S(`
		<script src="/static/view.js"></script>
	</body>
</html>
`)
//line views/stuff.qtpl:87
}

//line views/stuff.qtpl:87
func WriteBaseHTML(qq422016 qtio422016.Writer, title, body string, lc *l18n.Localizer, u *user.User, headElements ...string) {
//line views/stuff.qtpl:87
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:87
	StreamBaseHTML(qw422016, title, body, lc, u, headElements...)
//line views/stuff.qtpl:87
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:87
}

//line views/stuff.qtpl:87
func BaseHTML(title, body string, lc *l18n.Localizer, u *user.User, headElements ...string) string {
//line views/stuff.qtpl:87
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:87
	WriteBaseHTML(qb422016, title, body, lc, u, headElements...)
//line views/stuff.qtpl:87
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:87
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:87
	return qs422016
//line views/stuff.qtpl:87
}

//line views/stuff.qtpl:89
func StreamTitleSearchHTML(qw422016 *qt422016.Writer, query string, generator func(string) <-chan string, lc *l18n.Localizer) {
//line views/stuff.qtpl:89
	qw422016.N().S(`
<div class="layout">
<main class="main-width title-search">
	<h1>`)
//line views/stuff.qtpl:92
	qw422016.E().S(lc.Get("ui.search_results_query", &l18n.Replacements{"query": query}))
//line views/stuff.qtpl:92
	qw422016.N().S(`</h1>
	<p>`)
//line views/stuff.qtpl:93
	qw422016.E().S(lc.Get("ui.search_results_desc"))
//line views/stuff.qtpl:93
	qw422016.N().S(`</p>
	<ul class="title-search__results">
	`)
//line views/stuff.qtpl:95
	for hyphaName := range generator(query) {
//line views/stuff.qtpl:95
		qw422016.N().S(`
		<li class="title-search__entry">
			<a class="title-search__link wikilink" href="/hypha/`)
//line views/stuff.qtpl:97
		qw422016.E().S(hyphaName)
//line views/stuff.qtpl:97
		qw422016.N().S(`">`)
//line views/stuff.qtpl:97
		qw422016.E().S(util.BeautifulName(hyphaName))
//line views/stuff.qtpl:97
		qw422016.N().S(`</a>
		</li>
	`)
//line views/stuff.qtpl:99
	}
//line views/stuff.qtpl:99
	qw422016.N().S(`
</main>
</div>
`)
//line views/stuff.qtpl:102
}

//line views/stuff.qtpl:102
func WriteTitleSearchHTML(qq422016 qtio422016.Writer, query string, generator func(string) <-chan string, lc *l18n.Localizer) {
//line views/stuff.qtpl:102
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:102
	StreamTitleSearchHTML(qw422016, query, generator, lc)
//line views/stuff.qtpl:102
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:102
}

//line views/stuff.qtpl:102
func TitleSearchHTML(query string, generator func(string) <-chan string, lc *l18n.Localizer) string {
//line views/stuff.qtpl:102
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:102
	WriteTitleSearchHTML(qb422016, query, generator, lc)
//line views/stuff.qtpl:102
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:102
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:102
	return qs422016
//line views/stuff.qtpl:102
}

// It outputs a poorly formatted JSON, but it works and is valid.

//line views/stuff.qtpl:105
func StreamTitleSearchJSON(qw422016 *qt422016.Writer, query string, generator func(string) <-chan string) {
//line views/stuff.qtpl:105
	qw422016.N().S(`
`)
//line views/stuff.qtpl:107
	// Lol
	counter := 0

//line views/stuff.qtpl:109
	qw422016.N().S(`
{
	"source_query": "`)
//line views/stuff.qtpl:111
	qw422016.E().S(query)
//line views/stuff.qtpl:111
	qw422016.N().S(`",
	"results": [
	`)
//line views/stuff.qtpl:113
	for hyphaName := range generator(query) {
//line views/stuff.qtpl:113
		qw422016.N().S(`
		`)
//line views/stuff.qtpl:114
		if counter > 0 {
//line views/stuff.qtpl:114
			qw422016.N().S(`, `)
//line views/stuff.qtpl:114
		}
//line views/stuff.qtpl:114
		qw422016.N().S(`{
			"canonical_name": "`)
//line views/stuff.qtpl:115
		qw422016.E().S(hyphaName)
//line views/stuff.qtpl:115
		qw422016.N().S(`",
			"beautiful_name": "`)
//line views/stuff.qtpl:116
		qw422016.E().S(util.BeautifulName(hyphaName))
//line views/stuff.qtpl:116
		qw422016.N().S(`",
			"url": "`)
//line views/stuff.qtpl:117
		qw422016.E().S(cfg.URL + "/hypha/" + hyphaName)
//line views/stuff.qtpl:117
		qw422016.N().S(`"
		}`)
//line views/stuff.qtpl:118
		counter++

//line views/stuff.qtpl:118
		qw422016.N().S(`
	`)
//line views/stuff.qtpl:119
	}
//line views/stuff.qtpl:119
	qw422016.N().S(`
	]
}
`)
//line views/stuff.qtpl:122
}

//line views/stuff.qtpl:122
func WriteTitleSearchJSON(qq422016 qtio422016.Writer, query string, generator func(string) <-chan string) {
//line views/stuff.qtpl:122
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:122
	StreamTitleSearchJSON(qw422016, query, generator)
//line views/stuff.qtpl:122
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:122
}

//line views/stuff.qtpl:122
func TitleSearchJSON(query string, generator func(string) <-chan string) string {
//line views/stuff.qtpl:122
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:122
	WriteTitleSearchJSON(qb422016, query, generator)
//line views/stuff.qtpl:122
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:122
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:122
	return qs422016
//line views/stuff.qtpl:122
}

//line views/stuff.qtpl:124
func StreamBacklinksHTML(qw422016 *qt422016.Writer, query string, generator func(string) <-chan string, lc *l18n.Localizer) {
//line views/stuff.qtpl:124
	qw422016.N().S(`
<div class="layout">
<main class="main-width backlinks">
	<h1>`)
//line views/stuff.qtpl:127
	qw422016.E().S(lc.Get("ui.backlinks_query", &l18n.Replacements{"query": query}))
//line views/stuff.qtpl:127
	qw422016.N().S(`</h1>
	<p>`)
//line views/stuff.qtpl:128
	qw422016.E().S(lc.Get("ui.backlinks_desc"))
//line views/stuff.qtpl:128
	qw422016.N().S(`</p>
	<ul class="backlinks__list">
	`)
//line views/stuff.qtpl:130
	for hyphaName := range generator(query) {
//line views/stuff.qtpl:130
		qw422016.N().S(`
		<li class="backlinks__entry">
			<a class="backlinks__link wikilink" href="/hypha/`)
//line views/stuff.qtpl:132
		qw422016.E().S(hyphaName)
//line views/stuff.qtpl:132
		qw422016.N().S(`">`)
//line views/stuff.qtpl:132
		qw422016.E().S(util.BeautifulName(hyphaName))
//line views/stuff.qtpl:132
		qw422016.N().S(`</a>
		</li>
	`)
//line views/stuff.qtpl:134
	}
//line views/stuff.qtpl:134
	qw422016.N().S(`
</main>
</div>
`)
//line views/stuff.qtpl:137
}

//line views/stuff.qtpl:137
func WriteBacklinksHTML(qq422016 qtio422016.Writer, query string, generator func(string) <-chan string, lc *l18n.Localizer) {
//line views/stuff.qtpl:137
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:137
	StreamBacklinksHTML(qw422016, query, generator, lc)
//line views/stuff.qtpl:137
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:137
}

//line views/stuff.qtpl:137
func BacklinksHTML(query string, generator func(string) <-chan string, lc *l18n.Localizer) string {
//line views/stuff.qtpl:137
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:137
	WriteBacklinksHTML(qb422016, query, generator, lc)
//line views/stuff.qtpl:137
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:137
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:137
	return qs422016
//line views/stuff.qtpl:137
}

//line views/stuff.qtpl:139
func StreamHelpHTML(qw422016 *qt422016.Writer, content, lang string, lc *l18n.Localizer) {
//line views/stuff.qtpl:139
	qw422016.N().S(`
<div class="layout">
<main class="main-width help">
	<article>
	`)
//line views/stuff.qtpl:143
	qw422016.N().S(content)
//line views/stuff.qtpl:143
	qw422016.N().S(`
	</article>
</main>
`)
//line views/stuff.qtpl:146
	qw422016.N().S(helpTopicsHTML(lang, lc))
//line views/stuff.qtpl:146
	qw422016.N().S(`
</div>
`)
//line views/stuff.qtpl:148
}

//line views/stuff.qtpl:148
func WriteHelpHTML(qq422016 qtio422016.Writer, content, lang string, lc *l18n.Localizer) {
//line views/stuff.qtpl:148
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:148
	StreamHelpHTML(qw422016, content, lang, lc)
//line views/stuff.qtpl:148
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:148
}

//line views/stuff.qtpl:148
func HelpHTML(content, lang string, lc *l18n.Localizer) string {
//line views/stuff.qtpl:148
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:148
	WriteHelpHTML(qb422016, content, lang, lc)
//line views/stuff.qtpl:148
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:148
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:148
	return qs422016
//line views/stuff.qtpl:148
}

//line views/stuff.qtpl:150
func StreamHelpEmptyErrorHTML(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:150
	qw422016.N().S(`
<h1>`)
//line views/stuff.qtpl:151
	qw422016.E().S(lc.Get("help.empty_error_title"))
//line views/stuff.qtpl:151
	qw422016.N().S(`</h1>
<p>`)
//line views/stuff.qtpl:152
	qw422016.E().S(lc.Get("help.empty_error_line_1"))
//line views/stuff.qtpl:152
	qw422016.N().S(`</p>
<p>`)
//line views/stuff.qtpl:153
	qw422016.E().S(lc.Get("help.empty_error_line_2a"))
//line views/stuff.qtpl:153
	qw422016.N().S(` <a class="wikilink wikilink_external wikilink_https" href="https://github.com/bouncepaw/mycorrhiza">`)
//line views/stuff.qtpl:153
	qw422016.E().S(lc.Get("help.empty_error_link"))
//line views/stuff.qtpl:153
	qw422016.N().S(`</a> `)
//line views/stuff.qtpl:153
	qw422016.E().S(lc.Get("help.empty_error_line_2b"))
//line views/stuff.qtpl:153
	qw422016.N().S(`</p>
`)
//line views/stuff.qtpl:154
}

//line views/stuff.qtpl:154
func WriteHelpEmptyErrorHTML(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:154
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:154
	StreamHelpEmptyErrorHTML(qw422016, lc)
//line views/stuff.qtpl:154
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:154
}

//line views/stuff.qtpl:154
func HelpEmptyErrorHTML(lc *l18n.Localizer) string {
//line views/stuff.qtpl:154
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:154
	WriteHelpEmptyErrorHTML(qb422016, lc)
//line views/stuff.qtpl:154
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:154
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:154
	return qs422016
//line views/stuff.qtpl:154
}

//line views/stuff.qtpl:156
func streamhelpTopicsHTML(qw422016 *qt422016.Writer, lang string, lc *l18n.Localizer) {
//line views/stuff.qtpl:156
	qw422016.N().S(`
<aside class="help-topics layout-card">
	<h2 class="layout-card__title">`)
//line views/stuff.qtpl:158
	qw422016.E().S(lc.GetWithLocale(lang, "help.topics"))
//line views/stuff.qtpl:158
	qw422016.N().S(`</h2>
	<ul class="help-topics__list">
		<li><a href="/help/`)
//line views/stuff.qtpl:160
	qw422016.E().S(lang)
//line views/stuff.qtpl:160
	qw422016.N().S(`">`)
//line views/stuff.qtpl:160
	qw422016.E().S(lc.GetWithLocale(lang, "help.main"))
//line views/stuff.qtpl:160
	qw422016.N().S(`</a></li>
		<li><a href="/help/`)
//line views/stuff.qtpl:161
	qw422016.E().S(lang)
//line views/stuff.qtpl:161
	qw422016.N().S(`/hypha">`)
//line views/stuff.qtpl:161
	qw422016.E().S(lc.GetWithLocale(lang, "help.hypha"))
//line views/stuff.qtpl:161
	qw422016.N().S(`</a>
			<ul>
				<li><a href="/help/`)
//line views/stuff.qtpl:163
	qw422016.E().S(lang)
//line views/stuff.qtpl:163
	qw422016.N().S(`/attachment">`)
//line views/stuff.qtpl:163
	qw422016.E().S(lc.GetWithLocale(lang, "help.attachment"))
//line views/stuff.qtpl:163
	qw422016.N().S(`</a></li>
			</ul>
		</li>
		<li><a href="/help/`)
//line views/stuff.qtpl:166
	qw422016.E().S(lang)
//line views/stuff.qtpl:166
	qw422016.N().S(`/mycomarkup">`)
//line views/stuff.qtpl:166
	qw422016.E().S(lc.GetWithLocale(lang, "help.mycomarkup"))
//line views/stuff.qtpl:166
	qw422016.N().S(`</a></li>
		<li>`)
//line views/stuff.qtpl:167
	qw422016.E().S(lc.GetWithLocale(lang, "help.interface"))
//line views/stuff.qtpl:167
	qw422016.N().S(`
			<ul>
				<li><a href="/help/`)
//line views/stuff.qtpl:169
	qw422016.E().S(lang)
//line views/stuff.qtpl:169
	qw422016.N().S(`/prevnext">`)
//line views/stuff.qtpl:169
	qw422016.E().S(lc.GetWithLocale(lang, "help.prevnext"))
//line views/stuff.qtpl:169
	qw422016.N().S(`</a></li>
				<li><a href="/help/`)
//line views/stuff.qtpl:170
	qw422016.E().S(lang)
//line views/stuff.qtpl:170
	qw422016.N().S(`/top_bar">`)
//line views/stuff.qtpl:170
	qw422016.E().S(lc.GetWithLocale(lang, "help.top_bar"))
//line views/stuff.qtpl:170
	qw422016.N().S(`</a></li>
				<li><a href="/help/`)
//line views/stuff.qtpl:171
	qw422016.E().S(lang)
//line views/stuff.qtpl:171
	qw422016.N().S(`/sibling_hyphae_section">`)
//line views/stuff.qtpl:171
	qw422016.E().S(lc.GetWithLocale(lang, "help.sibling_hyphae"))
//line views/stuff.qtpl:171
	qw422016.N().S(`</a></li>
				<li>...</li>
			</ul>
		</li>
		<li>`)
//line views/stuff.qtpl:175
	qw422016.E().S(lc.GetWithLocale(lang, "help.special_pages"))
//line views/stuff.qtpl:175
	qw422016.N().S(`
			<ul>
				<li><a href="/help/`)
//line views/stuff.qtpl:177
	qw422016.E().S(lang)
//line views/stuff.qtpl:177
	qw422016.N().S(`/recent_changes">`)
//line views/stuff.qtpl:177
	qw422016.E().S(lc.GetWithLocale(lang, "help.recent_changes"))
//line views/stuff.qtpl:177
	qw422016.N().S(`</a></li>
			</ul>
		</li>
		<li>`)
//line views/stuff.qtpl:180
	qw422016.E().S(lc.GetWithLocale(lang, "help.configuration"))
//line views/stuff.qtpl:180
	qw422016.N().S(`
			<ul>
				<li><a href="/help/`)
//line views/stuff.qtpl:182
	qw422016.E().S(lang)
//line views/stuff.qtpl:182
	qw422016.N().S(`/lock">`)
//line views/stuff.qtpl:182
	qw422016.E().S(lc.GetWithLocale(lang, "help.lock"))
//line views/stuff.qtpl:182
	qw422016.N().S(`</a></li>
				<li><a href="/help/`)
//line views/stuff.qtpl:183
	qw422016.E().S(lang)
//line views/stuff.qtpl:183
	qw422016.N().S(`/whitelist">`)
//line views/stuff.qtpl:183
	qw422016.E().S(lc.GetWithLocale(lang, "help.whitelist"))
//line views/stuff.qtpl:183
	qw422016.N().S(`</a></li>
				<li><a href="/help/`)
//line views/stuff.qtpl:184
	qw422016.E().S(lang)
//line views/stuff.qtpl:184
	qw422016.N().S(`/telegram">`)
//line views/stuff.qtpl:184
	qw422016.E().S(lc.GetWithLocale(lang, "help.telegram"))
//line views/stuff.qtpl:184
	qw422016.N().S(`</a></li>
				<li>...</li>
			</ul>
		</li>
	</ul>
</aside>
`)
//line views/stuff.qtpl:190
}

//line views/stuff.qtpl:190
func writehelpTopicsHTML(qq422016 qtio422016.Writer, lang string, lc *l18n.Localizer) {
//line views/stuff.qtpl:190
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:190
	streamhelpTopicsHTML(qw422016, lang, lc)
//line views/stuff.qtpl:190
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:190
}

//line views/stuff.qtpl:190
func helpTopicsHTML(lang string, lc *l18n.Localizer) string {
//line views/stuff.qtpl:190
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:190
	writehelpTopicsHTML(qb422016, lang, lc)
//line views/stuff.qtpl:190
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:190
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:190
	return qs422016
//line views/stuff.qtpl:190
}

//line views/stuff.qtpl:192
func streamhelpTopicBadgeHTML(qw422016 *qt422016.Writer, lang, topic string) {
//line views/stuff.qtpl:192
	qw422016.N().S(`
<a class="help-topic-badge" href="/help/`)
//line views/stuff.qtpl:193
	qw422016.E().S(lang)
//line views/stuff.qtpl:193
	qw422016.N().S(`/`)
//line views/stuff.qtpl:193
	qw422016.E().S(topic)
//line views/stuff.qtpl:193
	qw422016.N().S(`">?</a>
`)
//line views/stuff.qtpl:194
}

//line views/stuff.qtpl:194
func writehelpTopicBadgeHTML(qq422016 qtio422016.Writer, lang, topic string) {
//line views/stuff.qtpl:194
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:194
	streamhelpTopicBadgeHTML(qw422016, lang, topic)
//line views/stuff.qtpl:194
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:194
}

//line views/stuff.qtpl:194
func helpTopicBadgeHTML(lang, topic string) string {
//line views/stuff.qtpl:194
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:194
	writehelpTopicBadgeHTML(qb422016, lang, topic)
//line views/stuff.qtpl:194
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:194
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:194
	return qs422016
//line views/stuff.qtpl:194
}

//line views/stuff.qtpl:196
func StreamUserListHTML(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:196
	qw422016.N().S(`
<div class="layout">
<main class="main-width user-list">
	<h1>`)
//line views/stuff.qtpl:199
	qw422016.E().S(lc.Get("ui.users_heading"))
//line views/stuff.qtpl:199
	qw422016.N().S(`</h1>
`)
//line views/stuff.qtpl:201
	var (
		admins     = make([]string, 0)
		moderators = make([]string, 0)
		editors    = make([]string, 0)
	)
	for u := range user.YieldUsers() {
		switch u.Group {
		case "admin":
			admins = append(admins, u.Name)
		case "moderator":
			moderators = append(moderators, u.Name)
		case "editor", "trusted":
			editors = append(editors, u.Name)
		}
	}

//line views/stuff.qtpl:216
	qw422016.N().S(`
	<section>
		<h2>`)
//line views/stuff.qtpl:218
	qw422016.E().S(lc.Get("ui.users_admins"))
//line views/stuff.qtpl:218
	qw422016.N().S(`</h2>
		<ol>`)
//line views/stuff.qtpl:219
	for _, name := range admins {
//line views/stuff.qtpl:219
		qw422016.N().S(`
			<li><a href="/hypha/`)
//line views/stuff.qtpl:220
		qw422016.E().S(cfg.UserHypha)
//line views/stuff.qtpl:220
		qw422016.N().S(`/`)
//line views/stuff.qtpl:220
		qw422016.E().S(name)
//line views/stuff.qtpl:220
		qw422016.N().S(`">`)
//line views/stuff.qtpl:220
		qw422016.E().S(name)
//line views/stuff.qtpl:220
		qw422016.N().S(`</a></li>
		`)
//line views/stuff.qtpl:221
	}
//line views/stuff.qtpl:221
	qw422016.N().S(`</ol>
	</section>
	<section>
		<h2>`)
//line views/stuff.qtpl:224
	qw422016.E().S(lc.Get("ui.users_moderators"))
//line views/stuff.qtpl:224
	qw422016.N().S(`</h2>
		<ol>`)
//line views/stuff.qtpl:225
	for _, name := range moderators {
//line views/stuff.qtpl:225
		qw422016.N().S(`
			<li><a href="/hypha/`)
//line views/stuff.qtpl:226
		qw422016.E().S(cfg.UserHypha)
//line views/stuff.qtpl:226
		qw422016.N().S(`/`)
//line views/stuff.qtpl:226
		qw422016.E().S(name)
//line views/stuff.qtpl:226
		qw422016.N().S(`">`)
//line views/stuff.qtpl:226
		qw422016.E().S(name)
//line views/stuff.qtpl:226
		qw422016.N().S(`</a></li>
		`)
//line views/stuff.qtpl:227
	}
//line views/stuff.qtpl:227
	qw422016.N().S(`</ol>
	</section>
	<section>
		<h2>`)
//line views/stuff.qtpl:230
	qw422016.E().S(lc.Get("ui.users_editors"))
//line views/stuff.qtpl:230
	qw422016.N().S(`</h2>
		<ol>`)
//line views/stuff.qtpl:231
	for _, name := range editors {
//line views/stuff.qtpl:231
		qw422016.N().S(`
			<li><a href="/hypha/`)
//line views/stuff.qtpl:232
		qw422016.E().S(cfg.UserHypha)
//line views/stuff.qtpl:232
		qw422016.N().S(`/`)
//line views/stuff.qtpl:232
		qw422016.E().S(name)
//line views/stuff.qtpl:232
		qw422016.N().S(`">`)
//line views/stuff.qtpl:232
		qw422016.E().S(name)
//line views/stuff.qtpl:232
		qw422016.N().S(`</a></li>
		`)
//line views/stuff.qtpl:233
	}
//line views/stuff.qtpl:233
	qw422016.N().S(`</ol>
	</section>
</main>
</div>
`)
//line views/stuff.qtpl:237
}

//line views/stuff.qtpl:237
func WriteUserListHTML(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:237
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:237
	StreamUserListHTML(qw422016, lc)
//line views/stuff.qtpl:237
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:237
}

//line views/stuff.qtpl:237
func UserListHTML(lc *l18n.Localizer) string {
//line views/stuff.qtpl:237
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:237
	WriteUserListHTML(qb422016, lc)
//line views/stuff.qtpl:237
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:237
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:237
	return qs422016
//line views/stuff.qtpl:237
}

//line views/stuff.qtpl:239
func StreamHyphaListHTML(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:239
	qw422016.N().S(`
<div class="layout">
<main class="main-width">
	<h1>`)
//line views/stuff.qtpl:242
	qw422016.E().S(lc.Get("ui.list_heading"))
//line views/stuff.qtpl:242
	qw422016.N().S(`</h1>
	<p>`)
//line views/stuff.qtpl:243
	qw422016.E().S(lc.GetPlural("ui.list_desc", hyphae.Count()))
//line views/stuff.qtpl:243
	qw422016.N().S(`</p>
	<ul class="hypha-list">
		`)
//line views/stuff.qtpl:245
	for h := range hyphae.YieldExistingHyphae() {
//line views/stuff.qtpl:245
		qw422016.N().S(`
		<li class="hypha-list__entry">
			<a class="hypha-list__link" href="/hypha/`)
//line views/stuff.qtpl:247
		qw422016.E().S(h.Name)
//line views/stuff.qtpl:247
		qw422016.N().S(`">`)
//line views/stuff.qtpl:247
		qw422016.E().S(util.BeautifulName(h.Name))
//line views/stuff.qtpl:247
		qw422016.N().S(`</a>
			`)
//line views/stuff.qtpl:248
		if h.BinaryPath != "" {
//line views/stuff.qtpl:248
			qw422016.N().S(`
			<span class="hypha-list__amnt-type">`)
//line views/stuff.qtpl:249
			qw422016.E().S(filepath.Ext(h.BinaryPath)[1:])
//line views/stuff.qtpl:249
			qw422016.N().S(`</span>
			`)
//line views/stuff.qtpl:250
		}
//line views/stuff.qtpl:250
		qw422016.N().S(`
		</li>
		`)
//line views/stuff.qtpl:252
	}
//line views/stuff.qtpl:252
	qw422016.N().S(`
	</ul>
</main>
</div>
`)
//line views/stuff.qtpl:256
}

//line views/stuff.qtpl:256
func WriteHyphaListHTML(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:256
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:256
	StreamHyphaListHTML(qw422016, lc)
//line views/stuff.qtpl:256
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:256
}

//line views/stuff.qtpl:256
func HyphaListHTML(lc *l18n.Localizer) string {
//line views/stuff.qtpl:256
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:256
	WriteHyphaListHTML(qb422016, lc)
//line views/stuff.qtpl:256
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:256
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:256
	return qs422016
//line views/stuff.qtpl:256
}

//line views/stuff.qtpl:258
func StreamAboutHTML(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:258
	qw422016.N().S(`
<div class="layout">
<main class="main-width">
	<section>
		<h1>`)
//line views/stuff.qtpl:262
	qw422016.E().S(lc.Get("ui.about_title", &l18n.Replacements{"name": cfg.WikiName}))
//line views/stuff.qtpl:262
	qw422016.N().S(`</h1>
		<ul>
			<li><b>`)
//line views/stuff.qtpl:264
	qw422016.N().S(lc.Get("ui.about_version", &l18n.Replacements{"pre": "<a href=\"https://mycorrhiza.wiki\">", "post": "</a>"}))
//line views/stuff.qtpl:264
	qw422016.N().S(`</b> 1.5.0</li>
`)
//line views/stuff.qtpl:265
	if cfg.UseAuth {
//line views/stuff.qtpl:265
		qw422016.N().S(`			<li><b>`)
//line views/stuff.qtpl:266
		qw422016.E().S(lc.Get("ui.about_usercount"))
//line views/stuff.qtpl:266
		qw422016.N().S(`</b> `)
//line views/stuff.qtpl:266
		qw422016.N().DUL(user.Count())
//line views/stuff.qtpl:266
		qw422016.N().S(`</li>
			<li><b>`)
//line views/stuff.qtpl:267
		qw422016.E().S(lc.Get("ui.about_homepage"))
//line views/stuff.qtpl:267
		qw422016.N().S(`</b> <a href="/">`)
//line views/stuff.qtpl:267
		qw422016.E().S(cfg.HomeHypha)
//line views/stuff.qtpl:267
		qw422016.N().S(`</a></li>
			<li><b>`)
//line views/stuff.qtpl:268
		qw422016.E().S(lc.Get("ui.about_admins"))
//line views/stuff.qtpl:268
		qw422016.N().S(`</b>`)
//line views/stuff.qtpl:268
		for i, username := range user.ListUsersWithGroup("admin") {
//line views/stuff.qtpl:269
			if i > 0 {
//line views/stuff.qtpl:269
				qw422016.N().S(`<span aria-hidden="true">, </span>
`)
//line views/stuff.qtpl:270
			}
//line views/stuff.qtpl:270
			qw422016.N().S(`				<a href="/hypha/`)
//line views/stuff.qtpl:271
			qw422016.E().S(cfg.UserHypha)
//line views/stuff.qtpl:271
			qw422016.N().S(`/`)
//line views/stuff.qtpl:271
			qw422016.E().S(username)
//line views/stuff.qtpl:271
			qw422016.N().S(`">`)
//line views/stuff.qtpl:271
			qw422016.E().S(username)
//line views/stuff.qtpl:271
			qw422016.N().S(`</a>`)
//line views/stuff.qtpl:271
		}
//line views/stuff.qtpl:271
		qw422016.N().S(`</li>
`)
//line views/stuff.qtpl:272
	} else {
//line views/stuff.qtpl:272
		qw422016.N().S(`			<li>`)
//line views/stuff.qtpl:273
		qw422016.E().S(lc.Get("ui.about_noauth"))
//line views/stuff.qtpl:273
		qw422016.N().S(`</li>
`)
//line views/stuff.qtpl:274
	}
//line views/stuff.qtpl:274
	qw422016.N().S(`		</ul>
		<p>`)
//line views/stuff.qtpl:276
	qw422016.N().S(lc.Get("ui.about_hyphae", &l18n.Replacements{"link": "<a href=\"/list\">/list</a>"}))
//line views/stuff.qtpl:276
	qw422016.N().S(`</p>
	</section>
</main>
</div>
`)
//line views/stuff.qtpl:280
}

//line views/stuff.qtpl:280
func WriteAboutHTML(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line views/stuff.qtpl:280
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:280
	StreamAboutHTML(qw422016, lc)
//line views/stuff.qtpl:280
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:280
}

//line views/stuff.qtpl:280
func AboutHTML(lc *l18n.Localizer) string {
//line views/stuff.qtpl:280
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:280
	WriteAboutHTML(qb422016, lc)
//line views/stuff.qtpl:280
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:280
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:280
	return qs422016
//line views/stuff.qtpl:280
}

//line views/stuff.qtpl:282
func StreamCommonScripts(qw422016 *qt422016.Writer) {
//line views/stuff.qtpl:282
	qw422016.N().S(`
`)
//line views/stuff.qtpl:283
	for _, scriptPath := range cfg.CommonScripts {
//line views/stuff.qtpl:283
		qw422016.N().S(`
<script src="`)
//line views/stuff.qtpl:284
		qw422016.E().S(scriptPath)
//line views/stuff.qtpl:284
		qw422016.N().S(`"></script>
`)
//line views/stuff.qtpl:285
	}
//line views/stuff.qtpl:285
	qw422016.N().S(`
`)
//line views/stuff.qtpl:286
}

//line views/stuff.qtpl:286
func WriteCommonScripts(qq422016 qtio422016.Writer) {
//line views/stuff.qtpl:286
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/stuff.qtpl:286
	StreamCommonScripts(qw422016)
//line views/stuff.qtpl:286
	qt422016.ReleaseWriter(qw422016)
//line views/stuff.qtpl:286
}

//line views/stuff.qtpl:286
func CommonScripts() string {
//line views/stuff.qtpl:286
	qb422016 := qt422016.AcquireByteBuffer()
//line views/stuff.qtpl:286
	WriteCommonScripts(qb422016)
//line views/stuff.qtpl:286
	qs422016 := string(qb422016.B)
//line views/stuff.qtpl:286
	qt422016.ReleaseByteBuffer(qb422016)
//line views/stuff.qtpl:286
	return qs422016
//line views/stuff.qtpl:286
}
