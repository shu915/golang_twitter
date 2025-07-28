package controllers

import (
	query "golang_twitter/db/query"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Router  *gin.Engine
	Queries *query.Queries
	RedisClient *redis.Client
}

func NewServer(queries *query.Queries, redisClient *redis.Client) *Server {
	return &Server{
		Router:  gin.Default(),
		Queries: queries,
		RedisClient: redisClient,
	}
}
