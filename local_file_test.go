package static

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalFile(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := os.CreateTemp(testRoot, "")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString("Gin Web Framework")
	if err != nil {
		t.Error(err)
	}
	f.Close()

	dir, filename := filepath.Split(f.Name())
	router := gin.New()
	router.Use(Serve("/", LocalFile(dir, true)))

	w := PerformRequest(router, "GET", "/"+filename)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Gin Web Framework", w.Body.String())

	w = PerformRequest(router, "GET", "/")
	assert.Contains(t, w.Body.String(), `<a href="`+filename)
}

func TestLocalFileMountPrefixes(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "nested"), 0o700))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "plain"), 0o700))
	for name, body := range map[string]string{
		"file.txt":          "static file",
		INDEX:               "root index",
		"nested/index.html": "nested index",
		"plain/file.txt":    "directory file",
	} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600))
	}

	for _, mount := range []struct {
		name   string
		prefix string
	}{
		{"empty", ""},
		{"root", "/"},
		{"nested", "/assets"},
	} {
		t.Run(mount.name, func(t *testing.T) {
			prefix := mount.prefix
			for _, indexes := range []bool{false, true} {
				name := "without_listing"
				if indexes {
					name = "with_listing"
				}
				t.Run(name, func(t *testing.T) {
					router := gin.New()
					router.Use(Serve(prefix, LocalFile(dir, indexes)))
					router.NoRoute(func(c *gin.Context) {
						c.String(http.StatusNotFound, "fallback")
					})

					for _, tt := range []struct {
						name   string
						method string
						path   string
						status int
						body   string
					}{
						{"file", http.MethodGet, "/file.txt", http.StatusOK, "static file"},
						{"head", http.MethodHead, "/file.txt", http.StatusOK, ""},
						{"root_index", http.MethodGet, "/", http.StatusOK, "root index"},
						{"nested_index", http.MethodGet, "/nested/", http.StatusOK, "nested index"},
						{"missing_file", http.MethodGet, "/missing.txt", http.StatusNotFound, "fallback"},
					} {
						t.Run(tt.name, func(t *testing.T) {
							path := tt.path
							if prefix != "/" {
								path = prefix + path
							}
							w := PerformRequest(router, tt.method, path)
							assert.Equal(t, tt.status, w.Code)
							assert.Equal(t, tt.body, w.Body.String())
						})
					}

					t.Run("directory_without_index", func(t *testing.T) {
						path := "/plain/"
						if prefix != "/" {
							path = prefix + path
						}
						w := PerformRequest(router, http.MethodGet, path)
						if indexes {
							assert.Contains(t, w.Body.String(), `<a href="file.txt">`)
						} else {
							assert.Equal(t, http.StatusNotFound, w.Code)
							assert.Equal(t, "fallback", w.Body.String())
						}
					})

					if prefix == "/assets" {
						t.Run("outside_prefix", func(t *testing.T) {
							w := PerformRequest(router, http.MethodGet, "/file.txt")
							assert.Equal(t, http.StatusNotFound, w.Code)
							assert.Equal(t, "fallback", w.Body.String())
						})
					}
				})
			}
		})
	}
}
