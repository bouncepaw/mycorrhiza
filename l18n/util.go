package l18n

import (
	"fmt"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"net/http"
	"strings"
)

// matcher is a language.Matcher configured for all supported languages.
var locales = language.NewMatcher([]language.Tag{
	language.Make("en"),
	language.Make("ru"),
})

// FromRequest takes a HTTP request and picks the most appropriate localizer (with English fallback)
func FromRequest(r *http.Request) *Localizer {
	t, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	tag, _, _ := locales.Match(t...)
	// TODO: support subtags such as en-US, en-GB, zh-Hans
	base, _ := tag.Base()
	return New(base.String(), "en")
}

var formNames = map[plural.Form]string{
	plural.Other: "other",
	plural.Zero:  "zero",
	plural.One:   "one",
	plural.Two:   "two",
	plural.Few:   "few",
	plural.Many:  "many",
}

func (t Localizer) rawPlural(lang, rawKey string, n int) (string, bool) {
	key := t.getLocalizationKey(lang, rawKey)
	str, ok := t.Localizations[key]
	if !ok {
		return key, false
	}
	var (
		formIdx = plural.Cardinal.MatchPlural(language.Make(lang), n, 0, 0, 0, 0)
		form    = formNames[formIdx]
	)
	plural, plOk := t.Localizations[fmt.Sprintf("%v+%v", key, form)]
	if !plOk {
		return key, false
	}
	return fmt.Sprintf(str, plural), true
}

// GetPlural gets a translated string respecting locale-specific plural rules. Technically, it replaces %s token with +form subkey and proceed as usual.
func (t Localizer) GetPlural(key string, n int, replacements ...*Replacements) string {
	str, ok := t.rawPlural(t.Locale, key, n)
	if !ok {
		str, ok = t.rawPlural(t.FallbackLocale, key, n)
		if !ok {
			return key
		}
	}

	// As in the original, we skip templating if we have nothing to replace (however, it's strange case for plurals)
	if strings.Index(str, "}}") == -1 {
		return str
	}

	return t.replace(str, append(replacements, &Replacements{"n": n})...)
}

// GetPlural64 is ditto for int64
func (t Localizer) GetPlural64(key string, n int64, replacements ...*Replacements) string {
	str, ok := t.rawPlural(t.Locale, key, int(n%1000000))
	if !ok {
		str, ok = t.rawPlural(t.FallbackLocale, key, int(n%1000000))
		if !ok {
			return key
		}
	}

	// As in the original, we skip templating if we have nothing to replace (however, it's strange case for plurals)
	if strings.Index(str, "}}") == -1 {
		return str
	}

	return t.replace(str, append(replacements, &Replacements{"n": n})...)
}

func getLocalizationKey(locale string, key string) string {
	return fmt.Sprintf("%v.%v", locale, key)
}

/* chekoopa: Missing translation features:
- history records (they use Git description, the possible solution is to parse and translate)
- history dates (history.WithRevisions doesn't consider locale, Monday package is bad idea)
- probably error messages (which are scattered across the code)
- default top bar (it is static from one-shot cfg.SetDefaultHeaderLinks, but it is possible to track default-ness in templates)
	- alt solution is implementing "special" links
- dynamic UI (JS are static, though we may send some strings through templates)
- help switches, like,
  - "Read in your language"
  - "Try reading it in English", if no page found in a foreign locale
- feeds (it seems diffcult to pull locale here)
	We do not translate:
- stdout traces (logging is English-only)
*/
