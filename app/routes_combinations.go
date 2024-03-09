package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func (a *Application) CombinationsIndex(c echo.Context, page int, limit int) error {
	list, err := a.DB.CombinationList(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (a *Application) CombinationsCreate(c echo.Context) error {
	combination := &Combination{}
	if err := c.Bind(combination); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
	}

	if err := a.DB.Combination.Save(combination); err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
	}

	return c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) CombinationsShow(c echo.Context, name string) error {
	list, err := a.DB.CombinationChildren(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (a *Application) CombinationsUpdate(c echo.Context, id string) error {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Combinations.Get(id)
	return c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) CombinationsSettings(c echo.Context, id string) error {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Combinations.Get(id)
	return c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) CombinationsDelete(c echo.Context, id string) error {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Combinations.Get(id)
	return c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}
