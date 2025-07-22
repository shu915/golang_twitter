package main

import (
	"golang_twitter/db"
	"golang_twitter/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*")
	routes.InitRoutes(r)

	r.Run()
}
