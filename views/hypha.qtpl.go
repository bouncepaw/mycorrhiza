// Code generated by qtc from "hypha.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/hypha.qtpl:1
package views

//line views/hypha.qtpl:1
import "path/filepath"

//line views/hypha.qtpl:2
import "strings"

//line views/hypha.qtpl:4
import "github.com/bouncepaw/mycorrhiza/cfg"

//line views/hypha.qtpl:5
import "github.com/bouncepaw/mycorrhiza/hyphae"

//line views/hypha.qtpl:6
import "github.com/bouncepaw/mycorrhiza/l18n"

//line views/hypha.qtpl:7
import "github.com/bouncepaw/mycorrhiza/user"

//line views/hypha.qtpl:8
import "github.com/bouncepaw/mycorrhiza/util"

//line views/hypha.qtpl:10
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/hypha.qtpl:10
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/hypha.qtpl:10
func streambeautifulLink(qw422016 *qt422016.Writer, hyphaName string) {
//line views/hypha.qtpl:10
	qw422016.N().S(`<a href="/hypha/`)
//line views/hypha.qtpl:10
	qw422016.N().S(hyphaName)
//line views/hypha.qtpl:10
	qw422016.N().S(`">`)
//line views/hypha.qtpl:10
	qw422016.E().S(util.BeautifulName(hyphaName))
//line views/hypha.qtpl:10
	qw422016.N().S(`</a>`)
//line views/hypha.qtpl:10
}

//line views/hypha.qtpl:10
func writebeautifulLink(qq422016 qtio422016.Writer, hyphaName string) {
//line views/hypha.qtpl:10
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/hypha.qtpl:10
	streambeautifulLink(qw422016, hyphaName)
//line views/hypha.qtpl:10
	qt422016.ReleaseWriter(qw422016)
//line views/hypha.qtpl:10
}

//line views/hypha.qtpl:10
func beautifulLink(hyphaName string) string {
//line views/hypha.qtpl:10
	qb422016 := qt422016.AcquireByteBuffer()
//line views/hypha.qtpl:10
	writebeautifulLink(qb422016, hyphaName)
//line views/hypha.qtpl:10
	qs422016 := string(qb422016.B)
//line views/hypha.qtpl:10
	qt422016.ReleaseByteBuffer(qb422016)
//line views/hypha.qtpl:10
	return qs422016
//line views/hypha.qtpl:10
}

//line views/hypha.qtpl:12
func streammycoLink(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line views/hypha.qtpl:12
	qw422016.N().S(`<a href="/help/en/mycomarkup" class="shy-link">`)
//line views/hypha.qtpl:12
	qw422016.E().S(lc.Get("ui.notexist_write_myco"))
//line views/hypha.qtpl:12
	qw422016.N().S(`</a>`)
//line views/hypha.qtpl:12
}

//line views/hypha.qtpl:12
func writemycoLink(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line views/hypha.qtpl:12
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/hypha.qtpl:12
	streammycoLink(qw422016, lc)
//line views/hypha.qtpl:12
	qt422016.ReleaseWriter(qw422016)
//line views/hypha.qtpl:12
}

//line views/hypha.qtpl:12
func mycoLink(lc *l18n.Localizer) string {
//line views/hypha.qtpl:12
	qb422016 := qt422016.AcquireByteBuffer()
//line views/hypha.qtpl:12
	writemycoLink(qb422016, lc)
//line views/hypha.qtpl:12
	qs422016 := string(qb422016.B)
//line views/hypha.qtpl:12
	qt422016.ReleaseByteBuffer(qb422016)
//line views/hypha.qtpl:12
	return qs422016
//line views/hypha.qtpl:12
}

