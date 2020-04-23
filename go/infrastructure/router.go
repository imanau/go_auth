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
	api := e.Group("/admin")
	api.Use(middleware.JWTWithConfig(controllers.Config))
	api.GET("/users", controllers.Index)
	api.POST("/users", controllers.Signup)
	api.PATCH("/users/:id", controllers.UpdateUser)
	api.PATCH("/users/change_password/:id", controllers.ChangePassword)
	api.GET("/me", controllers.UserIDFromToken)
	e.Start(":3000")
}
