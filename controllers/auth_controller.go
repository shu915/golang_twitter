package controllers

import (
	"context"
	query "golang_twitter/db/query"
	"golang_twitter/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

// AuthController 構造体
type AuthController struct {
	queries *query.Queries
}

// NewAuthController コンストラクタ
func NewAuthController(queries *query.Queries) *AuthController {
	return &AuthController{queries: queries}
}

func (ac *AuthController) SignupPage(c *gin.Context) {
	c.HTML(200, "auth/signup", gin.H{
		"csrf_token": csrf.GetToken(c),
	})
}

func (ac *AuthController) Signup(c *gin.Context) {
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
	_, err := ac.queries.GetUserByEmail(context.Background(), req.Email)
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

	// パスワードのハッシュ化
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	CreateUserParams := query.CreateUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	ac.queries.CreateUser(context.Background(), CreateUserParams)

	// 成功時は成功ページを表示
	c.HTML(http.StatusOK, "auth/signup_success", gin.H{})
}

func (ac *AuthController) SignupSuccessPage(c *gin.Context) {
	c.HTML(200, "auth/signup_success", gin.H{})
}