//line views/hypha.qtpl:14
func streamnonExistentHyphaNotice(qw422016 *qt422016.Writer, h hyphae.Hypha, u *user.User, lc *l18n.Localizer) {
//line views/hypha.qtpl:14
	qw422016.N().S(`
<section class="non-existent-hypha">
	<h2 class="non-existent-hypha__title">`)
//line views/hypha.qtpl:16
	qw422016.E().S(lc.Get("ui.notexist_heading"))
//line views/hypha.qtpl:16
	qw422016.N().S(`</h2>
	`)
//line views/hypha.qtpl:17
	if cfg.UseAuth && u.Group == "anon" {
//line views/hypha.qtpl:17
		qw422016.N().S(`
	<p>`)
//line views/hypha.qtpl:18
		qw422016.E().S(lc.Get("ui.notexist_norights"))
//line views/hypha.qtpl:18
		qw422016.N().S(`</p>
	<ul>
		<li><a href="/login">`)
//line views/hypha.qtpl:20
		qw422016.E().S(lc.Get("ui.notexist_login"))
//line views/hypha.qtpl:20
		qw422016.N().S(`</a></li>
		`)
//line views/hypha.qtpl:21
		if cfg.AllowRegistration {
//line views/hypha.qtpl:21
			qw422016.N().S(`<li><a href="/register">`)
//line views/hypha.qtpl:21
			qw422016.E().S(lc.Get("ui.notexist_register"))
//line views/hypha.qtpl:21
			qw422016.N().S(`</a></li>`)
//line views/hypha.qtpl:21
		}
//line views/hypha.qtpl:21
		qw422016.N().S(`
	</ul>
	`)
//line views/hypha.qtpl:23
	} else {
//line views/hypha.qtpl:23
		qw422016.N().S(`

	<div class="non-existent-hypha__ways">
	<section class="non-existent-hypha__way">
		<h3 class="non-existent-hypha__subtitle">📝 `)
//line views/hypha.qtpl:27
		qw422016.E().S(lc.Get("ui.notexist_write"))
//line views/hypha.qtpl:27
		qw422016.N().S(`</h3>
		<p>`)
//line views/hypha.qtpl:28
		qw422016.N().S(lc.Get("ui.notexist_write_tip1", &l18n.Replacements{"myco": mycoLink(lc)}))
//line views/hypha.qtpl:28
		qw422016.N().S(`</p>
		<p>`)
//line views/hypha.qtpl:29
		qw422016.E().S(lc.Get("ui.notexist_write_tip2"))
//line views/hypha.qtpl:29
		qw422016.N().S(`</p>
		<a class="btn btn_accent stick-to-bottom" href="/edit/`)
//line views/hypha.qtpl:30
		qw422016.E().S(h.CanonicalName())
//line views/hypha.qtpl:30
		qw422016.N().S(`">`)
//line views/hypha.qtpl:30
		qw422016.E().S(lc.Get("ui.notexist_write_button"))
//line views/hypha.qtpl:30
		qw422016.N().S(`</a>
	</section>

	<section class="non-existent-hypha__way">
		<h3 class="non-existent-hypha__subtitle">🖼 `)
//line views/hypha.qtpl:34
		qw422016.E().S(lc.Get("ui.notexist_media"))
//line views/hypha.qtpl:34
		qw422016.N().S(`</h3>
		<p>`)
//line views/hypha.qtpl:35
		qw422016.E().S(lc.Get("ui.notexist_media_tip1"))
//line views/hypha.qtpl:35
		qw422016.N().S(`</p>
		<form action="/upload-binary/`)
//line views/hypha.qtpl:36
		qw422016.E().S(h.CanonicalName())
//line views/hypha.qtpl:36
		qw422016.N().S(`"
        		method="post" enctype="multipart/form-data"
        		class="upload-binary">
        	<label for="upload-binary__input"></label>
        	<input type="file" id="upload-binary__input" name="binary">

        	<button type="submit" class="btn stick-to-bottom" value="Upload">`)
//line views/hypha.qtpl:42
		qw422016.E().S(lc.Get("ui.attach_upload"))
//line views/hypha.qtpl:42
		qw422016.N().S(`</button>
        </form>
	</section>
	</div>
	`)
//line views/hypha.qtpl:46
	}
//line views/hypha.qtpl:46
	qw422016.N().S(`
</section>
`)
//line views/hypha.qtpl:48
}

//line views/hypha.qtpl:48
func writenonExistentHyphaNotice(qq422016 qtio422016.Writer, h hyphae.Hypha, u *user.User, lc *l18n.Localizer) {
//line views/hypha.qtpl:48
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/hypha.qtpl:48
	streamnonExistentHyphaNotice(qw422016, h, u, lc)
//line views/hypha.qtpl:48
	qt422016.ReleaseWriter(qw422016)
//line views/hypha.qtpl:48
}

//line views/hypha.qtpl:48
func nonExistentHyphaNotice(h hyphae.Hypha, u *user.User, lc *l18n.Localizer) string {
//line views/hypha.qtpl:48
	qb422016 := qt422016.AcquireByteBuffer()
//line views/hypha.qtpl:48
	writenonExistentHyphaNotice(qb422016, h, u, lc)
//line views/hypha.qtpl:48
	qs422016 := string(qb422016.B)
//line views/hypha.qtpl:48
	qt422016.ReleaseByteBuffer(qb422016)
//line views/hypha.qtpl:48
	return qs422016
//line views/hypha.qtpl:48
}

