// Code generated by qtc from "auth.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/auth.qtpl:1
package views

//line views/auth.qtpl:1
import "net/http"

//line views/auth.qtpl:2
import "github.com/bouncepaw/mycorrhiza/user"

//line views/auth.qtpl:3
import "github.com/bouncepaw/mycorrhiza/cfg"

//line views/auth.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/auth.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/auth.qtpl:5
func StreamRegisterHTML(qw422016 *qt422016.Writer, rq *http.Request) {
//line views/auth.qtpl:5
	qw422016.N().S(`
<div class="layout">
<main class="main-width">
	<section>
	`)
//line views/auth.qtpl:9
	if cfg.UseRegistration {
//line views/auth.qtpl:9
		qw422016.N().S(`
		<form class="modal" method="post" action="/register?`)
//line views/auth.qtpl:10
		qw422016.E().S(rq.URL.RawQuery)
//line views/auth.qtpl:10
		qw422016.N().S(`" id="register-form" enctype="multipart/form-data" autocomplete="off">
			<fieldset class="modal__fieldset">
				<legend class="modal__title">Register to `)
//line views/auth.qtpl:12
		qw422016.E().S(cfg.WikiName)
//line views/auth.qtpl:12
		qw422016.N().S(`</legend>

				<label for="register-form__username">Username</label>
				<br>
				<input type="text" required autofocus id="login-form__username" name="username">
				<br>
				<label for="login-form__password">Password</label>
				<br>
				<input type="password" required name="password">
				<p>The server stores your password in an encrypted form; even administrators cannot read it.</p>
				<p>By submitting this form you give this wiki a permission to store cookies in your browser. It lets the engine associate your edits with you. You will stay logged in until you log out.</p>
				<input class="btn" type="submit" value="Register">
				<a class="btn btn_weak" href="/`)
//line views/auth.qtpl:24
		qw422016.E().S(rq.URL.RawQuery)
//line views/auth.qtpl:24
		qw422016.N().S(`">Cancel</a>
			</fieldset>
		</form>
	`)
//line views/auth.qtpl:27
	} else if cfg.UseFixedAuth {
//line views/auth.qtpl:27
		qw422016.N().S(`
		<p>Administrators have forbidden registration for this wiki. Administrators can make an account for you by hand; contact them.</p>
		<p><a href="/`)
//line views/auth.qtpl:29
		qw422016.E().S(rq.URL.RawQuery)
//line views/auth.qtpl:29
		qw422016.N().S(`">← Go back</a></p>
	`)
//line views/auth.qtpl:30
	} else {
//line views/auth.qtpl:30
		qw422016.N().S(`
		<p>Administrators of this wiki have not configured any authorization method. You can make edits anonymously.</p>
		<p><a href="/`)
//line views/auth.qtpl:32
		qw422016.E().S(rq.URL.RawQuery)
//line views/auth.qtpl:32
		qw422016.N().S(`">← Go back</a></p>
	`)
//line views/auth.qtpl:33
	}
//line views/auth.qtpl:33
	qw422016.N().S(`
	</section>
</main>
</div>
`)
//line views/auth.qtpl:37
}

//line views/auth.qtpl:37
func WriteRegisterHTML(qq422016 qtio422016.Writer, rq *http.Request) {
//line views/auth.qtpl:37
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/auth.qtpl:37
	StreamRegisterHTML(qw422016, rq)
//line views/auth.qtpl:37
	qt422016.ReleaseWriter(qw422016)
//line views/auth.qtpl:37
}

//line views/auth.qtpl:37
func RegisterHTML(rq *http.Request) string {
//line views/auth.qtpl:37
	qb422016 := qt422016.AcquireByteBuffer()
//line views/auth.qtpl:37
	WriteRegisterHTML(qb422016, rq)
//line views/auth.qtpl:37
	qs422016 := string(qb422016.B)
//line views/auth.qtpl:37
	qt422016.ReleaseByteBuffer(qb422016)
//line views/auth.qtpl:37
	return qs422016
//line views/auth.qtpl:37
}

//line views/auth.qtpl:39
func StreamLoginHTML(qw422016 *qt422016.Writer) {
//line views/auth.qtpl:39
	qw422016.N().S(`
<div class="layout">
<main class="main-width">
	<section>
	`)
//line views/auth.qtpl:43
	if user.AuthUsed {
//line views/auth.qtpl:43
		qw422016.N().S(`
		<form class="modal" method="post" action="/login-data" id="login-form" enctype="multipart/form-data" autocomplete="on">
			<fieldset class="modal__fieldset">
				<legend class="modal__title">Log in to `)
//line views/auth.qtpl:46
		qw422016.E().S(cfg.WikiName)
//line views/auth.qtpl:46
		qw422016.N().S(`</legend>
				<label for="login-form__username">Username</label>
				<br>
				<input type="text" required autofocus id="login-form__username" name="username" autocomplete="username">
				<br>
				<label for="login-form__password">Password</label>
				<br>
				<input type="password" required name="password" autocomplete="current-password">
				<p>By submitting this form you give this wiki a permission to store cookies in your browser. It lets the engine associate your edits with you. You will stay logged in until you log out.</p>
				<input class="btn" type="submit" value="Log in">
				<a class="btn btn_weak" href="/">Cancel</a>
			</fieldset>
		</form>
	`)
//line views/auth.qtpl:59
	} else {
//line views/auth.qtpl:59
		qw422016.N().S(`
		<p>Administrators of this wiki have not configured any authorization method. You can make edits anonymously.</p>
		<p><a class="modal__cancel" href="/">← Go home</a></p>
	`)
//line views/auth.qtpl:62
	}
//line views/auth.qtpl:62
	qw422016.N().S(`
	</section>
</main>
</div>
`)
//line views/auth.qtpl:66
}

