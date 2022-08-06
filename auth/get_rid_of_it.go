package auth

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
