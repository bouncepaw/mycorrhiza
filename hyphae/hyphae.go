package hyphae

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	markup.HyphaExists = func(hyphaName string) bool {
		return ByName(hyphaName).Exists
	}
	markup.HyphaAccess = func(hyphaName string) (rawText, binaryBlock string, err error) {
		if h := ByName(hyphaName); h.Exists {
			rawText, err = h.FetchTextPart()
			if h.BinaryPath != "" {
				binaryBlock = h.BinaryHtmlBlock()
			}
		} else {
			err = errors.New("Hypha " + hyphaName + " does not exist")
		}
		return
	}
	markup.HyphaIterate = func(λ func(string)) {
		for h := range YieldExistingHyphae() {
			λ(h.Name)
		}
	}
	markup.HyphaImageForOG = func(hyphaName string) string {
		if h := ByName(hyphaName); h.Exists && h.BinaryPath != "" {
			return util.URL + "/binary/" + hyphaName
		}
		return util.URL + "/favicon.ico"
	}
}

// HyphaPattern is a pattern which all hyphae must match.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"\'&%{}]+`)

type Hypha struct {
	sync.RWMutex

	Name       string
	Exists     bool
	TextPath   string
	BinaryPath string
	OutLinks   []*Hypha
	BackLinks  []*Hypha
}

var byNames = make(map[string]*Hypha)
var byNamesMutex = sync.Mutex{}

// YieldExistingHyphae iterates over all hyphae and yields all existing ones.
func YieldExistingHyphae() chan *Hypha {
	ch := make(chan *Hypha)
	go func() {
		for _, h := range byNames {
			if h.Exists {
				ch <- h
			}
		}
		close(ch)
	}()
	return ch
}

// FilterTextHyphae filters the source channel and yields only those hyphae than have text parts.
func FilterTextHyphae(src chan *Hypha) chan *Hypha {
	sink := make(chan *Hypha)
	go func() {
		for h := range src {
			if h.TextPath != "" {
				sink <- h
			}
		}
		close(sink)
	}()
	return sink
}

// Subhyphae returns slice of subhyphae.
func (h *Hypha) Subhyphae() []*Hypha {
	hyphae := []*Hypha{}
	for subh := range YieldExistingHyphae() {
		if strings.HasPrefix(subh.Name, h.Name+"/") {
			hyphae = append(hyphae, subh)
		}
	}
	return hyphae
}

// AreFreeNames checks if all given `hyphaNames` are not taken.
func AreFreeNames(hyphaNames ...string) (firstFailure string, ok bool) {
	for h := range YieldExistingHyphae() {
		for _, hn := range hyphaNames {
			if hn == h.Name {
				return hn, false
			}
		}
	}
	return "", true
}

// EmptyHypha returns an empty hypha struct with given name.
func EmptyHypha(hyphaName string) *Hypha {
	return &Hypha{
		Name:       hyphaName,
		Exists:     false,
		TextPath:   "",
		BinaryPath: "",
		OutLinks:   make([]*Hypha, 0),
		BackLinks:  make([]*Hypha, 0),
	}
}

// ByName returns a hypha by name. If h.Exists, the returned hypha pointer is known to be part of the hypha index (byNames map).
func ByName(hyphaName string) (h *Hypha) {
	h, exists := byNames[hyphaName]
	if exists {
		return h
	}
	return EmptyHypha(hyphaName)
}

// Insert inserts the hypha into the storage. It overwrites the previous record, if there was any, and returns false. If the was no previous record, return true.
func (h *Hypha) Insert() (justCreated bool) {
	hp := ByName(h.Name)

	byNamesMutex.Lock()
	defer byNamesMutex.Unlock()
	if hp.Exists {
		hp = h
	} else {
		h.Exists = true
		byNames[h.Name] = h
		IncrementCount()
	}

	return !hp.Exists
}

func (h *Hypha) InsertIfNew() (justCreated bool) {
	if !h.Exists {
		return h.Insert()
	}
	return false
}

func (h *Hypha) InsertIfNewKeepExistence() {
	hp := ByName(h.Name)

	byNamesMutex.Lock()
	defer byNamesMutex.Unlock()
	if hp.Exists {
		hp = h
	} else {
		byNames[h.Name] = h
	}
}

func (h *Hypha) delete() {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.Name)
	DecrementCount()
	byNamesMutex.Unlock()
	h.Unlock()
}

func (h *Hypha) renameTo(newName string) {
	byNamesMutex.Lock()
	h.Lock()
	delete(byNames, h.Name)
	h.Name = newName
	byNames[h.Name] = h
	byNamesMutex.Unlock()
	h.Unlock()
}

// MergeIn merges in content file paths from a different hypha object. Prints warnings sometimes.
func (h *Hypha) MergeIn(oh *Hypha) {
	if h.TextPath == "" && oh.TextPath != "" {
		h.TextPath = oh.TextPath
	}
	if oh.BinaryPath != "" {
		if h.BinaryPath != "" {
			log.Println("There is a file collision for binary part of a hypha:", h.BinaryPath, "and", oh.BinaryPath, "-- going on with the latter")
		}
		h.BinaryPath = oh.BinaryPath
	}
}
