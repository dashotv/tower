package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *Application) MessagesIndex(c echo.Context) error {
	page, err := QueryDefaultInteger(c, "page", 1)
	if err != nil {
		return err
	}

	limit, err := QueryDefaultInteger(c, "limit", 250)
	if err != nil {
		return err
	}

	list, err := app.DB.Message.Query().Desc("created_at").Skip((page - 1) * limit).Limit(limit).Run()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, list)
}

func (a *Application) MessagesCreate(c echo.Context) error {
	m := &Message{}
	if err := c.Bind(m); err != nil {
		return err
	}

	if err := app.DB.Message.Save(m); err != nil {
		return err
	}

	if err := a.Events.Send("tower.logs", &EventLogs{Event: "new", Id: m.ID.Hex(), Log: m}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}
