package routes

import (
	"golang_twitter/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(r *gin.Engine) {
	r.GET("/healthcheck", controllers.HealthCheck)
}