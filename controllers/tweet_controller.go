package controllers

import (
	query "golang_twitter/db/query"
	validation "golang_twitter/validation"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func (s *Server) PostTweets(c *gin.Context) {
	var tweetRequest validation.TweetRequest

	if err := c.ShouldBind(&tweetRequest); err != nil {
		errors := tweetRequest.Validate()
		if len(errors) > 0 {
			c.HTML(200, "home/index", gin.H{
				"csrf_token": csrf.GetToken(c),
				"errors":     errors,
			})
			return
		}
	}

	userIDStr, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}
	_, err = s.Queries.CreatePost(c.Request.Context(), query.CreatePostParams{
		UserID:  int32(userID),
		Content: tweetRequest.Tweet,
	})

	if err != nil {
		log.Printf("Failed to create post: %v", err)
		return
	}

	c.HTML(http.StatusOK, "home/index", gin.H{
		"csrf_token":      csrf.GetToken(c),
		"isAuthenticated": true,
		"success":         "ツイートしました",
	})
}
