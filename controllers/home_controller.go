package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func Home(c *gin.Context) {
	c.HTML(200, "home/top", gin.H{})
}

func Index(c *gin.Context) {
	isAuthenticated := c.GetBool("isAuthenticated")
	c.HTML(200, "home/index", gin.H{
		"isAuthenticated": isAuthenticated,
		"csrf_token":      csrf.GetToken(c),
	})
}
