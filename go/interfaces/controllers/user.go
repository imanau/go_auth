package controllers

import (
	"go_auth/domain"
	"go_auth/interfaces/model"
	"go_auth/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var signingKey = []byte("secretkeysample")

// Config jwtconfig
var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

// Signup Handler
func Signup(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	validateUser := utils.UserValidation{}
	validateUser.Email = user.UID
	validateUser.Password = user.Password
	// validation
	if !utils.UserValidate(validateUser) {
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
	t, err := CreateToken(u)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
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

// Show showActionHandler
func Show(c echo.Context) error {
	user := new(domain.User)
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	SetUser(user, db, c)
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// UpdateUser パスワード以外のユーザー情報の更新
func UpdateUser(c echo.Context) error {
	user := new(domain.User)
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	SetUser(user, db, c)
	if err = c.Bind(user); err != nil {
		return err
	}
	// validation
	validateUser := utils.EmailValidation{Email: user.UID}
	if !utils.EmailValidate(validateUser) {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid uid",
		}
	}
	// ユーザー更新処理
	user.Password = ""
	model.UpdateUser(db, user)
	return c.JSON(http.StatusOK, user)
}

// ChangePassword ユーザーのパスワード更新
func ChangePassword(c echo.Context) error {
	user := new(domain.User)
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	SetUser(user, db, c)
	passwordInfo := new(domain.PasswordInfo)
	if err = c.Bind(passwordInfo); err != nil {
		return err
	}
	// validation
	validateUser := utils.PasswordValidation{Password: passwordInfo.Password}
	if !utils.PasswordValidate(validateUser) || passwordInfo.Password != passwordInfo.PasswordConfirmation {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid uid",
		}
	}
	// パスワード暗号化処理
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordInfo.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	user.Password = string(hash)
	// ユーザー更新処理
	model.ChangePassword(db, user)
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// InitPassword ユーザーのパスワード初期化
func InitPassword(c echo.Context) error {
	user := new(domain.User)
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	SetUser(user, db, c)
	// パスワード暗号化処理
	hash, err := bcrypt.GenerateFromPassword([]byte(user.UID), bcrypt.DefaultCost)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	user.Password = string(hash)
	model.ChangePassword(db, user)
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// DestroyUser ユーザーの削除処理
func DestroyUser(c echo.Context) error {
	user := new(domain.User)
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	SetUser(user, db, c)
	model.DeleteUser(db, user)
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// UserIDFromToken jwt tokenでユーザーを認証し、そのユーザー情報を返却する
func UserIDFromToken(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*jwtCustomClaims)
	id := claims.ID
	name := claims.Name
	user := domain.UserForGeneral{ID: id, Name: name}
	return c.JSON(http.StatusOK, user)
}

// CreateToken jwt tokenを作成する処理
func CreateToken(u *domain.User) (string, error) {
	claims := &jwtCustomClaims{
		u.ID,
		u.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	return t, err
}

// AdminAuthMiddleware admin権限があるか判断するカスタムミドルウェア
func AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(domain.User)
		// jwt tokenからユーザーidを取得
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*jwtCustomClaims)
		id := claims.ID
		user.ID = id
		// db connect
		db, err := model.ConnectDB()
		defer db.Close()
		if err != nil {
			SQLError(c, err)
		}
		// ユーザー検証
		model.FindUser(db, user)
		if user.ID == 0 || user.Role != 1 {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "アクセス権限がありません",
			}
		}
		c.Set("authorizedUser")
		return next(c)
	}
}

// UserAuthMiddleware jwt認証後の一般ユーザー用middleware
func UserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(domain.User)
		// jwt tokenからユーザーidを取得
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*jwtCustomClaims)
		user.ID = claims.ID
		// request bodyからidを取得
		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "invalid id",
			}
		}
		if user.ID != uint(id) {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "アクセス権限がありません",
			}
		}
		c.Set("authorizedUser")
		return next(c)
	}
}

// SetUser ユーザー情報をセットする
func SetUser(user *domain.User, db *gorm.DB, c echo.Context) error {
	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}
	user.ID = uint(id)
	if err != nil {
		SQLError(c, err)
	}
	// ユーザー検証
	model.FindUser(db, user)
	if user.ID == 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "ユーザー情報が正しくありません",
		}
	}
	return err
}
