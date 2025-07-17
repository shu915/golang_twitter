package routes

import "github.com/gin-gonic/gin"

func RegisterAuthRoutes(r *gin.Engine) {
	r.GET("/signup", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Signup"})
	})
}