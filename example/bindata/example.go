package main

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func binaryFS(root string) *binaryFileSystem {
	return &binaryFileSystem{
		AssetFile(),
	}
}

// Usage
// $ go-bindata -prefix "data" -fs data/
// $ go build && ./bindata
//
func main() {
	r := gin.Default()

	r.Use(static.Serve("/static", binaryFS("data")))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
