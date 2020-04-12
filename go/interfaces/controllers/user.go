package controllers

import (
	"go_auth/interfaces/model"
	"net/http"

	"github.com/labstack/echo"
)

// Index indexActionHandler
func Index() echo.HandlerFunc {
	return func(c echo.Context) error {
		rows := model.AllUser()
		if rows.Error != nil {
			c.JSON(http.StatusInternalServerError, "SQL Error")
		}
		return c.JSON(http.StatusOK, rows)
	}
}
