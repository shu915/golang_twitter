package controllers

import (
	"context"
	"errors"
	query "golang_twitter/db/query"
	"golang_twitter/dto"
	"golang_twitter/services"
	"golang_twitter/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) SignupPage(c *gin.Context) {
	c.HTML(200, "auth/signup", gin.H{
		"csrf_token": csrf.GetToken(c),
	})
}

func (s *Server) Signup(c *gin.Context) {
	var req dto.SignupRequest

	// HTMLフォームからのデータをバインド
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/signup", gin.H{
			"error":      "リクエストの形式が正しくありません",
			"csrf_token": csrf.GetToken(c),
		})
		return
	}

	// DTOのカスタムバリデーションを使用
	if validationErrors := req.Validate(); validationErrors != nil {
		c.HTML(http.StatusBadRequest, "auth/signup", gin.H{
			"errors":     validationErrors,
			"email":      req.Email, // 入力値を保持
			"csrf_token": csrf.GetToken(c),
		})
		return
	}

	// Emailのユニーク制約チェック
	_, err := s.Queries.GetUserByEmail(context.Background(), req.Email)
	if err == nil {
		c.HTML(http.StatusBadRequest, "auth/signup", gin.H{
			"errors": []dto.ValidationError{
				{
					Field:   "email",
					Message: "このメールアドレスは既に使用されています",
				},
			},
			"email":      req.Email,
			"csrf_token": csrf.GetToken(c),
		})
		return
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("DBエラー: %v", err)
    c.AbortWithStatus(http.StatusInternalServerError)
    return
}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("パスワードのハッシュ化エラー: %v", err)
    c.AbortWithStatus(http.StatusInternalServerError)
    return
}

	token, err := utils.GenerateToken(32)
	if err != nil {
		log.Printf("トークンの生成エラー: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	CreateUserParams := query.CreateUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
		Token: pgtype.Text{String: token, Valid: true},
	}

	_, err = s.Queries.CreateUser(c.Request.Context(), CreateUserParams)
	if err != nil {
		log.Printf("ユーザーの作成エラー: %v", err)
    c.AbortWithStatus(http.StatusInternalServerError)
    return
}

	err = services.SendActivationEmail(req.Email, token)
	if err != nil {
		log.Printf("メールの送信エラー: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 成功時は成功ページを表示
	c.HTML(http.StatusOK, "auth/signup_success", gin.H{})
}

func (s *Server) SignupSuccessPage(c *gin.Context) {
	c.HTML(200, "auth/signup_success", gin.H{})
}

func (s *Server) Activate(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.HTML(http.StatusBadRequest, "auth/activate_error", gin.H{
			"error": "トークンが無効です",
		})
		return
	}

	_, err := s.Queries.GetUserByToken(c.Request.Context(), pgtype.Text{String: token, Valid: true})
	if err != nil {
		c.HTML(http.StatusBadRequest, "auth/activate_error", gin.H{
			"error": "トークンが無効です",
		})
		return
	}

	s.Queries.UpdateUserIsActive(c.Request.Context(), query.UpdateUserIsActiveParams{
		IsActive: pgtype.Bool{Bool: true, Valid: true},
		Token: pgtype.Text{String: token, Valid: true},
	})



	c.HTML(http.StatusOK, "auth/activate_success", gin.H{})
}

func (s *Server) ActivateSuccessPage(c *gin.Context) {
	c.HTML(200, "auth/activate_success", gin.H{})
}

func (s *Server) ActivateErrorPage(c *gin.Context) {
	c.HTML(200, "auth/activate_error", gin.H{})
}