//line views/hypha.qtpl:50
func StreamNaviTitleHTML(qw422016 *qt422016.Writer, h hyphae.Hypha) {
//line views/hypha.qtpl:50
	qw422016.N().S(`
`)
//line views/hypha.qtpl:52
	var (
		prevAcc = "/hypha/"
		parts   = strings.Split(h.CanonicalName(), "/")
	)

//line views/hypha.qtpl:56
	qw422016.N().S(`
<h1 class="navi-title">
`)
//line views/hypha.qtpl:58
	qw422016.N().S(`<a href="/hypha/`)
//line views/hypha.qtpl:59
	qw422016.E().S(cfg.HomeHypha)
//line views/hypha.qtpl:59
	qw422016.N().S(`">`)
//line views/hypha.qtpl:60
	qw422016.N().S(cfg.NaviTitleIcon)
//line views/hypha.qtpl:60
	qw422016.N().S(`<span aria-hidden="true" class="navi-title__colon">:</span></a>`)
//line views/hypha.qtpl:64
	for i, part := range parts {
//line views/hypha.qtpl:65
		if i > 0 {
//line views/hypha.qtpl:65
			qw422016.N().S(`<span aria-hidden="true" class="navi-title__separator">/</span>`)
//line views/hypha.qtpl:67
		}
//line views/hypha.qtpl:67
		qw422016.N().S(`<a href="`)
//line views/hypha.qtpl:69
		qw422016.E().S(prevAcc + part)
//line views/hypha.qtpl:69
		qw422016.N().S(`" rel="`)
//line views/hypha.qtpl:69
		if i == len(parts)-1 {
//line views/hypha.qtpl:69
			qw422016.N().S(`bookmark`)
//line views/hypha.qtpl:69
		} else {
//line views/hypha.qtpl:69
			qw422016.N().S(`up`)
//line views/hypha.qtpl:69
		}
//line views/hypha.qtpl:69
		qw422016.N().S(`">`)
//line views/hypha.qtpl:70
		qw422016.N().S(util.BeautifulName(part))
//line views/hypha.qtpl:70
		qw422016.N().S(`</a>`)
//line views/hypha.qtpl:72
		prevAcc += part + "/"

//line views/hypha.qtpl:73
	}
//line views/hypha.qtpl:74
	qw422016.N().S(`
</h1>
`)
//line views/hypha.qtpl:76
}

//line views/hypha.qtpl:76
func WriteNaviTitleHTML(qq422016 qtio422016.Writer, h hyphae.Hypha) {
//line views/hypha.qtpl:76
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/hypha.qtpl:76
	StreamNaviTitleHTML(qw422016, h)
//line views/hypha.qtpl:76
	qt422016.ReleaseWriter(qw422016)
//line views/hypha.qtpl:76
}

//line views/hypha.qtpl:76
func NaviTitleHTML(h hyphae.Hypha) string {
//line views/hypha.qtpl:76
	qb422016 := qt422016.AcquireByteBuffer()
//line views/hypha.qtpl:76
	WriteNaviTitleHTML(qb422016, h)
//line views/hypha.qtpl:76
	qs422016 := string(qb422016.B)
//line views/hypha.qtpl:76
	qt422016.ReleaseByteBuffer(qb422016)
//line views/hypha.qtpl:76
	return qs422016
//line views/hypha.qtpl:76
}

//line views/hypha.qtpl:78
func StreamAttachmentHTMLRaw(qw422016 *qt422016.Writer, h *hyphae.MediaHypha) {
//line views/hypha.qtpl:78
	StreamAttachmentHTML(qw422016, h, l18n.New("en", "en"))
//line views/hypha.qtpl:78
}

//line views/hypha.qtpl:78
func WriteAttachmentHTMLRaw(qq422016 qtio422016.Writer, h *hyphae.MediaHypha) {
//line views/hypha.qtpl:78
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/hypha.qtpl:78
	StreamAttachmentHTMLRaw(qw422016, h)
//line views/hypha.qtpl:78
	qt422016.ReleaseWriter(qw422016)
//line views/hypha.qtpl:78
}

//line views/hypha.qtpl:78
func AttachmentHTMLRaw(h *hyphae.MediaHypha) string {
//line views/hypha.qtpl:78
	qb422016 := qt422016.AcquireByteBuffer()
//line views/hypha.qtpl:78
	WriteAttachmentHTMLRaw(qb422016, h)
//line views/hypha.qtpl:78
	qs422016 := string(qb422016.B)
//line views/hypha.qtpl:78
	qt422016.ReleaseByteBuffer(qb422016)
//line views/hypha.qtpl:78
	return qs422016
//line views/hypha.qtpl:78
}

