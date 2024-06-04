// Code generated by qtc from "readers.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line hypview/readers.qtpl:1
package hypview

//line hypview/readers.qtpl:1
import "net/http"

//line hypview/readers.qtpl:2
import "strings"

//line hypview/readers.qtpl:3
import "path"

//line hypview/readers.qtpl:4
import "os"

//line hypview/readers.qtpl:6
import "github.com/bouncepaw/mycorrhiza/internal/cfg"

//line hypview/readers.qtpl:7
import "github.com/bouncepaw/mycorrhiza/internal/hyphae"

//line hypview/readers.qtpl:8
import "github.com/bouncepaw/mycorrhiza/categories"

//line hypview/readers.qtpl:9
import "github.com/bouncepaw/mycorrhiza/l18n"

//line hypview/readers.qtpl:10
import "github.com/bouncepaw/mycorrhiza/internal/mimetype"

//line hypview/readers.qtpl:11
import "github.com/bouncepaw/mycorrhiza/tree"

//line hypview/readers.qtpl:12
import "github.com/bouncepaw/mycorrhiza/internal/user"

//line hypview/readers.qtpl:13
import "github.com/bouncepaw/mycorrhiza/util"

//line hypview/readers.qtpl:14
import "github.com/bouncepaw/mycorrhiza/web/viewutil"

