package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /messages/
func (a *Application) MessagesIndex(c echo.Context, page int, limit int) error {
	list, total, err := a.DB.MessageList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error loading Messages"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: list, Total: total})
}

// POST /messages/
func (a *Application) MessagesCreate(c echo.Context, subject *Message) error {
	if err := a.DB.Message.Save(subject); err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{Error: true, Message: "error saving Messages"})
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: subject})
}
