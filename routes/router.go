package routes

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.Engine) {
	RegisterTopRoutes(r)
	RegisterHealthRoutes(r)
	RegisterAuthRoutes(r)
}