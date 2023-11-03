package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestsIndex(c *gin.Context) {
	list, err := db.Request.Query().Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func RequestsShow(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"message": "RequestsShow"})
}

func RequestsUpdate(c *gin.Context, id string) {
	req := &Request{}
	err := db.Request.Find(id, req)
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
	if err := db.Request.Update(req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}
