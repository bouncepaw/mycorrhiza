package storage

import (
	"github.com/hashicorp/go-memdb"
)

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

var DB *memdb.MemDB

func init() {
	var err error
	DB, err = memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
}

func ForEveryRecord(table string, λ func(obj interface{})) error {
	txn := DB.Txn(false)
	defer txn.Abort()

	it, err := txn.Get(table, "id")
	if err != nil {
		return err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		λ(obj)
	}

	return nil
}

func TxnW() *memdb.Txn { return DB.Txn(true) }
func TxnR() *memdb.Txn { return DB.Txn(false) }
