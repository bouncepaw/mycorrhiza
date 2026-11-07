package misc

import (
	"log/slog"
	"os"
	"strings"
	"text/template" // sic! TODO: make it html/template after the template library migration

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/internal/version"
)

const aboutTemplateString = `
<main class="main-width">
	<section class="about-page">
		<h1>About {{ .Cfg.WikiName }}</h1>
		<dl>
			<dt><a href="https://mycorrhiza.wiki">Mycorrhiza Wiki</a> version</dt>
			<dd>{{ .Version }}</dd>
		{{ if .Cfg.UseAuth }}
			<dt>Home hypha</dt>
			<dd><a href="/">{{ .Cfg.HomeHypha }}</a></dd>

			<dt>Authentication</dt>
			<dd>Authentication is on</dd>
			{{if .Cfg.TelegramEnabled}}<dd>Telegram authentication is on</dd>{{end}}

			<dt>User count</dt>
			<dd>{{ .UserCount }}</dd>
			{{if .Cfg.RegistrationLimit}}
			<dt>Registration limit</dt>
			<dd>{{.RegistrationLimit}}</dd>
			{{end}}

			<dt>Administrators</dt>
			{{$cfg := .Cfg}}{{ range $i, $username := .Admins }}
				<dd><a href="/hypha/{{ $cfg.UserHypha }}/{{ $username }}">{{ $username }}</a></dd>
			{{ end }}

		{{ else }}
			<dt>Authentication</dt>
			<dd>Authentication is off</dd>
		{{ end }}
		</dl>
	</section>
</main>`

var aboutData = struct {
	Version           string
	Cfg               map[string]any
	Admins            []string
	UserCount         uint64
	RegistrationLimit uint64
}{}

func AboutHTML() string {
	temp, err := template.New("about wiki").Parse(aboutTemplateString)
	if err != nil {
		slog.Error("Failed to parse About template", "err", err)
		os.Exit(1)
	}
	data := aboutData
	data.Version = version.Short
	data.Admins = usersvc.ListUsersWithGroup(userdomain.GroupAdmin)
	data.UserCount = usersvc.Count()
	data.RegistrationLimit = cfg.RegistrationLimit
	data.Cfg = map[string]any{
		"UseAuth":           cfg.UseAuth,
		"WikiName":          cfg.WikiName,
		"HomeHypha":         cfg.HomeHypha,
		"TelegramEnabled":   cfg.TelegramEnabled,
		"RegistrationLimit": cfg.RegistrationLimit,
	}
	var out strings.Builder
	err = temp.Execute(&out, data)
	if err != nil {
		slog.Error("Failed to execute About template", "err", err)
		os.Exit(1)
	}
	return out.String()
}
