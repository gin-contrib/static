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
