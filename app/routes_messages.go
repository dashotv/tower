package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) MessagesIndex(c *gin.Context) {
	list, err := app.DB.Message.Query().Desc("created_at").Limit(250).Run()
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
