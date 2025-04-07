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
