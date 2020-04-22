package utils

import (
	"gopkg.in/go-playground/validator.v9"
)

// EmailValidation struct of validation email
type EmailValidation struct {
	Email string `validate:"required,email"` //必須パラメータ、かつ、emailフォーマット
}

// PasswordValidation struct of validation password
type PasswordValidation struct {
	Password string `validate:"required,min=8,max=255"`
}

// UserValidation struct of validation user
type UserValidation struct {
	EmailValidation
	PasswordValidation
}

// UserValidate パスワード及びemailアドレスのバリデーション
func UserValidate(user UserValidation) (result bool) {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		result = false
	} else {
		result = true
	}
	return result
}

// EmailValidate emailアドレスのバリデーション
func EmailValidate(user EmailValidation) (result bool) {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		result = false
	} else {
		result = true
	}
	return result
}
