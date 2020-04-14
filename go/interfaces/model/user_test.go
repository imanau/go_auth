package model

import (
	"go_auth/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllUser(t *testing.T) {
	result := AllUser()
	if result.Error == nil {
		// レコードがある場合はvalueがdomain.Users型
		assert.IsType(t, result.Value, new(*domain.Users))
	} else {
		// エラーが有った場合はレコードが見つからないエラーである
		bool := result.RecordNotFound()
		assert.True(t, bool)
	}
}
