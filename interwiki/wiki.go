package interwikisvc

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/bouncepaw/mycorrhiza/interwiki/domain"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Wiki is an entry in the interwiki map.
type Wiki struct {
	interwikidomain.WikiDTO
}

func (w *Wiki) canonize() error {
	switch {
	case w.Name == "":
		return fmt.Errorf("the name is not set")
	case w.URL == "":
		return fmt.Errorf("the URL is not set")
	case !w.Engine.Valid():
		return fmt.Errorf("unknown engine %q", w.Engine)
	}
	if err := checkAddress(w.URL); err != nil {
		return fmt.Errorf("URL: %w", err)
	}

	w.Name = util.CanonicalName(w.Name)
	for i, alias := range w.Aliases {
		w.Aliases[i] = util.CanonicalName(alias)
	}

	if w.LinkHrefFormat == "" {
		switch w.Engine {
		case interwikidomain.Mycorrhiza:
			w.LinkHrefFormat = fmt.Sprintf("%s/hypha/{NAME}", w.URL)
		case interwikidomain.Betula:
			w.LinkHrefFormat = fmt.Sprintf("%s/{BETULA-NAME}", w.URL)
		case interwikidomain.Agora:
			w.LinkHrefFormat = fmt.Sprintf("%s/node/{NAME}", w.URL)
		default:
			w.LinkHrefFormat = fmt.Sprintf("%s/{NAME}", w.URL)
		}
	}

	if w.ImgSrcFormat == "" {
		switch w.Engine {
		case interwikidomain.Mycorrhiza:
			w.ImgSrcFormat = fmt.Sprintf("%s/binary/{NAME}", w.URL)
		default:
			w.ImgSrcFormat = fmt.Sprintf("%s/{NAME}", w.URL)
		}
	}

	if err := checkAddress(w.LinkHrefFormat); err != nil {
		return fmt.Errorf("link format: %w", err)
	}
	if err := checkAddress(w.ImgSrcFormat); err != nil {
		return fmt.Errorf("image format: %w", err)
	}
	return nil
}

func checkAddress(s string) error {
	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsControl(r) || strings.ContainsRune("\"'<>`\\", r) {
			return fmt.Errorf("forbidden character %q", r)
		}
	}

	s = strings.NewReplacer("{NAME}", "name", "{BETULA-NAME}", "name").Replace(s)
	u, err := url.Parse(s)
	switch {
	case err != nil:
		return fmt.Errorf("not a valid URL")
	case u.Scheme == "":
		return fmt.Errorf("no scheme found in URL")
	}
	return nil
}
