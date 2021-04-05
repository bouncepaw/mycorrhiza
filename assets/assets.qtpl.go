// Code generated by qtc from "assets.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line assets/assets.qtpl:1
package assets

//line assets/assets.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line assets/assets.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line assets/assets.qtpl:1
func StreamHelpMessage(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:1
	qw422016.N().S(`Usage of %s:
`)
//line assets/assets.qtpl:3
}

//line assets/assets.qtpl:3
func WriteHelpMessage(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:3
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:3
	StreamHelpMessage(qw422016)
//line assets/assets.qtpl:3
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:3
}

//line assets/assets.qtpl:3
func HelpMessage() string {
//line assets/assets.qtpl:3
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:3
	WriteHelpMessage(qb422016)
//line assets/assets.qtpl:3
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:3
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:3
	return qs422016
//line assets/assets.qtpl:3
}

//line assets/assets.qtpl:5
func StreamExampleConfig(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:6
	qw422016.N().S(`WikiName = My wiki
NaviTitleIcon = 🐑

[Hyphae]
HomeHypha = home
UserHypha = u
HeaderLinksHypha = header-links

[Network]
HTTPPort = 8080
URL = https://wiki
GeminiCertificatePath = /home/wiki/gemcerts

[Authorization]
UseFixedAuth = true
FixedAuthCredentialsPath = /home/wiki/mycocredentials.json
`)
//line assets/assets.qtpl:6
	qw422016.N().S(`
`)
//line assets/assets.qtpl:7
}

//line assets/assets.qtpl:7
func WriteExampleConfig(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:7
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:7
	StreamExampleConfig(qw422016)
//line assets/assets.qtpl:7
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:7
}

//line assets/assets.qtpl:7
func ExampleConfig() string {
//line assets/assets.qtpl:7
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:7
	WriteExampleConfig(qb422016)
//line assets/assets.qtpl:7
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:7
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:7
	return qs422016
//line assets/assets.qtpl:7
}

