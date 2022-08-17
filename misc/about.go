package misc

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/version"
	"log"
	"strings"
	"text/template" // sic!
)

type L10nEntry struct {
	_en string
	_ru string
}

func En(v string) L10nEntry {
	return L10nEntry{_en: v}
}

func (e L10nEntry) Ru(v string) L10nEntry {
	e._ru = v
	return e
}

func (e L10nEntry) Get(lang string) string {
	if lang == "ru" && e._ru != "" {
		return e._ru
	}
	return e._en
}

const aboutTemplateString = `
<main class="main-width">
	<section class="about-page">
		<h1>{{ printf (get .L.Title) .Cfg.WikiName }}</h1>
		<dl>
			<dt>{{ get .L.Version }}</dt>
			<dd>{{ .Version }}</dd>
		{{ if .Cfg.UseAuth }}
			<dt>{{ get .L.HomeHypha }}</dt>
			<dd><a href="/">{{ .Cfg.HomeHypha }}</a></dd>

			<dt>{{get .L.Auth}}</dt>
			<dd>{{ get .L.AuthOn }}</dd>
			{{if .Cfg.TelegramEnabled}}<dd>{{get .L.TelegramOn}}</dd>{{end}}

			<dt>{{ get .L.UserCount }}</dt>
			<dd>{{ .UserCount }}</dd>
			{{if .Cfg.RegistrationLimit}}
			<dt>{{get .L.RegistrationLimit}}</dt>
			<dd>{{.RegistrationLimit}}</dd>
			{{end}}

			<dt>{{ get .L.Admins }}</dt>
			{{$cfg := .Cfg}}{{ range $i, $username := .Admins }}
				<dd><a href="/hypha/{{ $cfg.UserHypha }}/{{ $username }}">{{ $username }}</a></dd>
			{{ end }}

		{{ else }}
			<dt>{{get .L.Auth}}</dt>
			<dd>{{ get .L.AuthOff }}</dd>
		{{ end }}
		</dl>
	</section>
</main>`

var aboutData = struct {
	L                 map[string]L10nEntry
	Version           string
	Cfg               map[string]interface{}
	Admins            []string
	UserCount         uint64
	RegistrationLimit uint64
}{
	L: map[string]L10nEntry{
		"Title":             En("About %s").Ru("О %s"),
		"Version":           En("<a href=\"https://mycorrhiza.wiki\">Mycorrhiza Wiki</a> version").Ru("Версия <a href=\"https://mycorrhiza.wiki\">Микоризы</a>"),
		"UserCount":         En("User count").Ru("Число пользователей"),
		"HomeHypha":         En("Home hypha").Ru("Домашняя гифа"),
		"RegistrationLimit": En("RegistrationLimit").Ru("Максимум пользователей"),
		"Admins":            En("Administrators").Ru("Администраторы"),

		"Auth":       En("Authentication").Ru("Аутентификация"),
		"AuthOn":     En("Authentication is on").Ru("Аутентификация включена"),
		"AuthOff":    En("Authentication is off").Ru("Аутентификация не включена"),
		"TelegramOn": En("Telegram authentication is on").Ru("Вход через Телеграм включён"),
	},
}

func AboutHTML(lc *l18n.Localizer) string {
	get := func(e L10nEntry) string {
		return e.Get(lc.Locale)
	}
	temp, err := template.New("about wiki").Funcs(template.FuncMap{"get": get}).Parse(aboutTemplateString)
	if err != nil {
		log.Fatalln(err)
	}
	data := aboutData
	data.Version = version.FormatVersion()
	data.Admins = user.ListUsersWithGroup("admin")
	data.UserCount = user.Count()
	data.RegistrationLimit = cfg.RegistrationLimit
	data.Cfg = map[string]interface{}{
		"UseAuth":           cfg.UseAuth,
		"WikiName":          cfg.WikiName,
		"HomeHypha":         cfg.HomeHypha,
		"TelegramEnabled":   cfg.TelegramEnabled,
		"RegistrationLimit": cfg.RegistrationLimit,
	}
	var out strings.Builder
	err = temp.Execute(&out, data)
	if err != nil {
		log.Println(err)
	}
	return out.String()
}
