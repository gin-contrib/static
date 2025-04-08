package static

import (
	"embed"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//go:embed test/data/server
var server embed.FS

// embedTests defines a set of test cases for testing the behavior of the embedded file system.
var embedTests = []struct {
	targetURL string // input URL
	httpCode  int    // expected HTTP status code
	httpBody  string // expected HTTP response body
	name      string // name of the test case
}{
	{"/404.html", 301, "<a href=\"/\">Moved Permanently</a>.\n\n", "Unknown file"},
	{"/", 200, "<h1>Hello Embed</h1>", "Root"},
	{"/index.html", 301, "", "Root by file name automatic redirect"},
	{"/static.html", 200, "<h1>Hello Gin Static</h1>", "Other file"},
}

// TestEmbedFolder tests the behavior of the embedded file system, ensuring correct handling of static files and directories.
func TestEmbedFolder(t *testing.T) {
	t.Run("EmbedFolder", func(t *testing.T) {
		router := gin.New()
		fs, err := EmbedFolder(server, "test/data/server")
		if err != nil {
			t.Fatalf("Failed to embed folder: %v", err)
		}
		router.Use(Serve("/", fs))
		router.NoRoute(func(c *gin.Context) {
			fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
			c.Redirect(301, "/")
		})

		for _, tt := range embedTests {
			w := PerformRequest(router, "GET", tt.targetURL)
			assert.Equal(t, tt.httpCode, w.Code, tt.name)
			assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
		}
	})

	t.Run("EmbedFolder with prefix", func(t *testing.T) {
		router := gin.New()
		fs, err := EmbedFolder(server, "test/data/server")
		if err != nil {
			t.Fatalf("Failed to embed folder: %v", err)
		}
		router.Use(Serve("/prefix", fs))
		router.NoRoute(func(c *gin.Context) {
			fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
			c.Redirect(301, "/")
		})

		for _, tt := range embedTests {
			w := PerformRequest(router, "GET", "/prefix"+tt.targetURL)
			assert.Equal(t, tt.httpCode, w.Code, tt.name)
			assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
		}
	})
}
