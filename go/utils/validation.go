package utils

import (
	"gopkg.in/go-playground/validator.v9"
)

// ValidateUser struct of validation user
type ValidateUser struct {
	Email    string `validate:"required,email"` //必須パラメータ、かつ、emailフォーマット
	Password string `validate:"required,min=8,max=255"`
}

// UserValidation パスワード及びemailアドレスのバリデーション
func UserValidation(user ValidateUser) (result bool) {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		result = false
	} else {
		result = true
	}
	return result
}
