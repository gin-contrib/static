package static

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Gin Web Framework")

	w = PerformRequest(router, "GET", "/")
	assert.Contains(t, w.Body.String(), `<a href="`+filename)
}

// TestLocalFileIndexConstants documents that the exported ListDirectory and
// NoListDirectory constants map to the historical true/false values of the
// indexes parameter, so existing call sites keep working.
func TestLocalFileIndexConstants(t *testing.T) {
	assert.True(t, ListDirectory)
	assert.False(t, NoListDirectory)

	fs := LocalFile("/tmp", ListDirectory)
	assert.True(t, fs.indexes)

	fs = LocalFile("/tmp", NoListDirectory)
	assert.False(t, fs.indexes)
}

// TestLocalFileNoListDirectory verifies the documented NoListDirectory
// behavior: directory listing is suppressed, but an index.html inside the
// directory is still served.
func TestLocalFileNoListDirectory(t *testing.T) {
	dir, err := os.MkdirTemp("", "static-nolist")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.Join(dir, INDEX), []byte("Gin Web Framework"), 0o600); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "secret.txt"), []byte("top secret"), 0o600); err != nil {
		t.Fatalf("Failed to write secret file: %v", err)
	}

	router := gin.New()
	router.Use(Serve("/", LocalFile(dir, NoListDirectory)))

	// The index.html is served for the directory root.
	w := PerformRequest(router, "GET", "/")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Gin Web Framework", w.Body.String())

	// The directory listing is not exposed: the sibling file name does not
	// appear as an auto-generated listing link.
	assert.NotContains(t, w.Body.String(), `<a href="secret.txt`)
}
