// Package l18n handles everything language-related for Mycorrhiza.
package l18n

/*
Some code is borrowed from github.com/m1/go-localize. The copyright notice is
included as required by the MIT License:

	Copyright (c) 2019 Miles Croxford <hello@milescroxford.com>

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:
*/

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

type Replacements map[string]interface{}

type Localizer struct {
	Locale         string
	FallbackLocale string
	Localizations  map[string]string
}

// locales is a filesystem containing all localization files.
//go:embed en ru
var locales embed.FS

// localization maps localizable keys to strings.
var localizations = make(map[string]string)

// supportedLocales is a language.Matcher configured for supported languages.
var supportedLocales = language.NewMatcher([]language.Tag{
	language.Make("en"),
	language.Make("ru"),
})

var formNames = map[plural.Form]string{
	plural.Other: "other",
	plural.Zero:  "zero",
	plural.One:   "one",
	plural.Two:   "two",
	plural.Few:   "few",
	plural.Many:  "many",
}

func init() {
	fs.WalkDir(locales, ".", func(path string, d fs.DirEntry, err error) error {
		ext := filepath.Ext(path)
		if !d.IsDir() && ext == ".json" {
			basename := path[:len(path)-len(ext)]
			// Note: embed.FS always uses a forward slash as the path separator.
			segments := strings.Split(basename, "/")
			prefix := strings.Join(segments, ".") + "."

			contents, err := locales.ReadFile(path)
			if err != nil {
				return nil
			}

			var strings map[string]string
			if err := json.Unmarshal(contents, &strings); err != nil {
				log.Fatalf("error while parsing %s: %v", path, err)
			}

			for key, value := range strings {
				localizations[prefix+key] = value
			}
		}
		return nil
	})
}

// New creates a new Localizer with locales set. This operation is cheap.
func New(locale string, fallbackLocale string) *Localizer {
	t := &Localizer{Locale: locale, FallbackLocale: fallbackLocale}
	t.Localizations = localizations
	return t
}

// FromRequest takes a HTTP request and picks the most appropriate localizer
// with English as the fallback language.
func FromRequest(r *http.Request) *Localizer {
	t, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	tag, _, _ := supportedLocales.Match(t...)
	// TODO: support subtags such as en-US, en-GB, zh-Hans
	base, _ := tag.Base()
	return New(base.String(), "en")
}

// GetWithLocale returns a localized string for the provided key in a specific
// locale with optional replacements executed on the string.
func (t Localizer) GetWithLocale(locale, key string, replacements ...*Replacements) string {
	str, ok := t.Localizations[getLocalizationKey(locale, key)]
	if !ok {
		str, ok = t.Localizations[getLocalizationKey(t.FallbackLocale, key)]
		if !ok {
			return key
		}
	}

	// If the str doesn't have any substitutions, no need to
	// template.Execute.
	if strings.Contains(str, "}}") {
		return str
	}

	return t.replace(str, replacements...)
}

// Get returns a localized string for the provided key with optional
// replacements executed on the string.
func (t Localizer) Get(key string, replacements ...*Replacements) string {
	str := t.GetWithLocale(t.Locale, key, replacements...)
	return str
}

// GetPlural returns a localized string respecting locale-specific plural rules.
// Technically, it replaces %s token with +form subkey and proceeds as usual.
func (t Localizer) GetPlural(key string, n int, replacements ...*Replacements) string {
	str, ok := t.rawPlural(t.Locale, key, n)
	if !ok {
		str, ok = t.rawPlural(t.FallbackLocale, key, n)
		if !ok {
			return key
		}
	}

	// As in the original, we skip templating if have nothing to replace
	// (however, it's strange case for plurals)
	if strings.Contains(str, "}}") {
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

	// As in the original, we skip templating if have nothing to replace
	// (however, it's strange case for plurals)
	if strings.Contains(str, "}}") {
		return str
	}

	return t.replace(str, append(replacements, &Replacements{"n": n})...)
}

func (t Localizer) replace(str string, replacements ...*Replacements) string {
	b := &bytes.Buffer{}
	tmpl, err := template.New("").Parse(str)
	if err != nil {
		return str
	}

	replacementsMerge := Replacements{}
	for _, replacement := range replacements {
		for k, v := range *replacement {
			replacementsMerge[k] = v
		}
	}

	err = template.Must(tmpl, err).Execute(b, replacementsMerge)
	if err != nil {
		return str
	}
	buff := b.String()
	return buff
}

func (t Localizer) rawPlural(lang, rawKey string, n int) (string, bool) {
	key := getLocalizationKey(lang, rawKey)
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

bouncepaw:
- more error messages
*/
