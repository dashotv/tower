package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /users/
func (a *Application) UsersIndex(c echo.Context) error {
	list, err := a.DB.User.Query().Run()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Users"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list})
}
