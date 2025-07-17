package routes

import (
	"golang_twitter/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	r.GET("/signup", controllers.SignupPage)
	r.POST("/signup", controllers.Signup)
}