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

// embedTests 定義了一組測試用例，用於測試嵌入文件系統的行為。
var embedTests = []struct {
	targetURL string // input
	httpCode  int    // expected http code
	httpBody  string // expected http body
	name      string // test name
}{
	{"/404.html", 301, "<a href=\"/\">Moved Permanently</a>.\n\n", "Unknown file"},
	{"/", 200, "<h1>Hello Embed</h1>", "Root"},
	{"/index.html", 301, "", "Root by file name automatic redirect"},
	{"/static.html", 200, "<h1>Hello Gin Static</h1>", "Other file"},
}

// TestEmbedFolder 測試 EmbedFolder 函數的行為，確保其能夠正確嵌入文件夾並提供靜態文件服務。
func TestEmbedFolder(t *testing.T) {
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
}
