// Code generated by qtc from "auth.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line auth/auth.qtpl:1
package auth

//line auth/auth.qtpl:1
import "net/http"

//line auth/auth.qtpl:2
import "github.com/bouncepaw/mycorrhiza/internal/cfg"

//line auth/auth.qtpl:3
import "github.com/bouncepaw/mycorrhiza/l18n"

//line auth/auth.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line auth/auth.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line auth/auth.qtpl:5
func StreamRegister(qw422016 *qt422016.Writer, rq *http.Request) {
//line auth/auth.qtpl:5
	qw422016.N().S(`
`)
//line auth/auth.qtpl:7
	lc := l18n.FromRequest(rq)

//line auth/auth.qtpl:8
	qw422016.N().S(`
<main class="main-width">
	<section>
	`)
//line auth/auth.qtpl:11
	if cfg.AllowRegistration {
//line auth/auth.qtpl:11
		qw422016.N().S(`
		<form class="modal" method="post" action="/register?`)
//line auth/auth.qtpl:12
		qw422016.E().S(rq.URL.RawQuery)
//line auth/auth.qtpl:12
		qw422016.N().S(`" id="register-form" enctype="multipart/form-data" autocomplete="off">
			<fieldset class="modal__fieldset">
				<legend class="modal__title">`)
//line auth/auth.qtpl:14
		qw422016.E().S(lc.Get("auth.register_header", &l18n.Replacements{"name": cfg.WikiName}))
//line auth/auth.qtpl:14
		qw422016.N().S(`</legend>

				<label for="register-form__username">`)
//line auth/auth.qtpl:16
		qw422016.E().S(lc.Get("auth.username"))
//line auth/auth.qtpl:16
		qw422016.N().S(`</label>
				<br>
				<input type="text" required autofocus id="login-form__username" name="username">
				<br>
				<label for="login-form__password">`)
//line auth/auth.qtpl:20
		qw422016.E().S(lc.Get("auth.password"))
//line auth/auth.qtpl:20
		qw422016.N().S(`</label>
				<br>
				<input type="password" required name="password">
				<p>`)
//line auth/auth.qtpl:23
		qw422016.E().S(lc.Get("auth.password_tip"))
//line auth/auth.qtpl:23
		qw422016.N().S(`</p>
				<p>`)
//line auth/auth.qtpl:24
		qw422016.E().S(lc.Get("auth.cookie_tip"))
//line auth/auth.qtpl:24
		qw422016.N().S(`</p>
				<button class="btn" type="submit" value="Register">`)
//line auth/auth.qtpl:25
		qw422016.E().S(lc.Get("auth.register_button"))
//line auth/auth.qtpl:25
		qw422016.N().S(`</button>
				<a class="btn btn_weak" href="/`)
//line auth/auth.qtpl:26
		qw422016.E().S(rq.URL.RawQuery)
//line auth/auth.qtpl:26
		qw422016.N().S(`">`)
//line auth/auth.qtpl:26
		qw422016.E().S(lc.Get("ui.cancel"))
//line auth/auth.qtpl:26
		qw422016.N().S(`</a>
			</fieldset>
		</form>
		`)
//line auth/auth.qtpl:29
		streamtelegramWidget(qw422016, lc)
//line auth/auth.qtpl:29
		qw422016.N().S(`
	`)
//line auth/auth.qtpl:30
	} else if cfg.UseAuth {
//line auth/auth.qtpl:30
		qw422016.N().S(`
		<p>`)
//line auth/auth.qtpl:31
		qw422016.E().S(lc.Get("auth.noregister"))
//line auth/auth.qtpl:31
		qw422016.N().S(`</p>
		<p><a href="/`)
//line auth/auth.qtpl:32
		qw422016.E().S(rq.URL.RawQuery)
//line auth/auth.qtpl:32
		qw422016.N().S(`">← `)
//line auth/auth.qtpl:32
		qw422016.E().S(lc.Get("auth.go_back"))
//line auth/auth.qtpl:32
		qw422016.N().S(`</a></p>
	`)
//line auth/auth.qtpl:33
	} else {
//line auth/auth.qtpl:33
		qw422016.N().S(`
		<p>`)
//line auth/auth.qtpl:34
		qw422016.E().S(lc.Get("auth.noauth"))
//line auth/auth.qtpl:34
		qw422016.N().S(`</p>
		<p><a href="/`)
//line auth/auth.qtpl:35
		qw422016.E().S(rq.URL.RawQuery)
//line auth/auth.qtpl:35
		qw422016.N().S(`">← `)
//line auth/auth.qtpl:35
		qw422016.E().S(lc.Get("auth.go_back"))
//line auth/auth.qtpl:35
		qw422016.N().S(`</a></p>
	`)
//line auth/auth.qtpl:36
	}
//line auth/auth.qtpl:36
	qw422016.N().S(`
	</section>
</main>
`)
//line auth/auth.qtpl:39
}