//line hypview/readers.qtpl:16
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line hypview/readers.qtpl:16
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line hypview/readers.qtpl:16
func StreamMediaMenu(qw422016 *qt422016.Writer, rq *http.Request, h hyphae.Hypha, u *user.User) {
//line hypview/readers.qtpl:16
	qw422016.N().S(`
`)
//line hypview/readers.qtpl:18
	lc := l18n.FromRequest(rq)

//line hypview/readers.qtpl:19
	qw422016.N().S(`
<main class="main-width media-tab">
	<h1>`)
//line hypview/readers.qtpl:21
	qw422016.N().S(lc.Get("ui.media_title", &l18n.Replacements{"name": beautifulLink(h.CanonicalName())}))
//line hypview/readers.qtpl:21
	qw422016.N().S(`</h1>
	`)
//line hypview/readers.qtpl:22
	switch h.(type) {
//line hypview/readers.qtpl:23
	case *hyphae.MediaHypha:
//line hypview/readers.qtpl:23
		qw422016.N().S(`
		<p class="explanation">`)
//line hypview/readers.qtpl:24
		qw422016.E().S(lc.Get("ui.media_tip"))
//line hypview/readers.qtpl:24
		qw422016.N().S(` <a href="/help/en/media" class="shy-link">`)
//line hypview/readers.qtpl:24
		qw422016.E().S(lc.Get("ui.media_what_is"))
//line hypview/readers.qtpl:24
		qw422016.N().S(`</a></p>
	`)
//line hypview/readers.qtpl:25
	default:
//line hypview/readers.qtpl:25
		qw422016.N().S(`
		<p class="explanation">`)
//line hypview/readers.qtpl:26
		qw422016.E().S(lc.Get("ui.media_empty"))
//line hypview/readers.qtpl:26
		qw422016.N().S(` <a href="/help/en/media" class="shy-link">`)
//line hypview/readers.qtpl:26
		qw422016.E().S(lc.Get("ui.media_what_is"))
//line hypview/readers.qtpl:26
		qw422016.N().S(`</a></p>
	`)
//line hypview/readers.qtpl:27
	}
//line hypview/readers.qtpl:27
	qw422016.N().S(`

	<section class="amnt-grid">
	`)
//line hypview/readers.qtpl:30
	switch h := h.(type) {
//line hypview/readers.qtpl:31
	case *hyphae.MediaHypha:
//line hypview/readers.qtpl:31
		qw422016.N().S(`
		`)
//line hypview/readers.qtpl:33
		mime := mimetype.FromExtension(path.Ext(h.MediaFilePath()))
		fileinfo, err := os.Stat(h.MediaFilePath())

//line hypview/readers.qtpl:34
		qw422016.N().S(`
		`)
//line hypview/readers.qtpl:35
		if err == nil {
//line hypview/readers.qtpl:35
			qw422016.N().S(`
		<fieldset class="amnt-menu-block">
			<legend class="modal__title modal__title_small">`)
//line hypview/readers.qtpl:37
			qw422016.E().S(lc.Get("ui.media_stat"))
//line hypview/readers.qtpl:37
			qw422016.N().S(`</legend>
			<p class="modal__confirmation-msg"><b>`)
//line hypview/readers.qtpl:38
			qw422016.E().S(lc.Get("ui.media_stat_size"))
//line hypview/readers.qtpl:38
			qw422016.N().S(`</b> `)
//line hypview/readers.qtpl:38
			qw422016.E().S(lc.GetPlural64("ui.media_size_value", fileinfo.Size()))
//line hypview/readers.qtpl:38
			qw422016.N().S(`</p>
			<p><b>`)
//line hypview/readers.qtpl:39
			qw422016.E().S(lc.Get("ui.media_stat_mime"))
//line hypview/readers.qtpl:39
			qw422016.N().S(`</b> `)
//line hypview/readers.qtpl:39
			qw422016.E().S(mime)
//line hypview/readers.qtpl:39
			qw422016.N().S(`</p>
		</fieldset>
		`)
//line hypview/readers.qtpl:41
		}
//line hypview/readers.qtpl:41
		qw422016.N().S(`

		`)
//line hypview/readers.qtpl:43
		if strings.HasPrefix(mime, "image/") {
//line hypview/readers.qtpl:43
			qw422016.N().S(`
		<fieldset class="amnt-menu-block">
			<legend class="modal__title modal__title_small">`)
//line hypview/readers.qtpl:45
			qw422016.E().S(lc.Get("ui.media_include"))
//line hypview/readers.qtpl:45
			qw422016.N().S(`</legend>
			<p class="modal__confirmation-msg">`)
//line hypview/readers.qtpl:46
			qw422016.E().S(lc.Get("ui.media_include_tip"))
//line hypview/readers.qtpl:46
			qw422016.N().S(`</p>
			<pre class="codeblock"><code>img { `)
//line hypview/readers.qtpl:47
			qw422016.E().S(h.CanonicalName())
//line hypview/readers.qtpl:47
			qw422016.N().S(` }</code></pre>
		</fieldset>
		`)
//line hypview/readers.qtpl:49
		}
//line hypview/readers.qtpl:49
		qw422016.N().S(`
	`)
//line hypview/readers.qtpl:50
	}
//line hypview/readers.qtpl:50
	qw422016.N().S(`

	`)
//line hypview/readers.qtpl:52
	if u.CanProceed("upload-binary") {
//line hypview/readers.qtpl:52
		qw422016.N().S(`
	<form action="/upload-binary/`)
//line hypview/readers.qtpl:53
		qw422016.E().S(h.CanonicalName())
//line hypview/readers.qtpl:53
		qw422016.N().S(`"
			method="post" enctype="multipart/form-data"
			class="upload-binary modal amnt-menu-block">
		<fieldset class="modal__fieldset">
			<legend class="modal__title modal__title_small">`)
//line hypview/readers.qtpl:57
		qw422016.E().S(lc.Get("ui.media_new"))
//line hypview/readers.qtpl:57
		qw422016.N().S(`</legend>
			<p class="modal__confirmation-msg">`)
//line hypview/readers.qtpl:58
		qw422016.E().S(lc.Get("ui.media_new_tip"))
//line hypview/readers.qtpl:58
		qw422016.N().S(`</p>
			<label for="upload-binary__input"></label>
			<input type="file" id="upload-binary__input" name="binary">

			<button type="submit" class="btn stick-to-bottom" value="Upload">`)
//line hypview/readers.qtpl:62
		qw422016.E().S(lc.Get("ui.media_upload"))
//line hypview/readers.qtpl:62
		qw422016.N().S(`</button>
		</fieldset>
	</form>
	`)
//line hypview/readers.qtpl:65
	}
//line hypview/readers.qtpl:65
	qw422016.N().S(`


	`)
//line hypview/readers.qtpl:68
	switch h := h.(type) {
//line hypview/readers.qtpl:69
	case *hyphae.MediaHypha:
//line hypview/readers.qtpl:69
		qw422016.N().S(`
		`)
//line hypview/readers.qtpl:70
		if u.CanProceed("remove-media") {
//line hypview/readers.qtpl:70
			qw422016.N().S(`
		<form action="/remove-media/`)
//line hypview/readers.qtpl:71
			qw422016.E().S(h.CanonicalName())
//line hypview/readers.qtpl:71
			qw422016.N().S(`" method="post" class="modal amnt-menu-block" method="POST">
			<fieldset class="modal__fieldset">
				<legend class="modal__title modal__title_small">`)
//line hypview/readers.qtpl:73
			qw422016.E().S(lc.Get("ui.media_remove"))
//line hypview/readers.qtpl:73
			qw422016.N().S(`</legend>
				<p class="modal__confirmation-msg">`)
//line hypview/readers.qtpl:74
			qw422016.E().S(lc.Get("ui.media_remove_tip"))
//line hypview/readers.qtpl:74
			qw422016.N().S(`</p>
				<button type="submit" class="btn" value="Remove media">`)
//line hypview/readers.qtpl:75
			qw422016.E().S(lc.Get("ui.media_remove_button"))
//line hypview/readers.qtpl:75
			qw422016.N().S(`</button>
			</fieldset>
		</form>
		`)
//line hypview/readers.qtpl:78
		}
//line hypview/readers.qtpl:78
		qw422016.N().S(`
	`)
//line hypview/readers.qtpl:79
	}
//line hypview/readers.qtpl:79
	qw422016.N().S(`

	</section>
</main>
`)
//line hypview/readers.qtpl:83
}

