package static

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) (ServeFileSystem, error) {
	fsys, _ := fs.Sub(fsEmbed, targetPath)
	_, err := os.Stat(targetPath)
	if err != nil {
		return nil, err
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}, nil
}
