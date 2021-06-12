package static

import (
	"embed"
	"io/fs"
	"log"
	"os"
)

//go:embed *.css *.js icon
var embedFS embed.FS

// FS serves all static files.
var FS HybridFS

// HybridFS is a filesystem that implements fs.FS. It can serve files
// from multiple filesystems, falling back on failures.
type HybridFS struct {
	fs []fs.FS
}

// Open tries to open the requested file using all filesystems provided.
// If neither succeeds, it returns the last error.
func (f HybridFS) Open(name string) (fs.File, error) {
	log.Printf("serving static file: %s\n", name)

	var file fs.File
	var err error

	for _, candidate := range f.fs {
		file, err = candidate.Open(name)
		if err == nil {
			log.Println("succeeded")
			return file, nil
		}
	}

	log.Printf("failed: %v\n", err)
	return nil, err
}

// InitFS initializes the global HybridFS singleton with the local wiki.
func InitFS(localPath string) {
	FS = HybridFS{
		fs: []fs.FS{
			os.DirFS(localPath),
			embedFS,
		},
	}
}