//line hypview/readers.qtpl:83
func WriteMediaMenu(qq422016 qtio422016.Writer, rq *http.Request, h hyphae.Hypha, u *user.User) {
//line hypview/readers.qtpl:83
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/readers.qtpl:83
	StreamMediaMenu(qw422016, rq, h, u)
//line hypview/readers.qtpl:83
	qt422016.ReleaseWriter(qw422016)
//line hypview/readers.qtpl:83
}

//line hypview/readers.qtpl:83
func MediaMenu(rq *http.Request, h hyphae.Hypha, u *user.User) string {
//line hypview/readers.qtpl:83
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/readers.qtpl:83
	WriteMediaMenu(qb422016, rq, h, u)
//line hypview/readers.qtpl:83
	qs422016 := string(qb422016.B)
//line hypview/readers.qtpl:83
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/readers.qtpl:83
	return qs422016
//line hypview/readers.qtpl:83
}

// If `contents` == "", a helpful message is shown instead.
//
// If you rename .prevnext, change the docs too.

//line hypview/readers.qtpl:88
func StreamHypha(qw422016 *qt422016.Writer, meta viewutil.Meta, h hyphae.Hypha, contents string) {
//line hypview/readers.qtpl:88
	qw422016.N().S(`
`)
//line hypview/readers.qtpl:90
	subhyphae, prevHyphaName, nextHyphaName := tree.Tree(h.CanonicalName())
	lc := meta.Lc

//line hypview/readers.qtpl:92
	qw422016.N().S(`
<main class="main-width">
	<section id="hypha">
		`)
//line hypview/readers.qtpl:95
	if meta.U.CanProceed("edit") {
//line hypview/readers.qtpl:95
		qw422016.N().S(`
		<div class="btn btn_navititle">
			<a class="btn__link_navititle" href="/edit/`)
//line hypview/readers.qtpl:97
		qw422016.E().S(h.CanonicalName())
//line hypview/readers.qtpl:97
		qw422016.N().S(`">`)
//line hypview/readers.qtpl:97
		qw422016.E().S(lc.Get("ui.edit_link"))
//line hypview/readers.qtpl:97
		qw422016.N().S(`</a>
		</div>
		`)
//line hypview/readers.qtpl:99
	}
//line hypview/readers.qtpl:99
	qw422016.N().S(`

		`)
//line hypview/readers.qtpl:101
	if cfg.UseAuth && util.IsProfileName(h.CanonicalName()) && meta.U.Name == strings.TrimPrefix(h.CanonicalName(), cfg.UserHypha+"/") {
//line hypview/readers.qtpl:101
		qw422016.N().S(`
		<div class="btn btn_navititle">
			<a class="btn__link_navititle" href="/logout">`)
//line hypview/readers.qtpl:103
		qw422016.E().S(lc.Get("ui.logout_link"))
//line hypview/readers.qtpl:103
		qw422016.N().S(`</a>
		</div>
		`)
//line hypview/readers.qtpl:105
		if meta.U.Group == "admin" {
//line hypview/readers.qtpl:105
			qw422016.N().S(`
		<div class="btn btn_navititle">
			<a class="btn__link_navititle" href="/admin">`)
//line hypview/readers.qtpl:107
			qw422016.E().S(lc.Get("ui.admin_panel"))
//line hypview/readers.qtpl:107
			qw422016.N().S(`<a>
		</div>
		`)
//line hypview/readers.qtpl:109
		}
//line hypview/readers.qtpl:109
		qw422016.N().S(`
		`)
//line hypview/readers.qtpl:110
	}
//line hypview/readers.qtpl:110
	qw422016.N().S(`

		`)
//line hypview/readers.qtpl:112
	qw422016.N().S(NaviTitle(meta, h.CanonicalName()))
//line hypview/readers.qtpl:112
	qw422016.N().S(`
		`)
//line hypview/readers.qtpl:113
	switch h.(type) {
//line hypview/readers.qtpl:114
	case *hyphae.EmptyHypha:
//line hypview/readers.qtpl:114
		qw422016.N().S(`
				`)
//line hypview/readers.qtpl:115
		qw422016.N().S(EmptyHypha(meta, h.CanonicalName()))
//line hypview/readers.qtpl:115
		qw422016.N().S(`
			`)
//line hypview/readers.qtpl:116
	default:
//line hypview/readers.qtpl:116
		qw422016.N().S(`
				`)
//line hypview/readers.qtpl:117
		qw422016.N().S(contents)
//line hypview/readers.qtpl:117
		qw422016.N().S(`
		`)
//line hypview/readers.qtpl:118
	}
//line hypview/readers.qtpl:118
	qw422016.N().S(`
	</section>
	<section class="prevnext">
		`)
//line hypview/readers.qtpl:121
	if prevHyphaName != "" {
//line hypview/readers.qtpl:121
		qw422016.N().S(`
		<a class="prevnext__el prevnext__prev" href="/hypha/`)
//line hypview/readers.qtpl:122
		qw422016.E().S(prevHyphaName)
//line hypview/readers.qtpl:122
		qw422016.N().S(`" rel="prev">← `)
//line hypview/readers.qtpl:122
		qw422016.E().S(util.BeautifulName(path.Base(prevHyphaName)))
//line hypview/readers.qtpl:122
		qw422016.N().S(`</a>
		`)
//line hypview/readers.qtpl:123
	}
//line hypview/readers.qtpl:123
	qw422016.N().S(`
		`)
//line hypview/readers.qtpl:124
	if nextHyphaName != "" {
//line hypview/readers.qtpl:124
		qw422016.N().S(`
		<a class="prevnext__el prevnext__next" href="/hypha/`)
//line hypview/readers.qtpl:125
		qw422016.E().S(nextHyphaName)
//line hypview/readers.qtpl:125
		qw422016.N().S(`" rel="next">`)
//line hypview/readers.qtpl:125
		qw422016.E().S(util.BeautifulName(path.Base(nextHyphaName)))
//line hypview/readers.qtpl:125
		qw422016.N().S(` →</a>
		`)
//line hypview/readers.qtpl:126
	}
//line hypview/readers.qtpl:126
	qw422016.N().S(`
	</section>
`)
//line hypview/readers.qtpl:128
	if strings.TrimSpace(subhyphae) != "" {
//line hypview/readers.qtpl:128
		qw422016.N().S(`
<section class="subhyphae">
	<h2 class="subhyphae__title">`)
//line hypview/readers.qtpl:130
		qw422016.E().S(lc.Get("ui.subhyphae"))
//line hypview/readers.qtpl:130
		qw422016.N().S(`</h2>
	<nav class="subhyphae__nav">
		<ul class="subhyphae__list">
		`)
//line hypview/readers.qtpl:133
		qw422016.N().S(subhyphae)
//line hypview/readers.qtpl:133
		qw422016.N().S(`
		</ul>
	</nav>
</section>
`)
//line hypview/readers.qtpl:137
	}
//line hypview/readers.qtpl:137
	qw422016.N().S(`
	<section id="hypha-bottom">
   		`)
//line hypview/readers.qtpl:139
	streamhyphaInfo(qw422016, meta, h)
//line hypview/readers.qtpl:139
	qw422016.N().S(`
	</section>
</main>
`)
//line hypview/readers.qtpl:142
	qw422016.N().S(categories.CategoryCard(meta, h.CanonicalName()))
//line hypview/readers.qtpl:142
	qw422016.N().S(`
`)
//line hypview/readers.qtpl:143
	streamviewScripts(qw422016)
//line hypview/readers.qtpl:143
	qw422016.N().S(`
`)
//line hypview/readers.qtpl:144
}

