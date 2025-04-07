package static

import (
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
)

// embedFileSystem is a struct that implements the http.FileSystem interface for embedding a file system.
type embedFileSystem struct {
	http.FileSystem
}

// Exists method checks if the given file path exists in the embedded file system.
// If the path exists, it returns true; otherwise, it returns false.
func (e embedFileSystem) Exists(prefix string, path string) bool {
	if len(prefix) > 1 && strings.HasPrefix(path, prefix) {
		path = strings.TrimPrefix(path, prefix)
	}

	_, err := e.Open(path)
	return err == nil
}

// EmbedFolder function embeds the target folder from the embedded file system into the ServeFileSystem.
// If an error occurs during the embedding process, it returns the error message.
func EmbedFolder(fsEmbed embed.FS, targetPath string) (ServeFileSystem, error) {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		slog.Error("Failed to embed folder",
			"targetPath", targetPath,
			"error", err,
		)
		return nil, errors.New("failed to embed folder: " + err.Error())
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}, nil
}
