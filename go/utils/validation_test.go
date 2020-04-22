package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
