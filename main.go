package main

import (
	"os"
	"golang_twitter/db"
	"golang_twitter/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	db.Init()
	r := gin.Default()

	// セッション設定（CSRFに必要）
	store := cookie.NewStore([]byte(os.Getenv("CSRF_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	// CSRF保護ミドルウェア
	r.Use(csrf.Middleware(csrf.Options{
		Secret: "secret-csrf-key-change-in-production",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*")
	routes.InitRoutes(r)

	r.Run()
}