//line assets/assets.qtpl:9
func StreamDefaultCSS(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:9
	qw422016.N().S(`
`)
//line assets/assets.qtpl:10
	qw422016.N().S(`.amnt-grid { display: grid; grid-template-columns: 1fr 1fr; }
.upload-binary__input { display: block; margin: .25rem 0; }

.modal__title { font-size: 2rem; }
.modal__title_small { font-size: 1.5rem; }
.modal__confirmation-msg { margin: 0 0 .5rem 0; }
.modal__action { display: inline-block; font-size: 1rem; padding: .25rem; border-radius: .25rem; }
.modal__submit { border: 1px #999 solid; }
.modal__cancel { border: 1px #999 dashed; text-decoration: none; }

.hypha-list { padding-left: 0; }
.hypha-list__entry { list-style-type: none; }
.hypha-list__link { text-decoration: none; display: inline-block; padding: .25rem; }
.hypha-list__link:hover { text-decoration: underline; }
.hypha-list__amnt-type { font-size: smaller; color: #999; }

/* General element positions, from small to big */
/* Phones and whatnot */
.layout { display: grid; row-gap: 1rem; }
header { width: 100%; margin-bottom: 1rem; }
.header-links__list, .hypha-tabs__flex { margin: 0; padding: 0; display: flex; flex-wrap: wrap; }
.header-links__entry, .hypha-tabs__tab { list-style-type: none; }

.header-links__entry { margin-right: .5rem; }
.header-links__entry_user { font-style:italic; }
.header-links__link { display: inline-block; padding: .25rem; text-decoration: none; }

.hypha-tabs { padding: 0; margin: 0; }
.hypha-tabs__tab { margin-right: .5rem; padding: 0; }
.hypha-tabs__link { display: inline-block; padding: .25rem; text-decoration: none; }
.hypha-tabs__selection { display: inline-block; padding: .25rem; font-weight: bold; }

.layout-card li { list-style-type: none; }
.backlinks__list { padding: 0; margin: 0; }
.backlinks__link { text-decoration: none; display: block; padding: .25rem; padding-left: 1.25rem; }

@media screen and (max-width: 800px) {
	.amnt-grid { grid-template-columns: 1fr; }
	.layout { grid-template-columns: auto; grid-template-rows: auto auto auto; }
	.main-width { width: 100%; }
	main { padding: 1rem; margin: 0; }
}

/* No longer a phone but still small screen: draw normal tabs, center main */
@media screen and (min-width: 801px) {
	.main-width { padding: 1rem 2rem; width: 800px; margin: 0 auto; }
	main { border-radius: .25rem; }
	.layout-card { width: 800px; margin: 0 auto; }

	.header-links { padding: 0; }
	.header-links__entry { margin-right: 1.5rem; }
	.header-links__entry_user { margin: 0 2rem 0 auto; }
	.header-links__entry:nth-of-type(1),

	.hypha-tabs { padding: 0; }
	.hypha-tabs__tab { border-radius: .25rem .25rem 0 0; margin-right: 0; }
	.hypha-tabs__selection, .hypha-tabs__link { padding: .25rem .5rem; }

	.header-links__entry:nth-of-type(1), .hypha-tabs__tab:nth-of-type(1) { margin-left: 2rem; }
}

/* Wide enough to fit two columns ok */
@media screen and (min-width: 1100px) {
	.layout { display: grid; grid-template-columns: auto 1fr; column-gap: 1rem; margin: 0 1rem; row-gap: 1rem; }
	.main-width { margin: 0; }
	main { grid-column: 1 / span 1; grid-row: 1 / span 2; }
	.relative-hyphae { grid-column: 2 / span 1; grid-row: 1 / span 1; }
	.layout-card { width: 100%; }
	.edit-toolbar { margin: 0 0 0 auto; }
	.edit-toolbar__buttons { display: grid; grid-template-columns: 1fr 1fr; }
}

@media screen and (min-width: 1250px) {
	.layout { grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr); }
	.layout-card { max-width: 18rem; }
	.main-width { margin: 0 auto; }
	.backlinks { grid-column: 1 / span 1; margin-right: 0; }
	main { grid-column: 2 / span 1; }
	.relative-hyphae { grid-column: 3 / span 1; margin-left: 0; }

	.backlinks__title { text-align: right; }
	.backlinks__link { text-align: right; padding-right: 1.25rem; padding-left: .25rem; }
}

*, *::before, *::after {box-sizing: border-box;}
html { height:100%; padding:0; }
body {height:100%; margin:0; }
body, input { font-size:16px; font-family: 'PT Sans', 'Liberation Sans', sans-serif;}
main > form {margin-bottom:1rem;}
textarea {font-size:16px; font-family: 'PT Sans', 'Liberation Sans', sans-serif;}

.edit { min-height: 80vh; }
.edit__title { margin-top: 0; }
.edit__preview { border: 2px dashed #ddd; }
.edit-form {height:70vh;}
.edit-form textarea {width:100%;height:95%;}
.edit-form__save { font-weight: bold; }
.edit-toolbar__buttons, .edit-toolbar__ad { margin: .5rem; }

.icon {margin-right: .25rem; vertical-align: bottom; }

main h1:not(.navi-title) {font-size:1.7rem;}
blockquote { margin-left: 0; padding-left: 1rem; }
.wikilink_external::before { display: inline-block; width: 18px; height: 16px; vertical-align: sub; }
/* .wikilink_external { padding-left: 16px; } */
.wikilink_gopher::before { content: url("/static/icon/gopher"); }
.wikilink_http::before { content: url("/static/icon/http"); }
.wikilink_https::before { content: url("/static/icon/http"); }
/* .wikilink_https { background: transparent url("/static/icon/http") center left no-repeat; } */
.wikilink_gemini::before { content: url("/static/icon/gemini"); }
.wikilink_mailto::before { content: url("/static/icon/mailto"); }

article { overflow-wrap: break-word; word-wrap: break-word; word-break: break-word; line-height: 150%; }
main h1, main h2, main h3, main h4, main h5, main h6 { margin: 1.5rem 0 0 0; }
.heading__link { text-decoration: none; display: inline-block; }
.heading__link::after { width: 1rem; content: "§"; color: transparent; }
.heading__link:hover::after, .heading__link:active::after { color: #999; }
article p { margin: .5rem 0; }
article ul, ol { padding-left: 1.5rem; margin: .5rem 0; }
article code { padding: .1rem .3rem; border-radius: .25rem; font-size: 90%; }
article pre.codeblock { padding:.5rem; white-space: pre-wrap; border-radius: .25rem;}
.codeblock code {padding:0; font-size:15px;}
.transclusion { border-radius: .25rem; }
.transclusion__content > *:not(.binary-container) {margin: 0.5rem; }
.transclusion__link {display: block; text-align: right; font-style: italic; margin-top: .5rem; margin-right: .25rem; text-decoration: none;}
.transclusion__link::before {content: "⇐ ";}

/* Derived from https://commons.wikimedia.org/wiki/File:U%2B21D2.svg */
.launchpad__entry { list-style-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' version='1.0' width='25' height='12'%3E%3Cg transform='scale(0.7,0.8) translate(-613.21429,-421)'%3E%3Cpath fill='%23999' d='M 638.06773,429.49751 L 631.01022,436.87675 L 630.1898,436.02774 L 632.416,433.30375 L 613.46876,433.30375 L 613.46876,431.66382 L 633.82089,431.66382 L 635.57789,429.5261 L 633.79229,427.35979 L 613.46876,427.35979 L 613.46876,425.71985 L 632.416,425.71985 L 630.1898,422.99587 L 631.01022,422.08788 L 638.06773,429.49751 z '/%3E%3C/g%3E%3C/svg%3E"); }

.binary-container_with-img img,
.binary-container_with-video video,
.binary-container_with-audio audio {width: 100%}

.subhyphae__title { padding-bottom: .5rem; clear: both; }
.navi-title { padding-bottom: .5rem; margin: .25rem 0; }
.navi-title a {text-decoration:none; }
.navi-title__separator { margin: 0 .25rem; }
.navi-title__colon { margin-right: .5rem; }
.upload-amnt { clear: both; padding: .5rem; border-radius: .25rem; }
.upload-amnt__unattach { display: block; }
aside { clear: both; }

.img-gallery { text-align: center; margin-top: .25rem; margin-bottom: .25rem; }
.img-gallery_many-images { border-radius: .25rem; padding: .5rem; }
.img-gallery img { max-width: 100%; max-height: 50vh; }
figure { margin: 0; }
figcaption { padding-bottom: .5rem; }

#new-name {width:100%;}


.rc-entry { display: grid; list-style-type: none; padding: .25rem; grid-template-columns: 1fr 1fr; border-radius: .25rem; }
.rc-entry__time { font-style: italic; }
.rc-entry__hash { font-style: italic; text-align: right; }
.rc-entry__links, .rc-entry__msg { grid-column: 1 / span 2; }
.rc-entry__author { font-style: italic; }

.prevnext__el { display: inline-block; min-width: 40%; padding: .5rem; margin-bottom: .25rem; text-decoration: none; border-radius: .25rem; }
.prevnext__prev { float: left; }
.prevnext__next { float: right; text-align: right; }

.page-separator { clear: both; }
.history__entries { background-color: #eee; margin: 0; padding: 0; border-radius: .25rem; }
.history__month-anchor { text-decoration: none; color: inherit; }
.history__entry { list-style-type: none; padding: .25rem; }
.history-entry { padding: .25rem; }
.history-entry__time { font-weight: bold; }
.history-entry__author { font-style: italic; }

table { border: #ddd 1px solid; border-radius: .25rem; min-width: 4rem; }
td { padding: .25rem; }
caption { caption-side: top; font-size: small; }

.subhyphae__list, .subhyphae__list ul { display: flex; padding: 0; margin: 0; flex-wrap: wrap; }
.subhyphae__entry { list-style-type: none; border: 1px solid #999; padding: 0; margin: .125rem; border-radius: .25rem; }
.subhyphae__link { display: block; padding: .25rem; text-decoration: none; }
.subhyphae__link:hover { background: #eee; }

.navitree { padding: 0; margin: 0; }
.navitree__entry { }
.navitree > .navitree__entry > a::before { display: inline-block; width: .5rem; color: #999; margin: 0 .25rem; }
.navitree > .navitree__entry_infertile > a::before { content: " "} /* nbsp, careful */
.navitree > .navitree__sibling_fertile > a::before { content: "▸"}
.navitree__trunk { border-left: 1px #999 solid; }
.navitree__link { text-decoration: none; display: block; padding: .25rem; }
.navitree__entry_this > span { display: block; padding: .25rem; font-weight: bold; }
.navitree__entry_this > span::before { content: " "; display: inline-block; width: 1rem; }


/* Color stuff */
/* Lighter stuff #eee */
article code,
article .codeblock,
.transclusion,
.img-gallery_many-images,
.rc-entry,
.prevnext__el,
table { background-color: #eee; }

.hypha-tabs__tab { background-color: #eee; }
.hypha-tabs__tab a { color: black; }
.hypha-tabs__tab_active { border-bottom: 2px white solid; background: white; }

@media screen and (max-width: 800px) {
	.hypha-tabs,
	.hypha-tabs__tab { background-color: white; }
}

@media screen and (min-width: 801px) {
	.hypha-tabs__tab { border: 1px #ddd solid; }
	.hypha-tabs__tab_active { border-bottom: 1px white solid; }
}

.layout-card { border-radius: .25rem; background-color: white; }
.layout-card__title { font-size: 1rem; margin: 0; padding: .25rem .5rem; border-radius: .25rem .25rem 0 0; }
.layout-card__title { background-color: #eee; }

/* Other stuff */
html { background-color: #ddd; 
background-image: url("data:image/svg+xml,%3Csvg width='42' height='44' viewBox='0 0 42 44' xmlns='http://www.w3.org/2000/svg'%3E%3Cg id='Page-1' fill='none' fill-rule='evenodd'%3E%3Cg id='brick-wall' fill='%23bbbbbb' fill-opacity='0.4'%3E%3Cpath d='M0 0h42v44H0V0zm1 1h40v20H1V1zM0 23h20v20H0V23zm22 0h20v20H22V23z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
} /* heropatterns.com */
header { background-color: #bbb; }
.header-links__link { color: black; }
.header-links__link:hover { background-color: #eee; }
main { background-color: white; }

blockquote { border-left: 4px black solid; }
.wikilink_new {color:#a55858;}
.transclusion code, .transclusion .codeblock {background-color:#ddd;}
.transclusion__link { color: black; }
.wikilink_new:visited {color:#a55858;}
.navi-title { border-bottom: #eee 1px solid; }
.upload-amnt { border: #eee 1px solid; }
td { border: #ddd 1px solid; }

.navitree__link:hover, .backlinks__link:hover { background-color: #eee; }

/* Dark theme! */
@media (prefers-color-scheme: dark) {
html { background: #222; color: #ddd; }
main,  article, .hypha-tabs__tab, header, .layout-card { background-color: #343434; color: #ddd; }

a, .wikilink_external { color: #f1fa8c; }
a:visited, .wikilink_external:visited { color: #ffb86c; }
.wikilink_new, .wikilink_new:visited { color: #dd4444; }
.subhyphae__link:hover, .navitree__link:hover, .backlinks__link:hover { background-color: #444; }

.header-links__link, .header-links__link:visited,
.prevnext__el, .prevnext__el:visited { color: #ddd; }
.header-links__link:hover { background-color: #444; }

.hypha-tabs__tab a, .hypha-tabs__tab { color: #ddd; background-color: #232323; border: 0; }
.layout-card__title, .hypha-tabs__tab_active { background-color: #343434; }

blockquote { border-left: 4px #ddd solid; }

.transclusion .transclusion__link { color: #ddd; }
article code, 
article .codeblock, 
.transclusion,
.img-gallery_many-images,
.rc-entry,
.history__entry, 
.prevnext__el, 
.upload-amnt, 
textarea,
table { border: 0; background-color: #444444; color: #ddd; }
.transclusion code,
.transclusion .codeblock { background-color: #454545; }
mark { background: rgba(130, 80, 30, 5); color: inherit; }
@media screen and (max-width: 800px) {
	.hypha-tabs { background-color: #232323; }
}
}
`)
//line assets/assets.qtpl:10
	qw422016.N().S(`
`)
//line assets/assets.qtpl:11
}

