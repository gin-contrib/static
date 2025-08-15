package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const INDEX = "index.html"

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

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
