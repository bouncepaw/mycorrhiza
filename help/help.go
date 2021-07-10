// Package help contains help messages and the utilities for retrieving them.
package help

import (
	"embed"
)

//go:embed en
var fs embed.FS

// Get determines what help text you need and returns it. The path is a substring from URL, it follows this form:
//     <language>/<topic>
func Get(path string) ([]byte, error) {
	return fs.ReadFile(path + ".myco")
}
