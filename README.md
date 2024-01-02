# Gin Static Middleware

[![Run Tests](https://github.com/gin-contrib/static/actions/workflows/go.yml/badge.svg)](https://github.com/gin-contrib/static/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gin-contrib/static/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/static)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/static)](https://goreportcard.com/report/github.com/gin-contrib/static)
[![GoDoc](https://godoc.org/github.com/gin-contrib/static?status.svg)](https://godoc.org/github.com/gin-contrib/static)

Static middleware, support both local files and embed filesystem.

No historical burden, using go 1.21, 100% code coverage.

## Quick Start

### Download and Import

Download and install it:

```bash
go get github.com/gin-contrib/static
```

Import it in your code:

```go
import "github.com/gin-contrib/static"
```

## Example

See the [example](example)

### Serve Local Files

[local files]:# (example/simple/main.go)

```go
package main

import (
	"log"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// if Allow DirectoryIndex
	// r.Use(static.Serve("/", static.LocalFile("./public", true)))
	// set prefix
	// r.Use(static.Serve("/static", static.LocalFile("./public", true)))

	r.Use(static.Serve("/", static.LocalFile("./public", false)))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

### Serve Embed folder

[embedmd]:# (example/embed/main.go)

```go
package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed public
var EmbedFS embed.FS

func main() {
	r := gin.Default()

	staticFiles, err := static.EmbedFolder(EmbedFS, "public/page")
	if err != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	} else {
		r.Use(static.Serve("/", staticFiles))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})

	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
```
