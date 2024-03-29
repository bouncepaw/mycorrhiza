// Code generated by qtc from "nav.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line hypview/nav.qtpl:1
package hypview

//line hypview/nav.qtpl:1
import "github.com/bouncepaw/mycorrhiza/backlinks"

//line hypview/nav.qtpl:2
import "github.com/bouncepaw/mycorrhiza/cfg"

//line hypview/nav.qtpl:3
import "github.com/bouncepaw/mycorrhiza/hyphae"

//line hypview/nav.qtpl:4
import "github.com/bouncepaw/mycorrhiza/user"

//line hypview/nav.qtpl:5
import "github.com/bouncepaw/mycorrhiza/util"

//line hypview/nav.qtpl:6
import "github.com/bouncepaw/mycorrhiza/viewutil"

//line hypview/nav.qtpl:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line hypview/nav.qtpl:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line hypview/nav.qtpl:8
func streamhyphaInfoEntry(qw422016 *qt422016.Writer, h hyphae.Hypha, u *user.User, action string, hasToExist bool, displayText string) {
//line hypview/nav.qtpl:8
	qw422016.N().S(`
`)
//line hypview/nav.qtpl:9
	flag := true

//line hypview/nav.qtpl:9
	qw422016.N().S(`
`)
//line hypview/nav.qtpl:10
	switch h.(type) {
//line hypview/nav.qtpl:11
	case *hyphae.EmptyHypha:
//line hypview/nav.qtpl:11
		qw422016.N().S(`
	`)
//line hypview/nav.qtpl:12
		flag = !hasToExist

//line hypview/nav.qtpl:12
		qw422016.N().S(`
`)
//line hypview/nav.qtpl:13
	}
//line hypview/nav.qtpl:13
	qw422016.N().S(`
`)
//line hypview/nav.qtpl:14
	if u.CanProceed(action) && flag {
//line hypview/nav.qtpl:14
		qw422016.N().S(`
<li class="hypha-info__entry hypha-info__entry_`)
//line hypview/nav.qtpl:15
		qw422016.E().S(action)
//line hypview/nav.qtpl:15
		qw422016.N().S(`">
	<a class="hypha-info__link" href="/`)
//line hypview/nav.qtpl:16
		qw422016.E().S(action)
//line hypview/nav.qtpl:16
		qw422016.N().S(`/`)
//line hypview/nav.qtpl:16
		qw422016.E().S(h.CanonicalName())
//line hypview/nav.qtpl:16
		qw422016.N().S(`">`)
//line hypview/nav.qtpl:16
		qw422016.E().S(displayText)
//line hypview/nav.qtpl:16
		qw422016.N().S(`</a>
</li>
`)
//line hypview/nav.qtpl:18
	}
//line hypview/nav.qtpl:18
	qw422016.N().S(`
`)
//line hypview/nav.qtpl:19
}

//line hypview/nav.qtpl:19
func writehyphaInfoEntry(qq422016 qtio422016.Writer, h hyphae.Hypha, u *user.User, action string, hasToExist bool, displayText string) {
//line hypview/nav.qtpl:19
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/nav.qtpl:19
	streamhyphaInfoEntry(qw422016, h, u, action, hasToExist, displayText)
//line hypview/nav.qtpl:19
	qt422016.ReleaseWriter(qw422016)
//line hypview/nav.qtpl:19
}

//line hypview/nav.qtpl:19
func hyphaInfoEntry(h hyphae.Hypha, u *user.User, action string, hasToExist bool, displayText string) string {
//line hypview/nav.qtpl:19
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/nav.qtpl:19
	writehyphaInfoEntry(qb422016, h, u, action, hasToExist, displayText)
//line hypview/nav.qtpl:19
	qs422016 := string(qb422016.B)
//line hypview/nav.qtpl:19
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/nav.qtpl:19
	return qs422016
//line hypview/nav.qtpl:19
}

//line hypview/nav.qtpl:21
func streamhyphaInfo(qw422016 *qt422016.Writer, meta viewutil.Meta, h hyphae.Hypha) {
//line hypview/nav.qtpl:21
	qw422016.N().S(`
`)
//line hypview/nav.qtpl:23
	u := meta.U
	lc := meta.Lc
	backs := backlinks.BacklinksCount(h.CanonicalName())

//line hypview/nav.qtpl:26
	qw422016.N().S(`
<nav class="hypha-info">
	<ul class="hypha-info__list">
		`)
//line hypview/nav.qtpl:29
	streamhyphaInfoEntry(qw422016, h, u, "history", false, lc.Get("ui.history_link"))
//line hypview/nav.qtpl:29
	qw422016.N().S(`
		`)
//line hypview/nav.qtpl:30
	streamhyphaInfoEntry(qw422016, h, u, "rename", true, lc.Get("ui.rename_link"))
//line hypview/nav.qtpl:30
	qw422016.N().S(`
		`)
//line hypview/nav.qtpl:31
	streamhyphaInfoEntry(qw422016, h, u, "delete", true, lc.Get("ui.delete_link"))
//line hypview/nav.qtpl:31
	qw422016.N().S(`
		`)
//line hypview/nav.qtpl:32
	streamhyphaInfoEntry(qw422016, h, u, "text", true, lc.Get("ui.text_link"))
//line hypview/nav.qtpl:32
	qw422016.N().S(`
	`)
//line hypview/nav.qtpl:33
	switch h := h.(type) {
//line hypview/nav.qtpl:34
	case *hyphae.TextualHypha:
//line hypview/nav.qtpl:34
		qw422016.N().S(`
		`)
//line hypview/nav.qtpl:35
		streamhyphaInfoEntry(qw422016, h, u, "media", true, lc.Get("ui.media_link_for_textual"))
//line hypview/nav.qtpl:35
		qw422016.N().S(`
	`)
//line hypview/nav.qtpl:36
	default:
//line hypview/nav.qtpl:36
		qw422016.N().S(`
		`)
//line hypview/nav.qtpl:37
		streamhyphaInfoEntry(qw422016, h, u, "media", true, lc.Get("ui.media_link"))
//line hypview/nav.qtpl:37
		qw422016.N().S(`
	`)
//line hypview/nav.qtpl:38
	}
//line hypview/nav.qtpl:38
	qw422016.N().S(`
		`)
//line hypview/nav.qtpl:39
	streamhyphaInfoEntry(qw422016, h, u, "backlinks", false, lc.GetPlural("ui.backlinks_link", backs))
//line hypview/nav.qtpl:39
	qw422016.N().S(`
	</ul>
</nav>
`)
//line hypview/nav.qtpl:42
}

