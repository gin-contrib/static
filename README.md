# static middleware

[![Run Tests](https://github.com/gin-contrib/static/actions/workflows/go.yml/badge.svg)](https://github.com/gin-contrib/static/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gin-contrib/static/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/static)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/static)](https://goreportcard.com/report/github.com/gin-contrib/static)
[![GoDoc](https://godoc.org/github.com/gin-contrib/static?status.svg)](https://godoc.org/github.com/gin-contrib/static)

Static middleware

## Usage

### Start using it

Download and install it:

```sh
go get github.com/gin-contrib/static
```

Import it in your code:

```go
import "github.com/gin-contrib/static"
```

### Canonical example

See the [example](_example)

#### Serve local file

```go
package main

import (
  "embed"
  "fmt"
  "net/http"

  "github.com/gin-contrib/static"
  "github.com/gin-gonic/gin"
)

//go:embed data
var server embed.FS

func main() {
  r := gin.Default()
  fs, err := static.EmbedFolder(server, "data/server")
  if err != nil {
    panic(err)
  }
  r.Use(static.Serve("/", fs))
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
