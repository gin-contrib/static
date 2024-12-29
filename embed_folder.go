package static

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		slog.Error("Failed to embed folder",
			"targetPath", targetPath,
			"error", err,
		)
		return nil
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
