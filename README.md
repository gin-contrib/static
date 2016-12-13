# static middleware

[![Build Status](https://travis-ci.org/gin-contrib/static.svg)](https://travis-ci.org/gin-contrib/static)
[![codecov](https://codecov.io/gh/gin-contrib/static/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/static)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/static)](https://goreportcard.com/report/github.com/gin-contrib/static)
[![GoDoc](https://godoc.org/github.com/gin-contrib/static?status.svg)](https://godoc.org/github.com/gin-contrib/static)

Static middleware

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/gin-contrib/static
```

Import it in your code:

```go
import "github.com/gin-contrib/static"
```

### Canonical example:

See the [example](example)

```go
package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	r := gin.Default()

	store := persistence.NewInMemoryStore(time.Second)
	// Cached Page
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	r.GET("/cache_ping", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	}))

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
```
