package controllers

import (
	"context"
	"errors"
	query "golang_twitter/db/query"
	"golang_twitter/services"
	"golang_twitter/utils"
	"golang_twitter/validation"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) SignupPage(c *gin.Context) {
	c.HTML(200, "auth/signup", gin.H{
		"csrf_token": csrf.GetToken(c),
	})
}

func (s *Server) Signup(c *gin.Context) {
	var req validation.SignupRequest

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
			"errors": []validation.ValidationError{
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
		Token:    pgtype.Text{String: token, Valid: true},
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
		Token:    pgtype.Text{String: token, Valid: true},
	})

	c.HTML(http.StatusOK, "auth/activate_success", gin.H{})
}

func (s *Server) ActivateSuccessPage(c *gin.Context) {
	c.HTML(200, "auth/activate_success", gin.H{})
}

func (s *Server) ActivateErrorPage(c *gin.Context) {
	c.HTML(200, "auth/activate_error", gin.H{})
}

func (s *Server) LoginPage(c *gin.Context) {
	c.HTML(200, "auth/login", gin.H{
		"csrf_token": csrf.GetToken(c),
	})
}

func (s *Server) Login(c *gin.Context) {
	var req validation.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/login", gin.H{
			"error":      "リクエストの形式が正しくありません",
			"csrf_token": csrf.GetToken(c),
			"email":      req.Email,
		})
		return
	}

	if validationErrors := req.Validate(); validationErrors != nil {
		c.HTML(http.StatusBadRequest, "auth/login", gin.H{
			"errors":     validationErrors,
			"csrf_token": csrf.GetToken(c),
			"email":      req.Email,
		})
		return
	}

	user, err := s.Queries.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.HTML(http.StatusBadRequest, "auth/login", gin.H{
			"errors": []validation.ValidationError{
				{
					Field:   "email",
					Message: "メールアドレスまたはパスワードが間違っています",
				},
			},
			"csrf_token": csrf.GetToken(c),
			"email":      req.Email,
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.HTML(http.StatusBadRequest, "auth/login", gin.H{
			"errors": []validation.ValidationError{
				{
					Field:   "email",
					Message: "メールアドレスまたはパスワードが間違っています",
				},
			},
			"csrf_token": csrf.GetToken(c),
			"email":      req.Email,
		})
		return
	}
	if !user.IsActive.Bool {
		c.HTML(http.StatusBadRequest, "auth/login", gin.H{
			"errors": []validation.ValidationError{
				{
					Field:   "email",
					Message: "アクティベーションしてください",
				},
			},
			"csrf_token": csrf.GetToken(c),
			"email":      req.Email,
		})
		return
	}

	uuid := uuid.New().String()

	c.SetCookie("session_id", uuid, 360000, "/", "", false, true)
	err = s.RedisClient.Set(c.Request.Context(), uuid, user.ID, 3600*time.Second).Err()
	if err != nil {
		log.Printf("Redisセットエラー: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(200, "home/index", gin.H{
		"csrf_token": csrf.GetToken(c),
	})
}

func (s *Server) Logout(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	s.RedisClient.Del(c.Request.Context(), c.GetHeader("session_id"))
	c.Redirect(http.StatusFound, "/login")
}