package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// "go_auth/domain"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// var (
//    testMockDB = domain.Users{
//        "ID": 1,
//        "Email": "test@example.com",
//        "Password": "testest"
//    }
// )

func TestIndex(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/")
	if assert.NoError(t, Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