//line auth/auth.qtpl:39
func WriteRegister(qq422016 qtio422016.Writer, rq *http.Request) {
//line auth/auth.qtpl:39
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth/auth.qtpl:39
	StreamRegister(qw422016, rq)
//line auth/auth.qtpl:39
	qt422016.ReleaseWriter(qw422016)
//line auth/auth.qtpl:39
}

//line auth/auth.qtpl:39
func Register(rq *http.Request) string {
//line auth/auth.qtpl:39
	qb422016 := qt422016.AcquireByteBuffer()
//line auth/auth.qtpl:39
	WriteRegister(qb422016, rq)
//line auth/auth.qtpl:39
	qs422016 := string(qb422016.B)
//line auth/auth.qtpl:39
	qt422016.ReleaseByteBuffer(qb422016)
//line auth/auth.qtpl:39
	return qs422016
//line auth/auth.qtpl:39
}

//line auth/auth.qtpl:41
func StreamLogin(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line auth/auth.qtpl:41
	qw422016.N().S(`
<main class="main-width">
	<section>
	`)
//line auth/auth.qtpl:44
	if cfg.UseAuth {
//line auth/auth.qtpl:44
		qw422016.N().S(`
		<form class="modal" method="post" action="/login" id="login-form" enctype="multipart/form-data" autocomplete="on">
			<fieldset class="modal__fieldset">
				<legend class="modal__title">`)
//line auth/auth.qtpl:47
		qw422016.E().S(lc.Get("auth.login_header", &l18n.Replacements{"name": cfg.WikiName}))
//line auth/auth.qtpl:47
		qw422016.N().S(`</legend>
				<label for="login-form__username">`)
//line auth/auth.qtpl:48
		qw422016.E().S(lc.Get("auth.username"))
//line auth/auth.qtpl:48
		qw422016.N().S(`</label>
				<br>
				<input type="text" required autofocus id="login-form__username" name="username" autocomplete="username">
				<br>
				<label for="login-form__password">`)
//line auth/auth.qtpl:52
		qw422016.E().S(lc.Get("auth.password"))
//line auth/auth.qtpl:52
		qw422016.N().S(`</label>
				<br>
				<input type="password" required name="password" autocomplete="current-password">
				<p>`)
//line auth/auth.qtpl:55
		qw422016.E().S(lc.Get("auth.cookie_tip"))
//line auth/auth.qtpl:55
		qw422016.N().S(`</p>
				<button class="btn" type="submit" value="Log in">`)
//line auth/auth.qtpl:56
		qw422016.E().S(lc.Get("auth.login_button"))
//line auth/auth.qtpl:56
		qw422016.N().S(`</button>
				<a class="btn btn_weak" href="/">`)
//line auth/auth.qtpl:57
		qw422016.E().S(lc.Get("ui.cancel"))
//line auth/auth.qtpl:57
		qw422016.N().S(`</a>
			</fieldset>
		</form>
		`)
//line auth/auth.qtpl:60
		streamtelegramWidget(qw422016, lc)
//line auth/auth.qtpl:60
		qw422016.N().S(`
	`)
//line auth/auth.qtpl:61
	} else {
//line auth/auth.qtpl:61
		qw422016.N().S(`
		<p>`)
//line auth/auth.qtpl:62
		qw422016.E().S(lc.Get("auth.noauth"))
//line auth/auth.qtpl:62
		qw422016.N().S(`</p>
		<p><a class="btn btn_weak" href="/">← `)
//line auth/auth.qtpl:63
		qw422016.E().S(lc.Get("auth.go_home"))
//line auth/auth.qtpl:63
		qw422016.N().S(`</a></p>
	`)
//line auth/auth.qtpl:64
	}
//line auth/auth.qtpl:64
	qw422016.N().S(`
	</section>
</main>
`)
//line auth/auth.qtpl:67
}