//line assets/assets.qtpl:11
func WriteDefaultCSS(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:11
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:11
	StreamDefaultCSS(qw422016)
//line assets/assets.qtpl:11
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:11
}

//line assets/assets.qtpl:11
func DefaultCSS() string {
//line assets/assets.qtpl:11
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:11
	WriteDefaultCSS(qb422016)
//line assets/assets.qtpl:11
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:11
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:11
	return qs422016
//line assets/assets.qtpl:11
}

//line assets/assets.qtpl:13
func StreamToolbarJS(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:13
	qw422016.N().S(`
`)
//line assets/assets.qtpl:14
	qw422016.N().S(`const editTextarea = document.getElementsByClassName('edit-form__textarea')[0]

function placeCursor(position, el = editTextarea) {
    el.selectionEnd = position
    el.selectionStart = el.selectionEnd
}

function getSelectedText(el = editTextarea) {
    const [start, end] = [el.selectionStart, el.selectionEnd]
    const text = el.value
    return text.substring(start, end)
}

function textInserter(text, cursorPosition = null, el = editTextarea) {
    return function() {
        const [start, end] = [el.selectionStart, el.selectionEnd]
        el.setRangeText(text, start, end, 'select')
        el.focus()
        if (cursorPosition == null) {
            placeCursor(end + text.length)
        } else {
            placeCursor(end + cursorPosition)
        }
    }
}

function selectionWrapper(cursorPosition, prefix, postfix = null, el = editTextarea) {
    return function() {
        const [start, end] = [el.selectionStart, el.selectionEnd]
        if (postfix == null) {
            postfix = prefix
        }
        text = getSelectedText(el)
        result = prefix + text + postfix
        el.setRangeText(result, start, end, 'select')
        el.focus()
        placeCursor(end + cursorPosition)
    }
}

const wrapBold = selectionWrapper(2, '**'), 
    wrapItalic = selectionWrapper(2, '//'), 
    wrapMonospace = selectionWrapper(1, '`)
//line assets/assets.qtpl:14
	qw422016.N().S("`")
//line assets/assets.qtpl:14
	qw422016.N().S(`'), 
    wrapHighlighted = selectionWrapper(2, '!!'), 
    wrapLifted = selectionWrapper(1, '^'), 
    wrapLowered = selectionWrapper(2, ',,'), 
    wrapStrikethrough = selectionWrapper(2, '~~'), 
    wrapLink = selectionWrapper(2, '[[', ']]')

const insertHorizontalBar = textInserter('----\n'),
    insertImgBlock = textInserter('img {\n\t\n}\n', 7), 
    insertTableBlock = textInserter('table {\n\t\n}\n', 9),
    insertRocket = textInserter('=> '),
    insertXcl = textInserter('<= ')

function insertDate() {
    let date = new Date().toISOString().split('T')[0]
    textInserter(date)()
}

function insertUserlink() {
    const userlink = document.querySelector('.header-links__entry_user a')
    const userHypha = userlink.getAttribute('href').substring(7) // no /hypha/
    textInserter('[[' + userHypha + ']]')()
}
`)
//line assets/assets.qtpl:14
	qw422016.N().S(`
`)
//line assets/assets.qtpl:15
}