//line hypview/readers.qtpl:144
func WriteHypha(qq422016 qtio422016.Writer, meta viewutil.Meta, h hyphae.Hypha, contents string) {
//line hypview/readers.qtpl:144
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/readers.qtpl:144
	StreamHypha(qw422016, meta, h, contents)
//line hypview/readers.qtpl:144
	qt422016.ReleaseWriter(qw422016)
//line hypview/readers.qtpl:144
}

//line hypview/readers.qtpl:144
func Hypha(meta viewutil.Meta, h hyphae.Hypha, contents string) string {
//line hypview/readers.qtpl:144
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/readers.qtpl:144
	WriteHypha(qb422016, meta, h, contents)
//line hypview/readers.qtpl:144
	qs422016 := string(qb422016.B)
//line hypview/readers.qtpl:144
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/readers.qtpl:144
	return qs422016
//line hypview/readers.qtpl:144
}

//line hypview/readers.qtpl:146
func StreamRevision(qw422016 *qt422016.Writer, meta viewutil.Meta, h hyphae.Hypha, contents, revHash string) {
//line hypview/readers.qtpl:146
	qw422016.N().S(`
<main class="main-width">
	<section>
		<p>`)
//line hypview/readers.qtpl:149
	qw422016.E().S(meta.Lc.Get("ui.revision_warning"))
//line hypview/readers.qtpl:149
	qw422016.N().S(` <a href="/rev-text/`)
//line hypview/readers.qtpl:149
	qw422016.E().S(revHash)
//line hypview/readers.qtpl:149
	qw422016.N().S(`/`)
//line hypview/readers.qtpl:149
	qw422016.E().S(h.CanonicalName())
//line hypview/readers.qtpl:149
	qw422016.N().S(`">`)
//line hypview/readers.qtpl:149
	qw422016.E().S(meta.Lc.Get("ui.revision_link"))
//line hypview/readers.qtpl:149
	qw422016.N().S(`</a></p>
		`)
//line hypview/readers.qtpl:150
	qw422016.N().S(NaviTitle(meta, h.CanonicalName()))
//line hypview/readers.qtpl:150
	qw422016.N().S(`
		`)
//line hypview/readers.qtpl:151
	qw422016.N().S(contents)
//line hypview/readers.qtpl:151
	qw422016.N().S(`
	</section>
</main>
`)
//line hypview/readers.qtpl:154
	streamviewScripts(qw422016)
//line hypview/readers.qtpl:154
	qw422016.N().S(`
`)
//line hypview/readers.qtpl:155
}

