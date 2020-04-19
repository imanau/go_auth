package controllers

import (
	"encoding/json"
	"go_auth/domain"
	"go_auth/interfaces/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
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
	okJSON := `{"name":"ok","uid":"test@example.com","password": "password","role": 1}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/sign_up", strings.NewReader(okJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sign_up")
	// テスト前準備
	user := new(domain.User)
	json.Unmarshal([]byte(okJSON), &user)
	// db接続＆後処理
	db, err := model.ConnectDB()
	if err != nil {
		t.Error("db connection error")
	}
	defer db.Close()
	defer phisDelete(db, user)
	if assert.NoError(t, Signup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

// Loginの正常系
func TestLoginOk(t *testing.T) {
	// param pattern
	okJSON := `{"uid":"test@example.com","password": "password","name":"test"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(okJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// db準備
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		t.Error("db connection errro")
	}
	c.SetPath("/login")
	if assert.NoError(t, Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

// テスト用レコード物理削除関数
func phisDelete(db *gorm.DB, user *domain.User) {
	if user.ID == 0 {
		model.FindUser(db, user)
	}
	db.Unscoped().Delete(user)
}
