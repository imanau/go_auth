package model

import (
	"go_auth/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllUser(t *testing.T) {
	result := AllUser()
	if result.Error == nil {
		assert.IsType(t, result.Value, new(*domain.Users))
	}
}