//line assets/assets.qtpl:15
func WriteToolbarJS(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:15
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:15
	StreamToolbarJS(qw422016)
//line assets/assets.qtpl:15
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:15
}

//line assets/assets.qtpl:15
func ToolbarJS() string {
//line assets/assets.qtpl:15
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:15
	WriteToolbarJS(qb422016)
//line assets/assets.qtpl:15
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:15
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:15
	return qs422016
//line assets/assets.qtpl:15
}

// Next three are from https://remixicon.com/

//line assets/assets.qtpl:18
func StreamIconHTTP(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:18
	qw422016.N().S(`
`)
//line assets/assets.qtpl:19
	qw422016.N().S(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16"><path fill="#999" d="M12 22C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm-2.29-2.333A17.9 17.9 0 0 1 8.027 13H4.062a8.008 8.008 0 0 0 5.648 6.667zM10.03 13c.151 2.439.848 4.73 1.97 6.752A15.905 15.905 0 0 0 13.97 13h-3.94zm9.908 0h-3.965a17.9 17.9 0 0 1-1.683 6.667A8.008 8.008 0 0 0 19.938 13zM4.062 11h3.965A17.9 17.9 0 0 1 9.71 4.333 8.008 8.008 0 0 0 4.062 11zm5.969 0h3.938A15.905 15.905 0 0 0 12 4.248 15.905 15.905 0 0 0 10.03 11zm4.259-6.667A17.9 17.9 0 0 1 15.973 11h3.965a8.008 8.008 0 0 0-5.648-6.667z"/></svg>
`)
//line assets/assets.qtpl:19
	qw422016.N().S(`
`)
//line assets/assets.qtpl:20
}

//line assets/assets.qtpl:20
func WriteIconHTTP(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:20
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:20
	StreamIconHTTP(qw422016)
//line assets/assets.qtpl:20
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:20
}

//line assets/assets.qtpl:20
func IconHTTP() string {
//line assets/assets.qtpl:20
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:20
	WriteIconHTTP(qb422016)
//line assets/assets.qtpl:20
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:20
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:20
	return qs422016
//line assets/assets.qtpl:20
}

//line assets/assets.qtpl:22
func StreamIconGemini(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:22
	qw422016.N().S(`
`)
//line assets/assets.qtpl:23
	qw422016.N().S(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16"><path fill="#999" d="M15.502 20A6.523 6.523 0 0 1 12 23.502 6.523 6.523 0 0 1 8.498 20h2.26c.326.489.747.912 1.242 1.243.495-.33.916-.754 1.243-1.243h2.259zM18 14.805l2 2.268V19H4v-1.927l2-2.268V9c0-3.483 2.504-6.447 6-7.545C15.496 2.553 18 5.517 18 9v5.805zM17.27 17L16 15.56V9c0-2.318-1.57-4.43-4-5.42C9.57 4.57 8 6.681 8 9v6.56L6.73 17h10.54zM12 11a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"/></svg>
`)
//line assets/assets.qtpl:23
	qw422016.N().S(`
`)
//line assets/assets.qtpl:24
}

