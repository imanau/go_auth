package controllers

import (
	"go_auth/domain"
	"go_auth/interfaces/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	ID   uint   `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var signingKey = []byte("secretkeysample")

// Config jwtconfig
var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

// Index indexActionHandler
func Index(c echo.Context) error {
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	rows := model.AllUser(db)
	if rows.Error != nil {
		SQLError(c, rows.Error)
	}
	return c.JSON(http.StatusOK, rows.Value)
}

// Signup Handler
func Signup(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	// validation
	if user.UID == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid uid or password",
		}
	}
	// db connect
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	if u := model.FindUser(db, user); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "uid already exists",
		}
	}
	// パスワード暗号化処理
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	user.Password = string(hash)
	model.CreateUser(db, user)
	user.Password = ""
	return c.JSON(http.StatusCreated, user)
}

// Login Handler return jwt
func Login(c echo.Context) error {
	u := new(domain.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	inputPass := u.Password
	// db connect
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	model.FindUser(db, u)
	if u.ID == 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "emailが正しくありません",
		}
	}
	// ハッシュ化したパスワードの比較
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPass))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "パスワードが正しくありません",
		}
	}
	claims := &jwtCustomClaims{
		u.ID,
		u.UID,
		u.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

// UserIDFromToken jwt tokenでユーザーを認証し、そのユーザー情報を返却する
func UserIDFromToken(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*jwtCustomClaims)
	id := claims.ID
	uid := claims.UID
	name := claims.Name
	user := domain.UserForGeneral{ID: id, UID: uid, Name: name}
	return c.JSON(http.StatusOK, user)
}
