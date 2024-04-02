package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func (a *Application) RequestsIndex(c echo.Context, page, limit int) error {
	list, err := app.DB.Request.Query().Desc("created_at").Run()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, list)
}

func (a *Application) RequestsShow(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, gin.H{"message": "RequestsShow"})
}

func (a *Application) RequestsCreate(c echo.Context, r *Request) error {
	return c.JSON(http.StatusOK, gin.H{"message": "RequestsCreate"})
}

func (a *Application) RequestsSettings(c echo.Context, id string, s *Setting) error {
	return c.JSON(http.StatusOK, gin.H{"message": "RequestsSettings"})
}

func (a *Application) RequestsDelete(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, gin.H{"message": "RequestsDelete"})
}

func (a *Application) RequestsUpdate(c echo.Context, id string, updated *Request) error {
	req := &Request{}
	err := app.DB.Request.Find(id, req)
	if err != nil {
		return err
	}

	req.Status = updated.Status
	if err := app.DB.Request.Update(req); err != nil {
		return err
	}

	if updated.Status == "approved" {
		if err := app.Workers.Enqueue(&CreateMediaFromRequests{}); err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, req)
}
