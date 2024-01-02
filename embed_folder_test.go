package static_test

import (
	"embed"
	"fmt"
	"log"
	"testing"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//go:embed test/data/server
var testFS embed.FS

func TestEmbedFolderWithRedir(t *testing.T) {
	var tests = []struct {
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

	router := gin.New()

	staticFiles, err := static.EmbedFolder(testFS, "test/data/server")
	if err != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	} else {
		router.Use(static.Serve("/", staticFiles))
	}

	router.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(301, "/")
	})

	for _, tt := range tests {
		w := PerformRequest(router, "GET", tt.targetURL)
		assert.Equal(t, tt.httpCode, w.Code, tt.name)
		assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
	}
}

func TestEmbedFolderWithoutRedir(t *testing.T) {
	var tests = []struct {
		targetURL string // input
		httpCode  int    // expected http code
		httpBody  string // expected http body
		name      string // test name
	}{
		{"/404.html", 404, "404 page not found", "Unknown file"},
		{"/", 200, "<h1>Hello Embed</h1>", "Root"},
		{"/index.html", 301, "", "Root by file name automatic redirect"},
		{"/static.html", 200, "<h1>Hello Gin Static</h1>", "Other file"},
	}

	router := gin.New()
	staticFiles, err := static.EmbedFolder(testFS, "test/data/server")
	if err != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	} else {
		router.Use(static.Serve("/", staticFiles))
	}

	for _, tt := range tests {
		w := PerformRequest(router, "GET", tt.targetURL)
		assert.Equal(t, tt.httpCode, w.Code, tt.name)
		assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
	}
}

func TestEmbedInitErrorPath(t *testing.T) {
	tests := []struct {
		name       string
		targetPath string
		haveErr    bool
		fs         embed.FS
	}{
		{
			name:       "ValidPath",
			targetPath: "test/data/server",
			haveErr:    false,
			fs:         testFS,
		},
		{
			name:       "InvalidPath",
			targetPath: "nonexistingdirectory/nonexistingdirectory",
			haveErr:    true,
			fs:         testFS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := static.EmbedFolder(tt.fs, tt.targetPath)
			assert.Equal(t, (err != nil), tt.haveErr, tt.name)
		})
	}
}
