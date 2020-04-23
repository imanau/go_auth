package controllers

import (
	"encoding/json"
	"fmt"
	"go_auth/domain"
	"go_auth/interfaces/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	// テスト前準備＆データ削除等後処理
	user := new(domain.User)
	json.Unmarshal([]byte(okJSON), &user)
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
	c.SetPath("/login")
	// テストデータ用意　＆　後始末
	user := new(domain.User)
	json.Unmarshal([]byte(okJSON), &user)
	db, err := model.ConnectDB()
	defer db.Close()
	defer phisDelete(db, user)
	user.Password = "$2a$10$Oowv3K1NeSMj78lKv9mHLuNu.QBoFHjtZv5UvEMtljBLyAImixx5q"
	model.CreateUser(db, user)
	if err != nil {
		t.Error("db connection error")
	}
	if assert.NoError(t, Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

// admin/meの正常系テスト
func TestMeOK(t *testing.T) {
	// param pattern
	okJSON := `{"uid":"test@example.com","password": "password","name":"test"}`
	// token＆headerセット
	c, rec := jwtAuth("/admin/me", okJSON, "GET")
	exec := middleware.JWTWithConfig(Config)(UserIDFromToken)(c)
	// テストデータ用意　＆　後始末
	user := new(domain.User)
	json.Unmarshal([]byte(okJSON), &user)
	db, err := model.ConnectDB()
	defer db.Close()
	defer phisDelete(db, user)
	model.CreateUser(db, user)
	if err != nil {
		t.Error("db create error")
	}
	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

// UpdateUserの正常系
func TestUpdateOK(t *testing.T) {
	// param
	okJSON := `{"uid":"test@example.com","name":"test2","role":2}`
	// テストデータ用意　＆　後始末
	user := new(domain.User)
	json.Unmarshal([]byte(okJSON), &user)
	user.UID = "test@example.com"
	user.Name = "test"
	user.Password = "password"
	user.Role = 1
	db, err := model.ConnectDB()
	if err != nil {
		t.Error("db connect error")
	}
	defer db.Close()
	defer phisDelete(db, user)
	model.CreateUser(db, user)
	// token＆headerセット
	url := "/admin/users/" + strconv.Itoa(int(user.ID))
	c, rec := jwtAuth(url, okJSON, "PATCH")
	exec := middleware.JWTWithConfig(Config)(UserIDFromToken)(c)
	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

// changePassswordの正常系
func TestChangePasswordOK(t *testing.T) {
	// param
	okJSON := `{"new_password":"password2", "password": "password","password_confirmation":"password"}`
	// テストデータ用意　＆　後始末
	user := new(domain.User)
	json.Unmarshal([]byte(okJSON), &user)
	user.UID = "test@example.com"
	user.Name = "test"
	user.Password = "password"
	user.Role = 1
	db, err := model.ConnectDB()
	if err != nil {
		t.Error("db connect error")
	}
	defer db.Close()
	defer phisDelete(db, user)
	model.CreateUser(db, user)
	// token＆headerセット
	url := "/admin/users/change_password/" + strconv.Itoa(int(user.ID))
	c, rec := jwtAuth(url, okJSON, "PATCH")
	exec := middleware.JWTWithConfig(Config)(UserIDFromToken)(c)
	if assert.NoError(t, exec) {
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

// jwt認証共通部分
func jwtAuth(path string, param string, method string) (echo.Context, *httptest.ResponseRecorder) {
	user := domain.User{UID: "test@example.com", Name: "test", ID: uint(1)}
	token, err := CreateToken(&user)
	if err != nil {
		fmt.Print("jwt error")
	}
	e := echo.New()
	var req *http.Request
	switch method {
	case "GET":
		req = httptest.NewRequest(echo.GET, path, strings.NewReader(param))
	case "POST":
		req = httptest.NewRequest(echo.POST, path, strings.NewReader(param))
	case "PATCH":
		req = httptest.NewRequest(echo.PATCH, path, strings.NewReader(param))
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	return c, rec
}
