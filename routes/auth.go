package routes

import (
	"golang_twitter/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, authController *controllers.AuthController) {
	r.GET("/signup", authController.SignupPage)
	r.POST("/signup", authController.Signup)
	r.GET("/signup/success", authController.SignupSuccessPage)
}
