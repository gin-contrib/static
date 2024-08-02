package static

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(strings.TrimPrefix(path, prefix))
	return err == nil
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		log.Fatalf("static.EmbedFolder - Invalid targetPath value - %s", err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
