package static

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequestWithContext(context.Background(), method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestEmptyDirectory(t *testing.T) {
	testRoot, _ := os.Getwd()
	f, err := os.CreateTemp(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("Gin Web Framework")
	f.Close()

	dir, filename := filepath.Split(f.Name())

	router := gin.New()
	router.Use(ServeRoot("/", dir))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "index")
	})
	router.GET("/a", func(c *gin.Context) {
		c.String(http.StatusOK, "a")
	})
	router.GET("/"+filename, func(c *gin.Context) {
		c.String(http.StatusOK, "this is not printed")
	})
	w := PerformRequest(router, "GET", "/")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "index", w.Body.String())

	w = PerformRequest(router, "GET", "/"+filename)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Gin Web Framework", w.Body.String())

	w = PerformRequest(router, "GET", "/"+filename+"a")
	assert.Equal(t, 404, w.Code)

	w = PerformRequest(router, "GET", "/a")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "a", w.Body.String())

	router2 := gin.New()
	router2.Use(ServeRoot("/static", dir))
	router2.GET("/"+filename, func(c *gin.Context) {
		c.String(http.StatusOK, "this is printed")
	})

	w = PerformRequest(router2, "GET", "/")
	assert.Equal(t, 404, w.Code)

	w = PerformRequest(router2, "GET", "/static")
	assert.Equal(t, 404, w.Code)
	router2.GET("/static", func(c *gin.Context) {
		c.String(http.StatusOK, "index")
	})

	w = PerformRequest(router2, "GET", "/static")
	assert.Equal(t, 200, w.Code)

	w = PerformRequest(router2, "GET", "/"+filename)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is printed", w.Body.String())

	w = PerformRequest(router2, "GET", "/static/"+filename)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Gin Web Framework", w.Body.String())
}

func TestIndex(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := os.Create(path.Join(testRoot, "index.html"))
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("index")
	f.Close()

	dir, filename := filepath.Split(f.Name())

	router := gin.New()
	router.Use(ServeRoot("/", dir))

	w := PerformRequest(router, "GET", "/"+filename)
	assert.Equal(t, 301, w.Code)

	w = PerformRequest(router, "GET", "/")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "index", w.Body.String())
}