//line auth/auth.qtpl:67
func WriteLogin(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line auth/auth.qtpl:67
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth/auth.qtpl:67
	StreamLogin(qw422016, lc)
//line auth/auth.qtpl:67
	qt422016.ReleaseWriter(qw422016)
//line auth/auth.qtpl:67
}

//line auth/auth.qtpl:67
func Login(lc *l18n.Localizer) string {
//line auth/auth.qtpl:67
	qb422016 := qt422016.AcquireByteBuffer()
//line auth/auth.qtpl:67
	WriteLogin(qb422016, lc)
//line auth/auth.qtpl:67
	qs422016 := string(qb422016.B)
//line auth/auth.qtpl:67
	qt422016.ReleaseByteBuffer(qb422016)
//line auth/auth.qtpl:67
	return qs422016
//line auth/auth.qtpl:67
}

// Telegram auth widget was requested by Yogurt. As you can see, we don't offer user administrators control over it. Of course we don't.

//line auth/auth.qtpl:70
func streamtelegramWidget(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line auth/auth.qtpl:70
	qw422016.N().S(`
`)
//line auth/auth.qtpl:71
	if cfg.TelegramEnabled {
//line auth/auth.qtpl:71
		qw422016.N().S(`
<p class="telegram-notice">`)
//line auth/auth.qtpl:72
		qw422016.E().S(lc.Get("auth.telegram_tip"))
//line auth/auth.qtpl:72
		qw422016.N().S(`</p>
<script async src="https://telegram.org/js/telegram-widget.js?15" data-telegram-login="`)
//line auth/auth.qtpl:73
		qw422016.E().S(cfg.TelegramBotName)
//line auth/auth.qtpl:73
		qw422016.N().S(`" data-size="medium" data-userpic="false" data-auth-url="`)
//line auth/auth.qtpl:73
		qw422016.E().S(cfg.URL)
//line auth/auth.qtpl:73
		qw422016.N().S(`/telegram-login"></script>
`)
//line auth/auth.qtpl:74
	}
//line auth/auth.qtpl:74
	qw422016.N().S(`
`)
//line auth/auth.qtpl:75
}

//line auth/auth.qtpl:75
func writetelegramWidget(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line auth/auth.qtpl:75
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth/auth.qtpl:75
	streamtelegramWidget(qw422016, lc)
//line auth/auth.qtpl:75
	qt422016.ReleaseWriter(qw422016)
//line auth/auth.qtpl:75
}

//line auth/auth.qtpl:75
func telegramWidget(lc *l18n.Localizer) string {
//line auth/auth.qtpl:75
	qb422016 := qt422016.AcquireByteBuffer()
//line auth/auth.qtpl:75
	writetelegramWidget(qb422016, lc)
//line auth/auth.qtpl:75
	qs422016 := string(qb422016.B)
//line auth/auth.qtpl:75
	qt422016.ReleaseByteBuffer(qb422016)
//line auth/auth.qtpl:75
	return qs422016
//line auth/auth.qtpl:75
}

