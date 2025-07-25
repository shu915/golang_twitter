package controllers

import (
	query "golang_twitter/db/query"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router  *gin.Engine
	Queries *query.Queries
}

func NewServer(queries *query.Queries) *Server {
	return &Server{
		Router:  gin.Default(),
		Queries: queries,
	}
}
