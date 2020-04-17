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
	e.POST("/sign_up", controllers.Signup)
	e.POST("/login", controllers.Login)
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(controllers.Config))
	api.GET("/users", controllers.Index)
	e.Start(":3000")
}
