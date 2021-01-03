package hyphae

import (
	"github.com/hashicorp/go-memdb"
)

type Hypha struct {
	Name       string
	Exists     bool
	TextPath   string
	BinaryPath string
	OutLinks   []string
	BackLinks  []string
}

// AddHypha adds a hypha named `name` with such `textPath` and `binaryPath`. Both paths can be empty. Does //not// check for hypha's existence beforehand. Count is handled.
func AddHypha(name, textPath, binaryPath string) {
	txn := db.Txn(true)
	txn.Insert("hyphae",
		&Hypha{
			Name:       name,
			TextPath:   textPath,
			BinaryPath: binaryPath,
			OutLinks:   make([]string, 0),
			BackLinks:  make([]string, 0),
		})
	txn.Commit()
	IncrementCount()
}

// DeleteHypha clears both paths and all out-links from the named hypha and marks it as non-existent. It does not actually delete it from the memdb. Count is handled.
func DeleteHypha(name string) {
}

// Create the DB schema
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"hyphae": &memdb.TableSchema{
			Name: "hyphae",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Name"},
				},
				"exists": &memdb.IndexSchema{
					Name:    "exists",
					Unique:  false,
					Indexer: &memdb.BoolFieldIndex{Field: "Exists"},
				},
				"text-path": &memdb.IndexSchema{
					Name:    "text-path",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "TextPath"},
				},
				"binary-path": &memdb.IndexSchema{
					Name:    "binary-path",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "BinaryPath"},
				},
				"out-links": &memdb.IndexSchema{
					Name:    "out-links",
					Unique:  false,
					Indexer: &memdb.StringSliceFieldIndex{Field: "OutLinks"},
				},
				"back-links": &memdb.IndexSchema{
					Name:    "back-links",
					Unique:  false,
					Indexer: &memdb.StringSliceFieldIndex{Field: "BackLinks"},
				},
			},
		},
	},
}

var db *memdb.MemDB

func init() {
	var err error
	db, err = memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
}
