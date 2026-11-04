package catstorage

import (
	"github.com/bouncepaw/mycorrhiza/internal/categories/domain"
	"github.com/bouncepaw/mycorrhiza/util/bxjson"
)

type record struct {
	Categories []catdomain.Category `json:"categories"`
}

type Storage struct {
	categories *bxjson.JSONFile[record]
}

func New(path string) *Storage {
	return &Storage{categories: bxjson.New[record](path)}
}

func (s *Storage) LoadCategories() ([]catdomain.Category, error) {
	rec, err := s.categories.Read()
	return rec.Categories, err
}

func (s *Storage) SaveCategories(categories []catdomain.Category) error {
	return s.categories.Write(record{Categories: categories})
}
