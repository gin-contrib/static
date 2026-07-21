package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const INDEX = "index.html"

// Named boolean values for the indexes parameter of LocalFile, so calls read
// clearly at the call site instead of using a bare true/false. For example:
//
//	router.Use(static.Serve("/", static.LocalFile("/www", static.NoListDirectory)))
const (
	// ListDirectory enables directory listing: requesting a directory returns
	// a listing of the files it contains.
	ListDirectory = true
	// NoListDirectory disables directory listing: requesting a directory is
	// served only when it contains an index.html file (which is then returned).
	NoListDirectory = false
)

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

// LocalFile serves files from the local directory rooted at root.
//
// The indexes parameter controls how directories are handled:
//   - when true (ListDirectory), directory listing is enabled and requesting a
//     directory returns a listing of its files;
//   - when false (NoListDirectory), directory listing is disabled and a
//     directory is served only when it contains an index.html file, which is
//     then returned.
//
// Note that an existing index.html is served in both cases; indexes only
// toggles the automatic directory listing.
func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, path string) bool {
	// Check if path starts with prefix
	p := strings.TrimPrefix(path, prefix)
	if len(p) >= len(path) {
		return false
	}

	name := filepath.Join(l.root, p)
	stats, err := os.Stat(name)
	if err != nil {
		return false
	}

	// If it's a directory and indexes are disabled, check for index file
	if stats.IsDir() && !l.indexes {
		indexPath := filepath.Join(name, INDEX)
		_, err := os.Stat(indexPath)
		return err == nil
	}

	return true
}
