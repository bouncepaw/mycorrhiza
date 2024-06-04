// Code generated by qtc from "view.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line history/view.qtpl:1
package history

//line history/view.qtpl:1
import "fmt"

//line history/view.qtpl:2
import "github.com/bouncepaw/mycorrhiza/internal/cfg"

// HyphaeLinksHTML returns a comma-separated list of hyphae that were affected by this revision as HTML string.

//line history/view.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line history/view.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line history/view.qtpl:5
func (rev Revision) StreamHyphaeLinksHTML(qw422016 *qt422016.Writer) {
//line history/view.qtpl:5
	qw422016.N().S(`
`)
//line history/view.qtpl:7
	for i, hyphaName := range rev.hyphaeAffected() {
//line history/view.qtpl:8
		if i > 0 {
//line history/view.qtpl:8
			qw422016.N().S(`<span aria-hidden="true">, </span>`)
//line history/view.qtpl:10
		}
//line history/view.qtpl:10
		qw422016.N().S(`<a href="/hypha/`)
//line history/view.qtpl:11
		qw422016.E().S(hyphaName)
//line history/view.qtpl:11
		qw422016.N().S(`">`)
//line history/view.qtpl:11
		qw422016.E().S(hyphaName)
//line history/view.qtpl:11
		qw422016.N().S(`</a>`)
//line history/view.qtpl:12
	}
//line history/view.qtpl:13
	qw422016.N().S(`
`)
//line history/view.qtpl:14
}

//line history/view.qtpl:14
func (rev Revision) WriteHyphaeLinksHTML(qq422016 qtio422016.Writer) {
//line history/view.qtpl:14
	qw422016 := qt422016.AcquireWriter(qq422016)
//line history/view.qtpl:14
	rev.StreamHyphaeLinksHTML(qw422016)
//line history/view.qtpl:14
	qt422016.ReleaseWriter(qw422016)
//line history/view.qtpl:14
}

//line history/view.qtpl:14
func (rev Revision) HyphaeLinksHTML() string {
//line history/view.qtpl:14
	qb422016 := qt422016.AcquireByteBuffer()
//line history/view.qtpl:14
	rev.WriteHyphaeLinksHTML(qb422016)
//line history/view.qtpl:14
	qs422016 := string(qb422016.B)
//line history/view.qtpl:14
	qt422016.ReleaseByteBuffer(qb422016)
//line history/view.qtpl:14
	return qs422016
//line history/view.qtpl:14
}

// HyphaeDiffsHTML returns a comma-separated list of diffs links of current revision for every affected file as HTML string.

//line history/view.qtpl:18
func (rev Revision) StreamHyphaeDiffsHTML(qw422016 *qt422016.Writer) {
//line history/view.qtpl:18
	qw422016.N().S(`
    `)
//line history/view.qtpl:19
	entries := rev.hyphaeAffected()

//line history/view.qtpl:19
	qw422016.N().S(`
`)
//line history/view.qtpl:21
	if len(entries) == 1 {
//line history/view.qtpl:21
		qw422016.N().S(`<a href="/primitive-diff/`)
//line history/view.qtpl:22
		qw422016.E().S(rev.Hash)
//line history/view.qtpl:22
		qw422016.N().S(`/`)
//line history/view.qtpl:22
		qw422016.E().S(entries[0])
//line history/view.qtpl:22
		qw422016.N().S(`">`)
//line history/view.qtpl:22
		qw422016.E().S(rev.Hash)
//line history/view.qtpl:22
		qw422016.N().S(`</a>`)
//line history/view.qtpl:23
	} else {
//line history/view.qtpl:24
		for i, hyphaName := range entries {
//line history/view.qtpl:25
			if i > 0 {
//line history/view.qtpl:25
				qw422016.N().S(`<span aria-hidden="true">, </span>`)
//line history/view.qtpl:27
			}
//line history/view.qtpl:27
			qw422016.N().S(`<a href="/primitive-diff/`)
//line history/view.qtpl:28
			qw422016.E().S(rev.Hash)
//line history/view.qtpl:28
			qw422016.N().S(`/`)
//line history/view.qtpl:28
			qw422016.E().S(hyphaName)
//line history/view.qtpl:28
			qw422016.N().S(`">`)
//line history/view.qtpl:29
			if i == 0 {
//line history/view.qtpl:30
				qw422016.E().S(rev.Hash)
//line history/view.qtpl:30
				qw422016.N().S(`&nbsp;`)
//line history/view.qtpl:31
			}
//line history/view.qtpl:32
			qw422016.E().S(hyphaName)
//line history/view.qtpl:32
			qw422016.N().S(`</a>`)
//line history/view.qtpl:33
		}
//line history/view.qtpl:34
	}
//line history/view.qtpl:35
	qw422016.N().S(`
`)
//line history/view.qtpl:36
}