//line auth/auth.qtpl:77
func StreamLoginError(qw422016 *qt422016.Writer, err string, lc *l18n.Localizer) {
//line auth/auth.qtpl:77
	qw422016.N().S(`
<main class="main-width">
	<section>
	`)
//line auth/auth.qtpl:80
	switch err {
//line auth/auth.qtpl:81
	case "unknown username":
//line auth/auth.qtpl:81
		qw422016.N().S(`
		<p class="error">`)
//line auth/auth.qtpl:82
		qw422016.E().S(lc.Get("auth.error_username"))
//line auth/auth.qtpl:82
		qw422016.N().S(`</p>
	`)
//line auth/auth.qtpl:83
	case "wrong password":
//line auth/auth.qtpl:83
		qw422016.N().S(`
		<p class="error">`)
//line auth/auth.qtpl:84
		qw422016.E().S(lc.Get("auth.error_password"))
//line auth/auth.qtpl:84
		qw422016.N().S(`</p>
	`)
//line auth/auth.qtpl:85
	default:
//line auth/auth.qtpl:85
		qw422016.N().S(`
		<p class="error">`)
//line auth/auth.qtpl:86
		qw422016.E().S(err)
//line auth/auth.qtpl:86
		qw422016.N().S(`</p>
	`)
//line auth/auth.qtpl:87
	}
//line auth/auth.qtpl:87
	qw422016.N().S(`
		<p><a href="/login">← `)
//line auth/auth.qtpl:88
	qw422016.E().S(lc.Get("auth.try_again"))
//line auth/auth.qtpl:88
	qw422016.N().S(`</a></p>
	</section>
</main>
`)
//line auth/auth.qtpl:91
}

//line auth/auth.qtpl:91
func WriteLoginError(qq422016 qtio422016.Writer, err string, lc *l18n.Localizer) {
//line auth/auth.qtpl:91
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth/auth.qtpl:91
	StreamLoginError(qw422016, err, lc)
//line auth/auth.qtpl:91
	qt422016.ReleaseWriter(qw422016)
//line auth/auth.qtpl:91
}

//line auth/auth.qtpl:91
func LoginError(err string, lc *l18n.Localizer) string {
//line auth/auth.qtpl:91
	qb422016 := qt422016.AcquireByteBuffer()
//line auth/auth.qtpl:91
	WriteLoginError(qb422016, err, lc)
//line auth/auth.qtpl:91
	qs422016 := string(qb422016.B)
//line auth/auth.qtpl:91
	qt422016.ReleaseByteBuffer(qb422016)
//line auth/auth.qtpl:91
	return qs422016
//line auth/auth.qtpl:91
}

//line auth/auth.qtpl:93
func StreamLogout(qw422016 *qt422016.Writer, can bool, lc *l18n.Localizer) {
//line auth/auth.qtpl:93
	qw422016.N().S(`
<main class="main-width">
	<section>
	`)
//line auth/auth.qtpl:96
	if can {
//line auth/auth.qtpl:96
		qw422016.N().S(`
		<h1>`)
//line auth/auth.qtpl:97
		qw422016.E().S(lc.Get("auth.logout_header"))
//line auth/auth.qtpl:97
		qw422016.N().S(`</h1>
		<form method="POST" action="/logout">
			<input class="btn btn_accent" type="submit" value="`)
//line auth/auth.qtpl:99
		qw422016.E().S(lc.Get("auth.logout_button"))
//line auth/auth.qtpl:99
		qw422016.N().S(`"/>
			<a class="btn btn_weak" href="/">`)
//line auth/auth.qtpl:100
		qw422016.E().S(lc.Get("auth.go_home"))
//line auth/auth.qtpl:100
		qw422016.N().S(`</a>
		</form>
	`)
//line auth/auth.qtpl:102
	} else {
//line auth/auth.qtpl:102
		qw422016.N().S(`
		<p>`)
//line auth/auth.qtpl:103
		qw422016.E().S(lc.Get("auth.logout_anon"))
//line auth/auth.qtpl:103
		qw422016.N().S(`</p>
		<p><a href="/login">`)
//line auth/auth.qtpl:104
		qw422016.E().S(lc.Get("auth.login_title"))
//line auth/auth.qtpl:104
		qw422016.N().S(`</a></p>
		<p><a href="/">← `)
//line auth/auth.qtpl:105
		qw422016.E().S(lc.Get("auth.go_home"))
//line auth/auth.qtpl:105
		qw422016.N().S(`</a></p>
	`)
//line auth/auth.qtpl:106
	}
//line auth/auth.qtpl:106
	qw422016.N().S(`
	</section>
</main>
`)
//line auth/auth.qtpl:109
}

