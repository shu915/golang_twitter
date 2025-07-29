package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        sessionID, err := c.Cookie("session_id") // UUIDが入ってるクッキー名
        if err != nil {
           c.Redirect(http.StatusFound, "/login")
           c.Abort()
           return
        }

        userID, err := redisClient.Get(c.Request.Context(), sessionID).Result()
        if err != nil {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        c.Set("userID", userID)
        c.Set("isAuthenticated", true)
        c.Next()
    }
}