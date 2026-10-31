package main

import (
	"embed"
	"log"
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
		log.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
