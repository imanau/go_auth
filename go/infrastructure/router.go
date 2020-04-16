package infrastructure

import (
	"go_auth/interfaces/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Init Start webserver
func Init() {
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// router
	e.GET("/", controllers.Index)
	e.POST("/sign_up", controllers.Signup)
	e.POST("/login", controllers.Login)
	e.Start(":3000")
}
