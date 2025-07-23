package main

import (
	"golang_twitter/controllers"
	"golang_twitter/db"
	query "golang_twitter/db/query"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	dbPool, err := db.Init()
	if err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	defer dbPool.Close()

	queries := query.New(dbPool)
	server := controllers.NewServer(queries) // ここでServer構造体を使う

	store := cookie.NewStore([]byte(os.Getenv("CSRF_SECRET")))
	server.Router.Use(sessions.Sessions("mysession", store))
	server.Router.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CSRF_SECRET"),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	server.Router.Static("/static", "./static")
	server.Router.LoadHTMLGlob("templates/**/*")

	server.Router.GET("/", controllers.Home)
	server.Router.GET("/healthcheck", controllers.HealthCheck)
	server.Router.GET("/signup", server.SignupPage)
	server.Router.POST("/signup", server.Signup)
	server.Router.GET("/signup_success", server.SignupSuccessPage)

	// サーバーを起動
	server.Router.Run()
}
