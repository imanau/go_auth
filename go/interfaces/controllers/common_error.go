package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// SQLError return 500 and message
func SQLError(c echo.Context, message interface{}) {
	c.JSON(http.StatusInternalServerError, message)
}
