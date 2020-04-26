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

	e.POST("/sign_up", controllers.Signup)
	e.POST("/login", controllers.Login)
	admin := e.Group("/admin")
	admin.Use(middleware.JWTWithConfig(controllers.Config))
	admin.Use(controllers.AdminAuthMiddleware)
	admin.GET("/users", controllers.Index)
	admin.POST("/users", controllers.Signup)
	admin.GET("/users/:id", controllers.Show)
	admin.PATCH("/users/:id", controllers.UpdateUser)
	admin.DELETE("/users/:id", controllers.DestroyUser)
	admin.PATCH("/users/change_password/:id", controllers.ChangePassword)
	admin.PATCH("/users/init_account/:id", controllers.InitPassword)
	admin.GET("/me", controllers.UserIDFromToken)
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(controllers.Config))
	api.Use(controllers.UserAuthMiddleware)
	api.GET("/users/:id", controllers.Show)
	e.Start(":3000")
}
