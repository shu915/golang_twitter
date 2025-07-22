package main

import (
	"golang_twitter/controllers"
	"golang_twitter/db"
	query "golang_twitter/db/query"
	"golang_twitter/routes"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	// DB接続を初期化
	dbPool, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()

	// Queriesを初期化
	queries := query.New(dbPool)

	// AuthControllerを初期化
	authController := controllers.NewAuthController(queries)

	r := gin.Default()

	// セッション設定（CSRFに必要）
	store := cookie.NewStore([]byte(os.Getenv("CSRF_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	// CSRF保護ミドルウェア
	r.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CSRF_SECRET"),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*")
	routes.InitRoutes(r, authController)

	r.Run()
}