//line hypview/readers.qtpl:155
func WriteRevision(qq422016 qtio422016.Writer, meta viewutil.Meta, h hyphae.Hypha, contents, revHash string) {
//line hypview/readers.qtpl:155
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/readers.qtpl:155
	StreamRevision(qw422016, meta, h, contents, revHash)
//line hypview/readers.qtpl:155
	qt422016.ReleaseWriter(qw422016)
//line hypview/readers.qtpl:155
}

//line hypview/readers.qtpl:155
func Revision(meta viewutil.Meta, h hyphae.Hypha, contents, revHash string) string {
//line hypview/readers.qtpl:155
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/readers.qtpl:155
	WriteRevision(qb422016, meta, h, contents, revHash)
//line hypview/readers.qtpl:155
	qs422016 := string(qb422016.B)
//line hypview/readers.qtpl:155
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/readers.qtpl:155
	return qs422016
//line hypview/readers.qtpl:155
}

//line hypview/readers.qtpl:157
func streamviewScripts(qw422016 *qt422016.Writer) {
//line hypview/readers.qtpl:157
	qw422016.N().S(`
`)
//line hypview/readers.qtpl:158
	for _, scriptPath := range cfg.ViewScripts {
//line hypview/readers.qtpl:158
		qw422016.N().S(`
<script src="`)
//line hypview/readers.qtpl:159
		qw422016.E().S(scriptPath)
//line hypview/readers.qtpl:159
		qw422016.N().S(`"></script>
`)
//line hypview/readers.qtpl:160
	}
//line hypview/readers.qtpl:160
	qw422016.N().S(`
`)
//line hypview/readers.qtpl:161
}

//line hypview/readers.qtpl:161
func writeviewScripts(qq422016 qtio422016.Writer) {
//line hypview/readers.qtpl:161
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/readers.qtpl:161
	streamviewScripts(qw422016)
//line hypview/readers.qtpl:161
	qt422016.ReleaseWriter(qw422016)
//line hypview/readers.qtpl:161
}

//line hypview/readers.qtpl:161
func viewScripts() string {
//line hypview/readers.qtpl:161
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/readers.qtpl:161
	writeviewScripts(qb422016)
//line hypview/readers.qtpl:161
	qs422016 := string(qb422016.B)
//line hypview/readers.qtpl:161
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/readers.qtpl:161
	return qs422016
//line hypview/readers.qtpl:161
}
