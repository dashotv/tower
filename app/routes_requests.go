package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) RequestsIndex(c *gin.Context, page, limit int) {
	list, err := app.DB.Request.Query().Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *Application) RequestsShow(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"message": "RequestsShow"})
}

func (a *Application) RequestsCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "RequestsCreate"})
}

func (a *Application) RequestsSettings(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"message": "RequestsSettings"})
}

func (a *Application) RequestsDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"message": "RequestsDelete"})
}

func (a *Application) RequestsUpdate(c *gin.Context, id string) {
	req := &Request{}
	err := app.DB.Request.Find(id, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated := &Request{}
	if err := c.BindJSON(updated); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Status = updated.Status
	if err := app.DB.Request.Update(req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updated.Status == "approved" {
		if err := app.Workers.Enqueue(&CreateMediaFromRequests{}); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, req)
}
