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
	r.Use(static.Serve("/", static.EmbedFolder(server, "data/server")))
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