//line auth/auth.qtpl:109
func WriteLogout(qq422016 qtio422016.Writer, can bool, lc *l18n.Localizer) {
//line auth/auth.qtpl:109
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth/auth.qtpl:109
	StreamLogout(qw422016, can, lc)
//line auth/auth.qtpl:109
	qt422016.ReleaseWriter(qw422016)
//line auth/auth.qtpl:109
}

//line auth/auth.qtpl:109
func Logout(can bool, lc *l18n.Localizer) string {
//line auth/auth.qtpl:109
	qb422016 := qt422016.AcquireByteBuffer()
//line auth/auth.qtpl:109
	WriteLogout(qb422016, can, lc)
//line auth/auth.qtpl:109
	qs422016 := string(qb422016.B)
//line auth/auth.qtpl:109
	qt422016.ReleaseByteBuffer(qb422016)
//line auth/auth.qtpl:109
	return qs422016
//line auth/auth.qtpl:109
}

//line auth/auth.qtpl:111
func StreamLock(qw422016 *qt422016.Writer, lc *l18n.Localizer) {
//line auth/auth.qtpl:111
	qw422016.N().S(`
<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>🔒 `)
//line auth/auth.qtpl:117
	qw422016.E().S(lc.Get("auth.lock_title"))
//line auth/auth.qtpl:117
	qw422016.N().S(`</title>
	<link rel="shortcut icon" href="/static/favicon.ico">
	<link rel="stylesheet" href="/static/style.css">
</head>
<body>
	<main class="locked-notice main-width">
		<section class="locked-notice__message">
			<p class="locked-notice__lock">🔒</p>
			<h1 class="locked-notice__title">`)
//line auth/auth.qtpl:125
	qw422016.E().S(lc.Get("auth.lock_title"))
//line auth/auth.qtpl:125
	qw422016.N().S(`</h1>
			<form class="locked-notice__login-form" method="post" action="/login" id="login-form" enctype="multipart/form-data" autocomplete="on">
				<label for="login-form__username">`)
//line auth/auth.qtpl:127
	qw422016.E().S(lc.Get("auth.username"))
//line auth/auth.qtpl:127
	qw422016.N().S(`</label>
				<br>
				<input type="text" required autofocus id="login-form__username" name="username" autocomplete="username">
				<br>
				<label for="login-form__password">`)
//line auth/auth.qtpl:131
	qw422016.E().S(lc.Get("auth.password"))
//line auth/auth.qtpl:131
	qw422016.N().S(`</label>
				<br>
				<input type="password" required name="password" autocomplete="current-password">
				<br>
				<button class="btn" type="submit" value="Log in">`)
//line auth/auth.qtpl:135
	qw422016.E().S(lc.Get("auth.login_button"))
//line auth/auth.qtpl:135
	qw422016.N().S(`</button>
			</form>
			`)
//line auth/auth.qtpl:137
	streamtelegramWidget(qw422016, lc)
//line auth/auth.qtpl:137
	qw422016.N().S(`
		</section>
	</main>
</body>
</html>
`)
//line auth/auth.qtpl:142
}

//line auth/auth.qtpl:142
func WriteLock(qq422016 qtio422016.Writer, lc *l18n.Localizer) {
//line auth/auth.qtpl:142
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth/auth.qtpl:142
	StreamLock(qw422016, lc)
//line auth/auth.qtpl:142
	qt422016.ReleaseWriter(qw422016)
//line auth/auth.qtpl:142
}

//line auth/auth.qtpl:142
func Lock(lc *l18n.Localizer) string {
//line auth/auth.qtpl:142
	qb422016 := qt422016.AcquireByteBuffer()
//line auth/auth.qtpl:142
	WriteLock(qb422016, lc)
//line auth/auth.qtpl:142
	qs422016 := string(qb422016.B)
//line auth/auth.qtpl:142
	qt422016.ReleaseByteBuffer(qb422016)
//line auth/auth.qtpl:142
	return qs422016
//line auth/auth.qtpl:142
}
