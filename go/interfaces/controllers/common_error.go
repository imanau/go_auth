package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

// SQLError return 500 and message
func SQLError(c echo.Context, message interface{}) {
	c.JSON(http.StatusInternalServerError, message)
}
