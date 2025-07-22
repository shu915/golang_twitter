package routes

import (
	"golang_twitter/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, authController *controllers.AuthController) {
	RegisterTopRoutes(r)
	RegisterHealthRoutes(r)
	RegisterAuthRoutes(r, authController)
}
