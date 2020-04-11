package main

import (
	"net/http"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// User User is model of users
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255);not null;"`
	UID       string `gorm:"type:varchar(255);not null;unique"`
	Pasword   string `gorm:"size:255;not null"`
	Role      int    `gorm:"not null"`
}

func main() {
	e := echo.New()

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// db接続とgorm準備
	db, err := gorm.Open("postgres", "host=db port=5432 user=postgres dbname=auth_service password=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("データベースへの接続に失敗しました")
	}
	defer db.Close()

	// スキーマのマイグレーション
	db.AutoMigrate(&User{})
	db.Create(&User{UID: "test@example.com", Pasword: "password", Role: 1})
	e.GET("/hello", mainPage())

	// Read
	var user User
	db.First(&user)
	fmt.Println(user)

	// サーバー起動
	e.Start(":3000") //ポート番号指定してね
}

func mainPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	}
}
