package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        sessionID, err := c.Cookie("session_id") // UUIDが入ってるクッキー名
        if err != nil {
            log.Printf("Error retrieving user ID: %v", err)
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        userIDStr, err := redisClient.Get(c.Request.Context(), sessionID).Result()
        if err != nil {
            log.Printf("Error retrieving user ID: %v", err)
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            log.Printf("Error converting user ID to int: %v", err)
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        //intでセットする
        c.Set("userID", userID)
        c.Set("isAuthenticated", true)
        c.Next()
    }
}