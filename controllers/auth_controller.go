package controllers

import "github.com/gin-gonic/gin"

func SignupPage(c *gin.Context) {
	c.HTML(200, "auth/signup", gin.H{})
}

func Signup(c *gin.Context) {
	
}