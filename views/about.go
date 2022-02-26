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

func e() l10nEntry {
	return l10nEntry{}
}

func (e l10nEntry) ru(v string) l10nEntry {
	e._ru = v
	return e
}

func (e l10nEntry) en(v string) l10nEntry {
	e._en = v
	return e
}

func (e l10nEntry) get(lang string) string {
	if lang == "ru" && e._ru != "" {
		return e._ru
	}
	return e._en
}

const aboutTemplateString = `<div class="layout">
<main class="main-width">
	<section>
		<h1>{{ printf (get .L.Title) .Cfg.WikiName }}</h1>
		<ul>
			<li><b>{{ get .L.Version }}</b> 1.8.2</li>
		{{ if .Cfg.UseAuth }}
			<li><b>{{ get .L.UserCount }}</b> {{ .UserCount }}</li>
			<li><b>{{ get .L.HomePage }}</b> <a href="/">{{ .Cfg.HomeHypha }}</a></li>
			<li><b>{{ get .L.Admins }}</b> {{$cfg := .Cfg}}{{ range $i, $username := .Admins }}
				{{ if gt $i 0 }}<span aria-hidden="true">, </span>{{ end }}
				<a href="/hypha/{{ $cfg.UserHypha }}/{{ $username }}">{{ $username }}</a>
			{{ end }}</li>
		{{ else }}
			<li>{{ get .L.NoAuth }}</li>
		{{ end }}
		</ul>
		<p>{{ get .L.AboutHyphae }}</p>
	</section>
</main>
</div>`

var aboutData = struct {
	L         map[string]l10nEntry
	Cfg       map[string]interface{}
	Admins    []string
	UserCount uint64
}{
	L: map[string]l10nEntry{
		"Title":       e().en("About %s").ru("О %s"),
		"Version":     e().en("<a href=\"https://mycorrhiza.wiki\">Mycorrhiza Wiki</a> version:").ru("Версия <a href=\"https://mycorrhiza.wiki\">Микоризы</a>:"),
		"UserCount":   e().en("User count:").ru("Число пользователей:"),
		"HomePage":    e().en("Home page:").ru("Домашняя гифа:"),
		"Admins":      e().en("Administrators:").ru("Администраторы:"),
		"NoAuth":      e().en("This wiki does not use authorization").ru("На этой вики не используется авторизация"),
		"AboutHyphae": e().en("See <a href=\"/list\">/list</a> for information about hyphae on this wiki.").ru("См. <a href=\"/list\">/list</a>, чтобы узнать о гифах в этой вики."),
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
		"UseAuth":   cfg.UseAuth,
		"WikiName":  cfg.WikiName,
		"HomeHypha": cfg.HomeHypha,
		"UserHypha": cfg.UserHypha,
	}
	var out strings.Builder
	err = temp.Execute(&out, data)
	if err != nil {
		log.Println(err)
	}
	return out.String()
}
