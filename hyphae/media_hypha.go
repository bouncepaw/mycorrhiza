package hyphae

import (
	"github.com/bouncepaw/mycorrhiza/files"
	"path/filepath"
	"sync"
)

type MediaHypha struct {
	sync.RWMutex

	canonicalName string
	mycoFilePath  string
	mediaFilePath string
}

func (m *MediaHypha) DoesExist() {
}

func (m *MediaHypha) CanonicalName() string {
	return m.canonicalName
}

func (m *MediaHypha) TextFilePath() string {
	if m.mycoFilePath == "" {
		return filepath.Join(files.HyphaeDir(), m.CanonicalName()+".myco")
	}
	return m.mycoFilePath
}

func (m *MediaHypha) HasTextFile() bool {
	return m.mycoFilePath != ""
}

func (m *MediaHypha) MediaFilePath() string {
	return m.mediaFilePath
}

func (m *MediaHypha) SetMediaFilePath(newPath string) {
	m.mediaFilePath = newPath
}

func ShrinkMediaToTextual(m *MediaHypha) *TextualHypha {
	return &TextualHypha{
		canonicalName: m.CanonicalName(),
		mycoFilePath:  m.TextFilePath(),
	}
}
