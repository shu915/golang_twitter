package controllers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	query "golang_twitter/db/query"
)

type Post struct {
	ID                 int32
	Content            string
	UserID             int32
	FormattedCreatedAt string
}

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

func (s *Server) Index(c *gin.Context) {
	currentPage, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		currentPage = 1
	}
	offset := (currentPage - 1) * 10

	postsCount, err := s.Queries.CountPosts(c.Request.Context())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	lastPage := int(math.Ceil(float64(postsCount) / 10))
	hasPreviousPage := currentPage > 1
	hasNextPage := currentPage < lastPage
	previousPage := currentPage - 1
	nextPage := currentPage + 1

	posts, err := s.Queries.GetPosts(c.Request.Context(), query.GetPostsParams{
		Limit:  int32(10),
		Offset: int32(offset),
	})
	if err != nil {
		c.HTML(200, "home/index", gin.H{"error": err.Error()})
		return
	}

	formattedPosts := make([]Post, len(posts))
	for i, post := range posts {
		var formattedTime string
		if post.CreatedAt.Valid {
			formattedTime = post.CreatedAt.Time.Format("2006-01-02 15:04:05")
		} else {
			formattedTime = "不明"
		}
		formattedPosts[i] = Post{
			ID:                 post.ID,
			Content:            post.Content,
			UserID:             post.UserID,
			FormattedCreatedAt: formattedTime,
		}
	}

	isAuthenticated := c.GetBool("isAuthenticated")
	c.HTML(200, "home/index", gin.H{
		"isAuthenticated": isAuthenticated,
		"csrf_token":      csrf.GetToken(c),
		"posts":           formattedPosts,
		"hasPreviousPage": hasPreviousPage,
		"hasNextPage":     hasNextPage,
		"currentPage":     currentPage,
		"lastPage":        lastPage,
		"previousPage":    previousPage,
		"nextPage":        nextPage,
	})
}
