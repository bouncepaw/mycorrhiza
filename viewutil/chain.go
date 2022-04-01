package viewutil

import "text/template"

// Chain represents a chain of different language versions of the same template.
type Chain struct {
	en *template.Template
	ru *template.Template
}

// En returns a new Chain. This is the only constructor of the type, so every view is forced to have an English representation.
func En(en *template.Template) Chain {
	return Chain{
		en: en,
	}
}

// Ru adds a Russian translation to the Chain.
func (c Chain) Ru(ru *template.Template) Chain {
	c.ru = ru
	return c
}

// Get returns an appropriate language representation for the given locale in meta.
func (c Chain) Get(meta Meta) *template.Template {
	switch meta.Locale() {
	case "en":
		return c.en
	case "ru":
		return c.ru
	}
	panic("unknown language " + meta.Locale())
}
