// Code generated by qtc from "http_delete.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// This dialog is to be shown to a user when they try to delete a hypha.

//line templates/http_delete.qtpl:2
package templates

//line templates/http_delete.qtpl:2
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/http_delete.qtpl:2
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/http_delete.qtpl:2
func StreamDeleteAskHTML(qw422016 *qt422016.Writer, hyphaName string, isOld bool) {
//line templates/http_delete.qtpl:2
	qw422016.N().S(`
<main>
`)
//line templates/http_delete.qtpl:4
	streamnavHTML(qw422016, hyphaName, "delete-ask")
//line templates/http_delete.qtpl:4
	qw422016.N().S(`
`)
//line templates/http_delete.qtpl:5
	if isOld {
//line templates/http_delete.qtpl:5
		qw422016.N().S(`
	<section>
		<h1>Delete `)
//line templates/http_delete.qtpl:7
		qw422016.E().S(hyphaName)
//line templates/http_delete.qtpl:7
		qw422016.N().S(`?</h1>
		<p>Do you really want to delete hypha <em>`)
//line templates/http_delete.qtpl:8
		qw422016.E().S(hyphaName)
//line templates/http_delete.qtpl:8
		qw422016.N().S(`</em>?</p>
		<p>In this version of MycorrhizaWiki you cannot undelete a deleted hypha but the history can still be accessed.</p>
		<p><a href="/delete-confirm/`)
//line templates/http_delete.qtpl:10
		qw422016.E().S(hyphaName)
//line templates/http_delete.qtpl:10
		qw422016.N().S(`"><strong>Confirm</strong></a></p>
		<p><a href="/page/`)
//line templates/http_delete.qtpl:11
		qw422016.E().S(hyphaName)
//line templates/http_delete.qtpl:11
		qw422016.N().S(`">Cancel</a></p>
	</section>
`)
//line templates/http_delete.qtpl:13
	} else {
//line templates/http_delete.qtpl:13
		qw422016.N().S(`
	`)
//line templates/http_delete.qtpl:14
		streamcannotDeleteDueToNonExistence(qw422016, hyphaName)
//line templates/http_delete.qtpl:14
		qw422016.N().S(`
`)
//line templates/http_delete.qtpl:15
	}
//line templates/http_delete.qtpl:15
	qw422016.N().S(`
</main>
`)
//line templates/http_delete.qtpl:17
}

//line templates/http_delete.qtpl:17
func WriteDeleteAskHTML(qq422016 qtio422016.Writer, hyphaName string, isOld bool) {
//line templates/http_delete.qtpl:17
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/http_delete.qtpl:17
	StreamDeleteAskHTML(qw422016, hyphaName, isOld)
//line templates/http_delete.qtpl:17
	qt422016.ReleaseWriter(qw422016)
//line templates/http_delete.qtpl:17
}

//line templates/http_delete.qtpl:17
func DeleteAskHTML(hyphaName string, isOld bool) string {
//line templates/http_delete.qtpl:17
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/http_delete.qtpl:17
	WriteDeleteAskHTML(qb422016, hyphaName, isOld)
//line templates/http_delete.qtpl:17
	qs422016 := string(qb422016.B)
//line templates/http_delete.qtpl:17
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/http_delete.qtpl:17
	return qs422016
//line templates/http_delete.qtpl:17
}

//line templates/http_delete.qtpl:19
func streamcannotDeleteDueToNonExistence(qw422016 *qt422016.Writer, hyphaName string) {
//line templates/http_delete.qtpl:19
	qw422016.N().S(`
	<section>
		<h1>Cannot delete `)
//line templates/http_delete.qtpl:21
	qw422016.E().S(hyphaName)
//line templates/http_delete.qtpl:21
	qw422016.N().S(`</h1>
		<p>This hypha does not exist.</p>
		<p><a href="/page/`)
//line templates/http_delete.qtpl:23
	qw422016.E().S(hyphaName)
//line templates/http_delete.qtpl:23
	qw422016.N().S(`">Go back</a></p>
	</section>
`)
//line templates/http_delete.qtpl:25
}

//line templates/http_delete.qtpl:25
func writecannotDeleteDueToNonExistence(qq422016 qtio422016.Writer, hyphaName string) {
//line templates/http_delete.qtpl:25
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/http_delete.qtpl:25
	streamcannotDeleteDueToNonExistence(qw422016, hyphaName)
//line templates/http_delete.qtpl:25
	qt422016.ReleaseWriter(qw422016)
//line templates/http_delete.qtpl:25
}

//line templates/http_delete.qtpl:25
func cannotDeleteDueToNonExistence(hyphaName string) string {
//line templates/http_delete.qtpl:25
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/http_delete.qtpl:25
	writecannotDeleteDueToNonExistence(qb422016, hyphaName)
//line templates/http_delete.qtpl:25
	qs422016 := string(qb422016.B)
//line templates/http_delete.qtpl:25
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/http_delete.qtpl:25
	return qs422016
//line templates/http_delete.qtpl:25
}