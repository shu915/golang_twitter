package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func (s *Server) Home(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err == nil {
		// セッションIDが存在する場合、Redisで確認
		userID, err := s.RedisClient.Get(c.Request.Context(), sessionID).Result()
		if err == nil && userID != "" {
			// ログインしていると判断し、/indexにリダイレクト
			c.Redirect(http.StatusFound, "/index")
			return
		}
	}

	c.HTML(200, "home/top", gin.H{})
}

func Index(c *gin.Context) {
	isAuthenticated := c.GetBool("isAuthenticated")
	c.HTML(200, "home/index", gin.H{
		"isAuthenticated": isAuthenticated,
		"csrf_token":      csrf.GetToken(c),
	})
}
