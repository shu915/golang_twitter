package controllers

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.HTML(200, "home/top", gin.H{})
}

func Index(c *gin.Context) {
	isAuthenticated := c.GetBool("isAuthenticated")
	c.HTML(200, "home/index", gin.H{
		"isAuthenticated": isAuthenticated,
	})
}