//line assets/assets.qtpl:24
func WriteIconGemini(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:24
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:24
	StreamIconGemini(qw422016)
//line assets/assets.qtpl:24
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:24
}

//line assets/assets.qtpl:24
func IconGemini() string {
//line assets/assets.qtpl:24
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:24
	WriteIconGemini(qb422016)
//line assets/assets.qtpl:24
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:24
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:24
	return qs422016
//line assets/assets.qtpl:24
}

//line assets/assets.qtpl:26
func StreamIconMailto(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:26
	qw422016.N().S(`
`)
//line assets/assets.qtpl:27
	qw422016.N().S(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16"><path fill="#999" d="M3 3h18a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zm17 4.238l-7.928 7.1L4 7.216V19h16V7.238zM4.511 5l7.55 6.662L19.502 5H4.511z"/></svg>
`)
//line assets/assets.qtpl:27
	qw422016.N().S(`
`)
//line assets/assets.qtpl:28
}

//line assets/assets.qtpl:28
func WriteIconMailto(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:28
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:28
	StreamIconMailto(qw422016)
//line assets/assets.qtpl:28
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:28
}

//line assets/assets.qtpl:28
func IconMailto() string {
//line assets/assets.qtpl:28
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:28
	WriteIconMailto(qb422016)
//line assets/assets.qtpl:28
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:28
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:28
	return qs422016
//line assets/assets.qtpl:28
}

// This is a modified version of https://www.svgrepo.com/svg/232085/rat

//line assets/assets.qtpl:31
func StreamIconGopher(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:31
	qw422016.N().S(`
`)
//line assets/assets.qtpl:32
	qw422016.N().S(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="16" height="16">
<path fill="#999" d="M447.238,204.944v-70.459c0-8.836-7.164-16-16-16c-34.051,0-64.414,21.118-75.079,55.286
C226.094,41.594,0,133.882,0,319.435c0,0.071,0.01,0.14,0.011,0.21c0.116,44.591,36.423,80.833,81.04,80.833h171.203
c8.836,0,16-7.164,16-16c0-8.836-7.164-16-16-16H81.051c-21.441,0-39.7-13.836-46.351-33.044H496c8.836,0,16-7.164,16-16
C512,271.82,486.82,228.692,447.238,204.944z M415.238,153.216v37.805c-10.318-2.946-19.556-4.305-29.342-4.937
C390.355,168.611,402.006,157.881,415.238,153.216z M295.484,303.435L295.484,303.435c-7.562-41.495-43.948-73.062-87.593-73.062
c-8.836,0-16,7.164-16,16c0,8.836,7.164,16,16,16c25.909,0,47.826,17.364,54.76,41.062H32.722
c14.415-159.15,218.064-217.856,315.136-90.512c3.545,4.649,9.345,6.995,15.124,6.118
c55.425-8.382,107.014,29.269,115.759,84.394H295.484z"/>
<circle fill="#999" cx="415.238" cy="260.05" r="21.166"/>
</svg>
`)
//line assets/assets.qtpl:32
	qw422016.N().S(`
`)
//line assets/assets.qtpl:33
}

