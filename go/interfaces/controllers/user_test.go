package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

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

// SignUpの正常系
func TestSignupOk(t *testing.T) {
	// param pattern
	okJSON := `{"name":"ok","uid":"test2@example.com","password": "password","role": 1}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/sign_up", strings.NewReader(okJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sign_up")
	if assert.NoError(t, Signup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

// Loginの正常系
func TestLoginOk(t *testing.T) {
	// param pattern
	okJSON := `{"uid":"test2@example.com","password": "password"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(okJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")
	if assert.NoError(t, Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
