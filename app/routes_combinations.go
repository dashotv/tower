package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) CombinationsIndex(c *gin.Context, page int, limit int) {
	list, err := a.DB.CombinationList(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (a *Application) CombinationsCreate(c *gin.Context) {
	combination := &Combination{}
	if err := c.BindJSON(combination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if err := a.DB.Combination.Save(combination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) CombinationsShow(c *gin.Context, name string) {
	list, err := a.DB.CombinationChildren(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (a *Application) CombinationsUpdate(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Combinations.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) CombinationsSettings(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Combinations.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) CombinationsDelete(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Combinations.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}
