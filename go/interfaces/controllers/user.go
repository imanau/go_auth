// Package controllers パッケージはMVCモデルのコントローラーにあたります
//
// ルーティングに紐付けられたアクションを実行します。
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

// jwt認証用の型です。jwt認証のclaimに保持する値を要素に持ちます。
type jwtCustomClaims struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

// jwt認証に利用する文字列です。（注：今回はサンプル文字列を利用しています。）
var signingKey = []byte("secretkeysample")

// Config jwt認証用の変数です。
var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

// Signup ユーザー作成用のハンドラです。
// domain.User型に準拠したjsonパラメーターを送信することでdbと連動し新規ユーザーを作成します。
// return json (new domain.User)
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

// Login jwt認証を用いないユーザーの認証を行います。
// domain.User型に準拠したjsonパラメーターを送信することでdbと連動し既存ユーザーを検索します。
// return json (domain.User)
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

// Index ユーザーの一覧情報をjsonで返却します。
// domain.User型に準拠したjsonパラメーターを送信することでdbと連動し既存ユーザーズを検索します。
// return json (array[domain.User])
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

// Show ユーザー情報をjsonで返却します。
// urlパラメーターのid情報に合致するユーザー情報を返却します。
// return json (domain.User)
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

// UpdateUser ユーザー情報を更新し、jsonで更新後ユーザー情報を返却します。
// domain.User型に準拠したjsonパラメーターを送信することでdbと連動し既存ユーザー情報を更新します。
// return json (domain.User)
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

// ChangePassword ユーザーのパスワードを変更します。
// domain.PasswordInfo型に準拠したjsonパラメーターを送信することでdbと連動し既存ユーザー情報を更新します。
// return json (domain.User)
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

// InitPassword ユーザーのパスワードを初期化します。
// domain.User型に準拠したjsonパラメーターを送信することでdbと連動し既存ユーザー情報を更新します。
// return json (domain.User)
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

// DestroyUser ユーザー情報を論理削除します。
// domain.User型に準拠したjsonパラメーターを送信することでdbと連動し既存ユーザー情報のdeleted_atを更新します。
// deleted_atが空ではないユーザーは、IndexやShowアクションで返却されません。
// return json (domain.User)
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

// UserIDFromToken jwt tokenから得たユーザー情報をjsonで返却します。
// request headerに付与されたtoken情報からユーザー情報を抽出し、domain.User型のjsonを返却します。
// return json (domain.User)
func UserIDFromToken(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*jwtCustomClaims)
	id := claims.ID
	name := claims.Name
	user := domain.UserForGeneral{ID: id, Name: name}
	return c.JSON(http.StatusOK, user)
}

// CreateToken *domain.User型を渡すことで、その情報からjwt tokenを生成します。
// ハンドラではなく、token情報の文字列及びerror情報を呼び出し元に返却します。
// return string t, err errror
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

// AdminAuthMiddleware 管理者権限用認可middle ware
// admin/以下のurlのアクションの前に呼び出され、token情報からユーザーの権限をチェックします。
// return echo.Context（各ハンドラーに引き渡す）
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
		c.Set("authorizedUser", user)
		return next(c)
	}
}

// UserAuthMiddleware 一般ユーザー権限用認可middle ware
// api/以下のurlのアクションの前に呼び出され、token情報からユーザーの権限をチェックします。
// 一般ユーザーは自身以外のユーザー情報にアクセスできないように制御します。
// return echo.Context（各ハンドラーに引き渡す）
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
		c.Set("authorizedUser", user)
		return next(c)
	}
}

// SetUser Show, Updateアクションの前に操作対象のユーザー情報を検索し、返却します。
// urlパラメータに指定するユーザーidからユーザー情報を検索します。
// return echo.Context（各ハンドラーに引き渡す）
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
