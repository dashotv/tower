package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) MessagesIndex(c *gin.Context) {
	page, err := QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	limit, err := QueryDefaultInteger(c, "limit", 250)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	list, err := app.DB.Message.Query().Desc("created_at").Skip((page - 1) * limit).Limit(limit).Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (a *Application) MessagesCreate(c *gin.Context) {
	m := &Message{}
	if err := c.BindJSON(m); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := app.DB.Message.Save(m); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := a.Events.Send("tower.logs", &EventLogs{Event: "new", Id: m.ID.Hex(), Log: m}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}
