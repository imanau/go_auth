package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// PasswordValidtionのテスト
func TestPasswordValidtion(t *testing.T) {
	okStr := "Dfaeytgd"
	ngStr := "1234567"
	res1 := PasswordValidtion(okStr)
	res2 := PasswordValidtion(ngStr)
	assert.Equal(t, res1, true)
	assert.Equal(t, res2, false)
}