//line views/auth.qtpl:66
func WriteLoginHTML(qq422016 qtio422016.Writer) {
//line views/auth.qtpl:66
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/auth.qtpl:66
	StreamLoginHTML(qw422016)
//line views/auth.qtpl:66
	qt422016.ReleaseWriter(qw422016)
//line views/auth.qtpl:66
}

//line views/auth.qtpl:66
func LoginHTML() string {
//line views/auth.qtpl:66
	qb422016 := qt422016.AcquireByteBuffer()
//line views/auth.qtpl:66
	WriteLoginHTML(qb422016)
//line views/auth.qtpl:66
	qs422016 := string(qb422016.B)
//line views/auth.qtpl:66
	qt422016.ReleaseByteBuffer(qb422016)
//line views/auth.qtpl:66
	return qs422016
//line views/auth.qtpl:66
}

//line views/auth.qtpl:68
func StreamLoginErrorHTML(qw422016 *qt422016.Writer, err string) {
//line views/auth.qtpl:68
	qw422016.N().S(`
<div class="layout">
<main class="main-width">
	<section>
	`)
//line views/auth.qtpl:72
	switch err {
//line views/auth.qtpl:73
	case "unknown username":
//line views/auth.qtpl:73
		qw422016.N().S(`
		<p class="error">Unknown username.</p>
	`)
//line views/auth.qtpl:75
	case "wrong password":
//line views/auth.qtpl:75
		qw422016.N().S(`
		<p class="error">Wrong password.</p>
	`)
//line views/auth.qtpl:77
	default:
//line views/auth.qtpl:77
		qw422016.N().S(`
		<p class="error">`)
//line views/auth.qtpl:78
		qw422016.E().S(err)
//line views/auth.qtpl:78
		qw422016.N().S(`</p>
	`)
//line views/auth.qtpl:79
	}
//line views/auth.qtpl:79
	qw422016.N().S(`
		<p><a href="/login">← Try again</a></p>
	</section>
</main>
</div>
`)
//line views/auth.qtpl:84
}

//line views/auth.qtpl:84
func WriteLoginErrorHTML(qq422016 qtio422016.Writer, err string) {
//line views/auth.qtpl:84
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/auth.qtpl:84
	StreamLoginErrorHTML(qw422016, err)
//line views/auth.qtpl:84
	qt422016.ReleaseWriter(qw422016)
//line views/auth.qtpl:84
}

//line views/auth.qtpl:84
func LoginErrorHTML(err string) string {
//line views/auth.qtpl:84
	qb422016 := qt422016.AcquireByteBuffer()
//line views/auth.qtpl:84
	WriteLoginErrorHTML(qb422016, err)
//line views/auth.qtpl:84
	qs422016 := string(qb422016.B)
//line views/auth.qtpl:84
	qt422016.ReleaseByteBuffer(qb422016)
//line views/auth.qtpl:84
	return qs422016
//line views/auth.qtpl:84
}

//line views/auth.qtpl:86
func StreamLogoutHTML(qw422016 *qt422016.Writer, can bool) {
//line views/auth.qtpl:86
	qw422016.N().S(`
<div class="layout">
<main class="main-width">
	<section>
	`)
//line views/auth.qtpl:90
	if can {
//line views/auth.qtpl:90
		qw422016.N().S(`
		<h1>Log out?</h1>
		<p><a href="/logout-confirm"><strong>Confirm</strong></a></p>
		<p><a href="/">Cancel</a></p>
	`)
//line views/auth.qtpl:94
	} else {
//line views/auth.qtpl:94
		qw422016.N().S(`
		<p>You cannot log out because you are not logged in.</p>
		<p><a href="/login">Login</a></p>
		<p><a href="/login">← Home</a></p>
	`)
//line views/auth.qtpl:98
	}
//line views/auth.qtpl:98
	qw422016.N().S(`
	</section>
</main>
</div>
`)
//line views/auth.qtpl:102
}

//line views/auth.qtpl:102
func WriteLogoutHTML(qq422016 qtio422016.Writer, can bool) {
//line views/auth.qtpl:102
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/auth.qtpl:102
	StreamLogoutHTML(qw422016, can)
//line views/auth.qtpl:102
	qt422016.ReleaseWriter(qw422016)
//line views/auth.qtpl:102
}

//line views/auth.qtpl:102
func LogoutHTML(can bool) string {
//line views/auth.qtpl:102
	qb422016 := qt422016.AcquireByteBuffer()
//line views/auth.qtpl:102
	WriteLogoutHTML(qb422016, can)
//line views/auth.qtpl:102
	qs422016 := string(qb422016.B)
//line views/auth.qtpl:102
	qt422016.ReleaseByteBuffer(qb422016)
//line views/auth.qtpl:102
	return qs422016
//line views/auth.qtpl:102
}
