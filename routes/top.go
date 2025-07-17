package routes

import (
	"golang_twitter/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterTopRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)
}
