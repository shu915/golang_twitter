package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type SignupRequest struct {
	Email           string `validate:"required,email,max=255" form:"email"`
	Password        string `validate:"required,min=8,max=255" form:"password"`
	PasswordConfirm string `validate:"required,eqfield=Password" form:"password_confirmation"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `validate:"required,email,max=255" form:"email"`
	Password string `validate:"required,min=8,max=255" form:"password"`
}

func handleBasicValidation(req interface{}) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			message := getErrorMessage(err)
			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}
	return errors
}
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate バリデーションを実行
func (s *SignupRequest) Validate() []ValidationError {
	var errors []ValidationError

	// 基本バリデーション
	errors = handleBasicValidation(s)

	// カスタムパスワードバリデーション
	if s.Password != "" { // パスワードが空でない場合のみチェック
		passwordErrors := validatePasswordComplexity(s.Password)
		errors = append(errors, passwordErrors...)
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// validatePasswordComplexity パスワードの複雑性をチェック
func validatePasswordComplexity(password string) []ValidationError {
	var errors []ValidationError

	// 1. 半角英数字が含まれる（英字だけはNG、数字だけもNG）
	hasAlpha := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasAlpha || !hasDigit {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "パスワードには英字と数字の両方を含む必要があります",
		})
	}

	// 2. 英字は小文字大文字混合（小文字だけはNG）
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)

	if hasAlpha && (!hasLower || !hasUpper) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "パスワードには大文字と小文字の両方を含む必要があります",
		})
	}

	// 3. 指定の記号が1文字以上含まれる (!?-_)
	hasSymbol := regexp.MustCompile(`[!?_-]`).MatchString(password)

	if !hasSymbol {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "パスワードには記号（!?-_）のいずれかを含む必要があります",
		})
	}

	return errors
}

func (l *LoginRequest) Validate() []ValidationError {
	errors := handleBasicValidation(l)

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// getErrorMessage バリデーションエラーメッセージを日本語で返す
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%sは必須項目です", getFieldNameJP(err.Field()))
	case "email":
		return "有効なメールアドレスを入力してください"
	case "min":
		return fmt.Sprintf("%sは%s文字以上で入力してください", getFieldNameJP(err.Field()), err.Param())
	case "max":
		return fmt.Sprintf("%sは%s文字以下で入力してください", getFieldNameJP(err.Field()), err.Param())
	case "eqfield":
		return "パスワードが一致しません"
	default:
		return fmt.Sprintf("%sの形式が正しくありません", getFieldNameJP(err.Field()))
	}
}

// getFieldNameJP フィールド名を日本語に変換
func getFieldNameJP(field string) string {
	switch field {
	case "Email":
		return "メールアドレス"
	case "Password":
		return "パスワード"
	case "PasswordConfirm":
		return "パスワード（確認用）"
	default:
		return field
	}
}