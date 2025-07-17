package controllers

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.HTML(200, "home/top", gin.H{})
}