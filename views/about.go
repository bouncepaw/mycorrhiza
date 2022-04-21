package views

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"log"
	"strings"
	"text/template" // sic!
)

type l10nEntry struct {
	_en string
	_ru string
}

func en(v string) l10nEntry {
	return l10nEntry{_en: v}
}

func (e l10nEntry) ru(v string) l10nEntry {
	e._ru = v
	return e
}

func (e l10nEntry) get(lang string) string {
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
			<dd>1.10.0</dd>
		{{ if .Cfg.UseAuth }}
			<dt>{{ get .L.HomeHypha }}</dt>
			<dd><a href="/">{{ .Cfg.HomeHypha }}</a></dd>

			<dt>{{get .L.Auth}}</dt>
			<dd>{{ get .L.AuthOn }}</dd>
			{{if .Cfg.TelegramEnabled}}<dd>{{get .L.TelegramOn}}</dd>{{end}}

			<dt>{{ get .L.UserCount }}</dt>
			<dd>{{ .UserCount }}</dd>

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
	L         map[string]l10nEntry
	Cfg       map[string]interface{}
	Admins    []string
	UserCount uint64
}{
	L: map[string]l10nEntry{
		"Title":     en("About %s").ru("О %s"),
		"Version":   en("<a href=\"https://mycorrhiza.wiki\">Mycorrhiza Wiki</a> version").ru("Версия <a href=\"https://mycorrhiza.wiki\">Микоризы</a>"),
		"UserCount": en("User count").ru("Число пользователей"),
		"HomeHypha": en("Home hypha").ru("Домашняя гифа"),
		"Admins":    en("Administrators").ru("Администраторы"),

		"Auth":       en("Authentication").ru("Аутентификация"),
		"AuthOn":     en("Authentication is on").ru("Аутентификация включена"),
		"AuthOff":    en("Authentication is off").ru("Аутентификация не включена"),
		"TelegramOn": en("Telegram authentication is on").ru("Вход через Телеграм включён"),
	},
}

func AboutHTML(lc *l18n.Localizer) string {
	get := func(e l10nEntry) string {
		return e.get(lc.Locale)
	}
	temp, err := template.New("about wiki").Funcs(template.FuncMap{"get": get}).Parse(aboutTemplateString)
	if err != nil {
		log.Fatalln(err)
	}
	data := aboutData
	data.Admins = user.ListUsersWithGroup("admin")
	data.UserCount = user.Count()
	data.Cfg = map[string]interface{}{
		"UseAuth":         cfg.UseAuth,
		"WikiName":        cfg.WikiName,
		"HomeHypha":       cfg.HomeHypha,
		"TelegramEnabled": cfg.TelegramEnabled,
	}
	var out strings.Builder
	err = temp.Execute(&out, data)
	if err != nil {
		log.Println(err)
	}
	return out.String()
}
