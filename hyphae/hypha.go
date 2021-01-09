package hyphae

import (
	"github.com/bouncepaw/mycorrhiza/storage"
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
	txn := storage.DB.Txn(true)
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
