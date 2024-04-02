package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func (a *Application) FeedsIndex(c echo.Context, page, limit int) error {
	results, err := app.DB.Feed.Query().
		Desc("processed").
		Limit(1000).
		Run()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}

func (a *Application) FeedsCreate(c echo.Context, data *Feed) error {
	err := app.DB.Feed.Save(data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "id": data.ID.Hex(), "result": data})
}

func (a *Application) FeedsShow(c echo.Context, id string) error {
	result := &Feed{}
	err := app.DB.Feed.Find(id, result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (a *Application) FeedsUpdate(c echo.Context, id string, data *Feed) error {
	err := app.DB.FeedUpdate(id, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "result": data})
}

func (a *Application) FeedsSettings(c echo.Context, id string, data *Setting) error {
	err := app.DB.FeedSetting(id, data.Name, data.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false, "result": data})
}

func (a *Application) FeedsDelete(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, gin.H{"error": false})
}