//line history/view.qtpl:36
func (rev Revision) WriteHyphaeDiffsHTML(qq422016 qtio422016.Writer) {
//line history/view.qtpl:36
	qw422016 := qt422016.AcquireWriter(qq422016)
//line history/view.qtpl:36
	rev.StreamHyphaeDiffsHTML(qw422016)
//line history/view.qtpl:36
	qt422016.ReleaseWriter(qw422016)
//line history/view.qtpl:36
}

//line history/view.qtpl:36
func (rev Revision) HyphaeDiffsHTML() string {
//line history/view.qtpl:36
	qb422016 := qt422016.AcquireByteBuffer()
//line history/view.qtpl:36
	rev.WriteHyphaeDiffsHTML(qb422016)
//line history/view.qtpl:36
	qs422016 := string(qb422016.B)
//line history/view.qtpl:36
	qt422016.ReleaseByteBuffer(qb422016)
//line history/view.qtpl:36
	return qs422016
//line history/view.qtpl:36
}

// descriptionForFeed generates a good enough HTML contents for a web feed.

//line history/view.qtpl:39
func (rev *Revision) streamdescriptionForFeed(qw422016 *qt422016.Writer) {
//line history/view.qtpl:39
	qw422016.N().S(`
<p><b>`)
//line history/view.qtpl:40
	qw422016.E().S(rev.Message)
//line history/view.qtpl:40
	qw422016.N().S(`</b> (by `)
//line history/view.qtpl:40
	qw422016.E().S(rev.Username)
//line history/view.qtpl:40
	qw422016.N().S(` at `)
//line history/view.qtpl:40
	qw422016.E().S(rev.TimeString())
//line history/view.qtpl:40
	qw422016.N().S(`)</p>
<p>Hyphae affected: `)
//line history/view.qtpl:41
	rev.StreamHyphaeLinksHTML(qw422016)
//line history/view.qtpl:41
	qw422016.N().S(`</p>
<pre><code>`)
//line history/view.qtpl:42
	qw422016.E().S(rev.textDiff())
//line history/view.qtpl:42
	qw422016.N().S(`</code></pre>
`)
//line history/view.qtpl:43
}

//line history/view.qtpl:43
func (rev *Revision) writedescriptionForFeed(qq422016 qtio422016.Writer) {
//line history/view.qtpl:43
	qw422016 := qt422016.AcquireWriter(qq422016)
//line history/view.qtpl:43
	rev.streamdescriptionForFeed(qw422016)
//line history/view.qtpl:43
	qt422016.ReleaseWriter(qw422016)
//line history/view.qtpl:43
}

//line history/view.qtpl:43
func (rev *Revision) descriptionForFeed() string {
//line history/view.qtpl:43
	qb422016 := qt422016.AcquireByteBuffer()
//line history/view.qtpl:43
	rev.writedescriptionForFeed(qb422016)
//line history/view.qtpl:43
	qs422016 := string(qb422016.B)
//line history/view.qtpl:43
	qt422016.ReleaseByteBuffer(qb422016)
//line history/view.qtpl:43
	return qs422016
//line history/view.qtpl:43
}

// WithRevisions returns an html representation of `revs` that is meant to be inserted in a history page.

