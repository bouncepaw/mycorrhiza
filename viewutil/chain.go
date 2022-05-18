package viewutil

import "text/template"

// Chain represents a chain of different language versions of the same template.
type Chain struct {
	_en *template.Template
	_ru *template.Template
}

// en returns a new Chain. This is the only constructor of the type, so every view is forced to have an English representation.
func en(en *template.Template) Chain {
	return Chain{
		_en: en,
	}
}

// ru adds a Russian translation to the Chain.
func (c Chain) ru(ru *template.Template) Chain {
	c._ru = ru
	return c
}

// Get returns an appropriate language representation for the given locale in meta.
func (c Chain) Get(meta Meta) *template.Template {
	switch meta.Locale() {
	case "_en":
		return c._en
	case "_ru":
		return c._ru
	}
	panic("unknown language " + meta.Locale())
}
