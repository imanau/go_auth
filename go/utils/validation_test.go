package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// EmailValidateのテスト
func TestEmailValidate(t *testing.T) {
	okuser := EmailValidation{Email: "test@example.com"}
	nguser1 := EmailValidation{Email: "tesdaffadre"}
	nguser2 := EmailValidation{Email: ""}
	nguser3 := EmailValidation{Email: "@dafete.com"}
	res1 := EmailValidate(okuser)
	res2 := EmailValidate(nguser1)
	res3 := EmailValidate(nguser2)
	res4 := EmailValidate(nguser3)
	assert.Equal(t, res1, true)
	assert.Equal(t, res2, false)
	assert.Equal(t, res3, false)
	assert.Equal(t, res4, false)
}

// PasswordValidateのテスト
func TestPasswordValidate(t *testing.T) {
	okuser := PasswordValidation{Password: "password"}
	nguser1 := PasswordValidation{Password: ""}
	nguser2 := PasswordValidation{Password: "paaa"}
	res1 := PasswordValidate(okuser)
	res2 := PasswordValidate(nguser1)
	res3 := PasswordValidate(nguser2)
	assert.Equal(t, res1, true)
	assert.Equal(t, res2, false)
	assert.Equal(t, res3, false)
}

// UserValidateのテスト
func TestUserValidate(t *testing.T) {
	okuser := UserValidation{}
	okuser.Email = "test@example.com"
	okuser.Password = "drdfaLdar"
	nguser1 := UserValidation{}
	nguser1.Email = "@tesple.com"
	nguser1.Password = "drdfaLdar"
	nguser2 := UserValidation{}
	nguser2.Email = "test@example.com"
	nguser2.Password = "ldafag"
	res1 := UserValidate(okuser)
	res2 := UserValidate(nguser1)
	res3 := UserValidate(nguser2)
	assert.Equal(t, res1, true)
	assert.Equal(t, res2, false)
	assert.Equal(t, res3, false)
}
