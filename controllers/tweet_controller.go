package controllers

import (
	query "golang_twitter/db/query"
	validation "golang_twitter/validation"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func (s *Server) PostTweets(c *gin.Context) {
	var tweetRequest validation.TweetRequest

	if err := c.ShouldBind(&tweetRequest); err != nil {
		c.HTML(http.StatusBadRequest, "home/index", gin.H{
			"csrf_token": csrf.GetToken(c),
			"errors":     map[string]string{"general": "リクエストの形式が正しくありません"},
		})
		return
		}
		
		errors := tweetRequest.Validate()
		if len(errors) > 0 {
			c.Redirect(http.StatusFound, "/index")
			return
		}

	userID := c.GetInt("userID")

	_, err := s.Queries.CreatePost(c.Request.Context(), query.CreatePostParams{
		UserID:  int32(userID),
		Content: tweetRequest.Tweet,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.Redirect(http.StatusFound, "/index")
}
