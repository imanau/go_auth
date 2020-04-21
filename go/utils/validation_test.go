package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// UserValidationのテスト
func TestUserValidation(t *testing.T) {
	okuser := ValidateUser{Email: "test@example.com", Password: "drdfaLdar"}
	nguser1 := ValidateUser{Email: "@tesple.com", Password: "drdfaLdar"}
	nguser2 := ValidateUser{Email: "test@example.com", Password: "ldafag"}
	res1 := UserValidation(okuser)
	res2 := UserValidation(nguser1)
	res3 := UserValidation(nguser2)
	assert.Equal(t, res1, true)
	assert.Equal(t, res2, false)
	assert.Equal(t, res3, false)
}