//line hypview/nav.qtpl:42
func writehyphaInfo(qq422016 qtio422016.Writer, meta viewutil.Meta, h hyphae.Hypha) {
//line hypview/nav.qtpl:42
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/nav.qtpl:42
	streamhyphaInfo(qw422016, meta, h)
//line hypview/nav.qtpl:42
	qt422016.ReleaseWriter(qw422016)
//line hypview/nav.qtpl:42
}

//line hypview/nav.qtpl:42
func hyphaInfo(meta viewutil.Meta, h hyphae.Hypha) string {
//line hypview/nav.qtpl:42
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/nav.qtpl:42
	writehyphaInfo(qb422016, meta, h)
//line hypview/nav.qtpl:42
	qs422016 := string(qb422016.B)
//line hypview/nav.qtpl:42
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/nav.qtpl:42
	return qs422016
//line hypview/nav.qtpl:42
}

//line hypview/nav.qtpl:44
func streamcommonScripts(qw422016 *qt422016.Writer) {
//line hypview/nav.qtpl:44
	qw422016.N().S(`
`)
//line hypview/nav.qtpl:45
	for _, scriptPath := range cfg.CommonScripts {
//line hypview/nav.qtpl:45
		qw422016.N().S(`
<script src="`)
//line hypview/nav.qtpl:46
		qw422016.E().S(scriptPath)
//line hypview/nav.qtpl:46
		qw422016.N().S(`"></script>
`)
//line hypview/nav.qtpl:47
	}
//line hypview/nav.qtpl:47
	qw422016.N().S(`
`)
//line hypview/nav.qtpl:48
}

//line hypview/nav.qtpl:48
func writecommonScripts(qq422016 qtio422016.Writer) {
//line hypview/nav.qtpl:48
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/nav.qtpl:48
	streamcommonScripts(qw422016)
//line hypview/nav.qtpl:48
	qt422016.ReleaseWriter(qw422016)
//line hypview/nav.qtpl:48
}

//line hypview/nav.qtpl:48
func commonScripts() string {
//line hypview/nav.qtpl:48
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/nav.qtpl:48
	writecommonScripts(qb422016)
//line hypview/nav.qtpl:48
	qs422016 := string(qb422016.B)
//line hypview/nav.qtpl:48
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/nav.qtpl:48
	return qs422016
//line hypview/nav.qtpl:48
}

//line hypview/nav.qtpl:50
func streambeautifulLink(qw422016 *qt422016.Writer, hyphaName string) {
//line hypview/nav.qtpl:50
	qw422016.N().S(`<a href="/hypha/`)
//line hypview/nav.qtpl:50
	qw422016.N().S(hyphaName)
//line hypview/nav.qtpl:50
	qw422016.N().S(`">`)
//line hypview/nav.qtpl:50
	qw422016.E().S(util.BeautifulName(hyphaName))
//line hypview/nav.qtpl:50
	qw422016.N().S(`</a>`)
//line hypview/nav.qtpl:50
}

//line hypview/nav.qtpl:50
func writebeautifulLink(qq422016 qtio422016.Writer, hyphaName string) {
//line hypview/nav.qtpl:50
	qw422016 := qt422016.AcquireWriter(qq422016)
//line hypview/nav.qtpl:50
	streambeautifulLink(qw422016, hyphaName)
//line hypview/nav.qtpl:50
	qt422016.ReleaseWriter(qw422016)
//line hypview/nav.qtpl:50
}

//line hypview/nav.qtpl:50
func beautifulLink(hyphaName string) string {
//line hypview/nav.qtpl:50
	qb422016 := qt422016.AcquireByteBuffer()
//line hypview/nav.qtpl:50
	writebeautifulLink(qb422016, hyphaName)
//line hypview/nav.qtpl:50
	qs422016 := string(qb422016.B)
//line hypview/nav.qtpl:50
	qt422016.ReleaseByteBuffer(qb422016)
//line hypview/nav.qtpl:50
	return qs422016
//line hypview/nav.qtpl:50
}
