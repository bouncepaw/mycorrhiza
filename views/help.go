package views

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"log"
	"strings"
	"text/template"
)

var helpTopicsL10n = map[string][]string{
	"topics":         {"Help topics", "Темы справки"},
	"main":           {"Main", "Введение"},
	"hypha":          {"Hypha", "Гифа"},
	"media":          {"Media", "Медиа"},
	"mycomarkup":     {"Mycomarkup", "Микоразметка"},
	"interface":      {"Interface", "Интерфейс"},
	"prevnext":       {"Previous/next", "Назад/далее"}, // пред след?
	"top_bar":        {"Top bar", "Верхняя панель"},
	"sibling_hyphae": {"Sibling hyphae", "Гифы-сиблинги"},
	"special_pages":  {"Special pages", "Специальные страницы"},
	"recent_changes": {"Recent changes", "Недавние изменения"}, // так ли? В медиавики свежие правки
	"feeds":          {"Feeds", "Ленты"},
	"configuration":  {"Configuration (for administrators)", "Конфигурация (для администраторов)"},
	"config_file":    {"Configuration file", "Файл конфигурации"},
	"lock":           {"Lock", "Блокировка"}, // Не Замок ли?
	"whitelist":      {"Whitelist", "Белый список"},
	"telegram":       {"Telegram authentication", "Вход через Телеграм"},
}

const helpTopicTemplate = `<aside class="help-topics layout-card">
	<h2 class="layout-card__title">{{l "topics"}}</h2>
	<ul class="help-topics__list">
		<li>{{l "main" | a ""}}</li>
		<li>{{l "hypha" | a "/hypha"}}
			<ul>
				{{l "media" | a "/media"}}
			</ul>
		</li>
		<li>{{l "mycomarkup" | a "/mycomarkup"}}</li>
		<li>{{l "interface"}}
			<ul>
				<li>{{l "prevnext" | a "/prevnext"}}</li>
				<li>{{l "top_bar" | a "/top_bar"}}</li>
				<li>{{l "sibling_hyphae" | a "/sibling_hyphae_section"}}</li>
			</ul>
		</li>
		<li>{{l "special_pages"}}
			<ul>
				<li>{{l "recent_changes" | a "/recent_changes"}}</li>
				<li>{{l "feeds" | a "/feeds"}}</li>
			</ul>
		</li>
		<li>{{l "configuration"}}
			<ul>
				<li>{{l "config_file" | a "/config_file"}}</li>
				<li>{{l "lock" | a "/lock"}}</li>
				<li>{{l "whitelist" | a "/whitelist"}}</li>
				<li>{{l "telegram" | a "/telegram"}}</li>
			</ul>
		</li>
	</ul>
</aside>`

// helpTopicsLinkWrapper wraps in <a>
func helpTopicsLinkWrapper(lang string) func(string, string) string {
	return func(path, contents string) string {
		return fmt.Sprintf(`<a href="/help/%s%s">%s</a>`, lang, path, contents)
	}
}

func helpTopicsLocalizedTopic(lang string) func(string) string {
	pos := 0
	if lang == "ru" {
		pos = 1
	}
	return func(topic string) string {
		return helpTopicsL10n[topic][pos]
	}
}

func helpTopicsHTML(lang string, lc *l18n.Localizer) string {
	temp, err := template.
		New("help topics").
		Funcs(template.FuncMap{
			"a": helpTopicsLinkWrapper(lang),
			"l": helpTopicsLocalizedTopic(lc.Locale),
		}).
		Parse(helpTopicTemplate)
	if err != nil {
		log.Println(err)
		return ""
	}

	// TODO: one day, it should write to a better place
	var out strings.Builder
	_ = temp.Execute(&out, nil) // Shall not fail!
	return out.String()
}