//line history/view.qtpl:46
func StreamWithRevisions(qw422016 *qt422016.Writer, hyphaName string, revs []Revision) {
//line history/view.qtpl:46
	qw422016.N().S(`
`)
//line history/view.qtpl:47
	for _, grp := range groupRevisionsByMonth(revs) {
//line history/view.qtpl:47
		qw422016.N().S(`
	`)
//line history/view.qtpl:49
		currentYear := grp[0].Time.Year()
		currentMonth := grp[0].Time.Month()
		sectionId := fmt.Sprintf("%04d-%02d", currentYear, currentMonth)

//line history/view.qtpl:52
		qw422016.N().S(`
<section class="history__month">
	<a href="#`)
//line history/view.qtpl:54
		qw422016.E().S(sectionId)
//line history/view.qtpl:54
		qw422016.N().S(`" class="history__month-anchor">
		<h2 id="`)
//line history/view.qtpl:55
		qw422016.E().S(sectionId)
//line history/view.qtpl:55
		qw422016.N().S(`" class="history__month-title">`)
//line history/view.qtpl:55
		qw422016.N().D(currentYear)
//line history/view.qtpl:55
		qw422016.N().S(` `)
//line history/view.qtpl:55
		qw422016.E().S(currentMonth.String())
//line history/view.qtpl:55
		qw422016.N().S(`</h2>
	</a>
	<ul class="history__entries">
        `)
//line history/view.qtpl:58
		for _, rev := range grp {
//line history/view.qtpl:58
			qw422016.N().S(`
            <li class="history__entry">
            	<a class="history-entry" href="/rev/`)
//line history/view.qtpl:60
			qw422016.E().S(rev.Hash)
//line history/view.qtpl:60
			qw422016.N().S(`/`)
//line history/view.qtpl:60
			qw422016.E().S(hyphaName)
//line history/view.qtpl:60
			qw422016.N().S(`">
                    <time class="history-entry__time">`)
//line history/view.qtpl:61
			qw422016.E().S(rev.timeToDisplay())
//line history/view.qtpl:61
			qw422016.N().S(`</time>
                </a>
            	<span class="history-entry__hash"><a href="/primitive-diff/`)
//line history/view.qtpl:63
			qw422016.E().S(rev.Hash)
//line history/view.qtpl:63
			qw422016.N().S(`/`)
//line history/view.qtpl:63
			qw422016.E().S(hyphaName)
//line history/view.qtpl:63
			qw422016.N().S(`">`)
//line history/view.qtpl:63
			qw422016.E().S(rev.Hash)
//line history/view.qtpl:63
			qw422016.N().S(`</a></span>
            	<span class="history-entry__msg">`)
//line history/view.qtpl:64
			qw422016.E().S(rev.Message)
//line history/view.qtpl:64
			qw422016.N().S(`</span>
            	`)
//line history/view.qtpl:65
			if rev.Username != "anon" {
//line history/view.qtpl:65
				qw422016.N().S(`
                    <span class="history-entry__author">by <a href="/hypha/`)
//line history/view.qtpl:66
				qw422016.E().S(cfg.UserHypha)
//line history/view.qtpl:66
				qw422016.N().S(`/`)
//line history/view.qtpl:66
				qw422016.E().S(rev.Username)
//line history/view.qtpl:66
				qw422016.N().S(`" rel="author">`)
//line history/view.qtpl:66
				qw422016.E().S(rev.Username)
//line history/view.qtpl:66
				qw422016.N().S(`</a></span>
                `)
//line history/view.qtpl:67
			}
//line history/view.qtpl:67
			qw422016.N().S(`
            </li>
        `)
//line history/view.qtpl:69
		}
//line history/view.qtpl:69
		qw422016.N().S(`
	</ul>
</section>
`)
//line history/view.qtpl:72
	}
//line history/view.qtpl:72
	qw422016.N().S(`
`)
//line history/view.qtpl:73
}

//line history/view.qtpl:73
func WriteWithRevisions(qq422016 qtio422016.Writer, hyphaName string, revs []Revision) {
//line history/view.qtpl:73
	qw422016 := qt422016.AcquireWriter(qq422016)
//line history/view.qtpl:73
	StreamWithRevisions(qw422016, hyphaName, revs)
//line history/view.qtpl:73
	qt422016.ReleaseWriter(qw422016)
//line history/view.qtpl:73
}

//line history/view.qtpl:73
func WithRevisions(hyphaName string, revs []Revision) string {
//line history/view.qtpl:73
	qb422016 := qt422016.AcquireByteBuffer()
//line history/view.qtpl:73
	WriteWithRevisions(qb422016, hyphaName, revs)
//line history/view.qtpl:73
	qs422016 := string(qb422016.B)
//line history/view.qtpl:73
	qt422016.ReleaseByteBuffer(qb422016)
//line history/view.qtpl:73
	return qs422016
//line history/view.qtpl:73
}
