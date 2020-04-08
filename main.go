package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", mainPage())

	// サーバー起動
	e.Start(":3000") //ポート番号指定してね
}

func mainPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	}
}
