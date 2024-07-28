package static

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed *.css *.js *.txt icon
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
	var file fs.File
	var err error

	for _, candidate := range f.fs {
		file, err = candidate.Open(name)
		if err == nil {
			return file, nil
		}
	}

	return nil, err
}

// InitFS initializes the global HybridFS singleton with the wiki's own static
// files directory as a primary filesystem and the embedded one as a fallback.
func InitFS(localPath string) {
	FS = HybridFS{
		fs: []fs.FS{
			os.DirFS(localPath),
			embedFS,
		},
	}
}
