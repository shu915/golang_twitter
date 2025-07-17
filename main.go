package main

import (
	"golang_twitter/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*")
	routes.InitRoutes(r)

	r.Run()
}
