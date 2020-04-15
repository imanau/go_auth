package controllers

import (
	"go_auth/interfaces/model"
	"net/http"

	"github.com/labstack/echo"
)

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
	return c.JSON(http.StatusOK, rows)
}
