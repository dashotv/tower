package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *Application) UsersIndex(c echo.Context) error {
	users, err := app.DB.User.Query().Run()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}
