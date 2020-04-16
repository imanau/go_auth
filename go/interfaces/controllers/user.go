package controllers

import (
	"go_auth/domain"
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

// Signup Handler
func Signup(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	// validation
	if user.UID == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid uid or password",
		}
	}
	// db connect
	db, err := model.ConnectDB()
	defer db.Close()
	if err != nil {
		SQLError(c, err)
	}
	if u := model.FindUser(db, user); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "uid already exists",
		}
	}

	model.CreateUser(db, user)
	user.Password = ""
	return c.JSON(http.StatusCreated, user)
}
