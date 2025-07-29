package routes

import (
	"golang_twitter/controllers"
	"golang_twitter/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, server *controllers.Server) {
	router.GET("/", controllers.Home)
	router.GET("/healthcheck", controllers.HealthCheck)
	router.GET("/signup", server.SignupPage)
	router.POST("/signup", server.Signup)
	router.GET("/signup_success", server.SignupSuccessPage)
	router.GET("/activate", server.Activate)
	router.GET("/activate_success", server.ActivateSuccessPage)
	router.GET("/login", server.LoginPage)
	router.POST("/login", server.Login)

	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(server.RedisClient))
	authRoutes.GET("/index", controllers.Index)
	authRoutes.GET("/logout", server.Logout)
}