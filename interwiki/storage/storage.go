package interwikistorage

import (
	"github.com/bouncepaw/mycorrhiza/interwiki/domain"
	"github.com/bouncepaw/mycorrhiza/util/bxjson"
)

type Storage struct {
	entries *bxjson.JSONFile[[]interwikidomain.WikiDTO]
}

func New(path string) *Storage {
	return &Storage{entries: bxjson.New[[]interwikidomain.WikiDTO](path)}
}

func (s *Storage) LoadEntries() ([]interwikidomain.WikiDTO, error) {
	return s.entries.Read()
}

func (s *Storage) SaveEntries(entries []interwikidomain.WikiDTO) error {
	return s.entries.Write(entries)
}
