package static

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

const INDEX = "index.html"

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

type CacheConfigs struct {
	Public    bool
	MaxAge    uint
	Immutable bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// GenericServe returns a middleware handler that serves static files in the given directory.
func GenericServe(urlPrefix string, fs ServeFileSystem, cc CacheConfigs) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			var cacheControl []string
			if cc.Public {
				cacheControl = append(cacheControl, "public")
			}
			if cc.MaxAge != 0 {
				cacheControl = append(cacheControl, fmt.Sprintf("max-age=%d", cc.MaxAge))
			} else {
				cacheControl = append(cacheControl, "no-store")
			}
			if cc.Immutable {
				cacheControl = append(cacheControl, "immutable")
			}
			c.Writer.Header().Add("Cache-Control", strings.Join(cacheControl, ", "))
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	return GenericServe(urlPrefix, fs, CacheConfigs{MaxAge: 0})
}

// ServeCached returns a middleware handler that similar as Serve but with the Cache-Control Header set as passed in the cacheAge parameter
func ServeCached(urlPrefix string, fs ServeFileSystem, cc CacheConfigs) gin.HandlerFunc {
	return GenericServe(urlPrefix, fs, cc)

}