//line views/hypha.qtpl:80
func StreamAttachmentHTML(qw422016 *qt422016.Writer, h *hyphae.MediaHypha, lc *l18n.Localizer) {
//line views/hypha.qtpl:80
	qw422016.N().S(`
	`)
//line views/hypha.qtpl:81
	switch filepath.Ext(h.MediaFilePath()) {
//line views/hypha.qtpl:83
	case ".jpg", ".gif", ".png", ".webp", ".svg", ".ico":
//line views/hypha.qtpl:83
		qw422016.N().S(`
	<div class="binary-container binary-container_with-img">
		<a href="/binary/`)
//line views/hypha.qtpl:85
		qw422016.N().S(h.CanonicalName())
//line views/hypha.qtpl:85
		qw422016.N().S(`"><img src="/binary/`)
//line views/hypha.qtpl:85
		qw422016.N().S(h.CanonicalName())
//line views/hypha.qtpl:85
		qw422016.N().S(`"/></a>
	</div>

	`)
//line views/hypha.qtpl:88
	case ".ogg", ".webm", ".mp4":
//line views/hypha.qtpl:88
		qw422016.N().S(`
	<div class="binary-container binary-container_with-video">
		<video controls>
			<source src="/binary/`)
//line views/hypha.qtpl:91
		qw422016.N().S(h.CanonicalName())
//line views/hypha.qtpl:91
		qw422016.N().S(`"/>
			<p>`)
//line views/hypha.qtpl:92
		qw422016.E().S(lc.Get("ui.media_novideo"))
//line views/hypha.qtpl:92
		qw422016.N().S(` <a href="/binary/`)
//line views/hypha.qtpl:92
		qw422016.N().S(h.CanonicalName())
//line views/hypha.qtpl:92
		qw422016.N().S(`">`)
//line views/hypha.qtpl:92
		qw422016.E().S(lc.Get("ui.media_novideo_link"))
//line views/hypha.qtpl:92
		qw422016.N().S(`</a></p>
		</video>
	</div>

	`)
//line views/hypha.qtpl:96
	case ".mp3":
//line views/hypha.qtpl:96
		qw422016.N().S(`
	<div class="binary-container binary-container_with-audio">
		<audio controls>
			<source src="/binary/`)
//line views/hypha.qtpl:99
		qw422016.N().S(h.CanonicalName())
//line views/hypha.qtpl:99
		qw422016.N().S(`"/>
			<p>`)
//line views/hypha.qtpl:100
		qw422016.E().S(lc.Get("ui.media_noaudio"))
//line views/hypha.qtpl:100
		qw422016.N().S(` <a href="/binary/`)
//line views/hypha.qtpl:100
		qw422016.N().S(h.CanonicalName())
//line views/hypha.qtpl:100
		qw422016.N().S(`">`)
//line views/hypha.qtpl:100
		qw422016.E().S(lc.Get("ui.media_noaudio_link"))
//line views/hypha.qtpl:100
		qw422016.N().S(`</a></p>
		</audio>
	</div>

	`)
//line views/hypha.qtpl:104
	default:
//line views/hypha.qtpl:104
		qw422016.N().S(`
	<div class="binary-container binary-container_with-nothing">
		<p><a href="/binary/`)
//line views/hypha.qtpl:106
		qw422016.N().S(h.CanonicalName())
//line views/hypha.qtpl:106
		qw422016.N().S(`">`)
//line views/hypha.qtpl:106
		qw422016.E().S(lc.Get("ui.media_download"))
//line views/hypha.qtpl:106
		qw422016.N().S(`</a></p>
	</div>
`)
//line views/hypha.qtpl:108
	}
//line views/hypha.qtpl:108
	qw422016.N().S(`
`)
//line views/hypha.qtpl:109
}

//line views/hypha.qtpl:109
func WriteAttachmentHTML(qq422016 qtio422016.Writer, h *hyphae.MediaHypha, lc *l18n.Localizer) {
//line views/hypha.qtpl:109
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/hypha.qtpl:109
	StreamAttachmentHTML(qw422016, h, lc)
//line views/hypha.qtpl:109
	qt422016.ReleaseWriter(qw422016)
//line views/hypha.qtpl:109
}

//line views/hypha.qtpl:109
func AttachmentHTML(h *hyphae.MediaHypha, lc *l18n.Localizer) string {
//line views/hypha.qtpl:109
	qb422016 := qt422016.AcquireByteBuffer()
//line views/hypha.qtpl:109
	WriteAttachmentHTML(qb422016, h, lc)
//line views/hypha.qtpl:109
	qs422016 := string(qb422016.B)
//line views/hypha.qtpl:109
	qt422016.ReleaseByteBuffer(qb422016)
//line views/hypha.qtpl:109
	return qs422016
//line views/hypha.qtpl:109
}
