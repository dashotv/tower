package releases

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	results, err := app.DB.Release.Query().Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Show(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Update(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Delete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}