//line assets/assets.qtpl:33
func WriteIconGopher(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:33
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:33
	StreamIconGopher(qw422016)
//line assets/assets.qtpl:33
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:33
}

//line assets/assets.qtpl:33
func IconGopher() string {
//line assets/assets.qtpl:33
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:33
	WriteIconGopher(qb422016)
//line assets/assets.qtpl:33
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:33
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:33
	return qs422016
//line assets/assets.qtpl:33
}

// https://upload.wikimedia.org/wikipedia/commons/4/46/Generic_Feed-icon.svg

//line assets/assets.qtpl:36
func StreamIconFeed(qw422016 *qt422016.Writer) {
//line assets/assets.qtpl:36
	qw422016.N().S(`
`)
//line assets/assets.qtpl:37
	qw422016.N().S(`<svg xmlns="http://www.w3.org/2000/svg"
     id="RSSicon"
     viewBox="0 0 8 8" width="256" height="256">

  <title>RSS feed icon</title>

  <style type="text/css">
    .button {stroke: none; fill: orange;}
    .symbol {stroke: none; fill: white;}
  </style>

  <rect   class="button" width="8" height="8" rx="1.5" />
  <circle class="symbol" cx="2" cy="6" r="1" />
  <path   class="symbol" d="m 1,4 a 3,3 0 0 1 3,3 h 1 a 4,4 0 0 0 -4,-4 z" />
  <path   class="symbol" d="m 1,2 a 5,5 0 0 1 5,5 h 1 a 6,6 0 0 0 -6,-6 z" />

</svg>
`)
//line assets/assets.qtpl:37
	qw422016.N().S(`
`)
//line assets/assets.qtpl:38
}

//line assets/assets.qtpl:38
func WriteIconFeed(qq422016 qtio422016.Writer) {
//line assets/assets.qtpl:38
	qw422016 := qt422016.AcquireWriter(qq422016)
//line assets/assets.qtpl:38
	StreamIconFeed(qw422016)
//line assets/assets.qtpl:38
	qt422016.ReleaseWriter(qw422016)
//line assets/assets.qtpl:38
}

//line assets/assets.qtpl:38
func IconFeed() string {
//line assets/assets.qtpl:38
	qb422016 := qt422016.AcquireByteBuffer()
//line assets/assets.qtpl:38
	WriteIconFeed(qb422016)
//line assets/assets.qtpl:38
	qs422016 := string(qb422016.B)
//line assets/assets.qtpl:38
	qt422016.ReleaseByteBuffer(qb422016)
//line assets/assets.qtpl:38
	return qs422016
//line assets/assets.qtpl:38
}
