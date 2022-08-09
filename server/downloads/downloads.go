package downloads

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	results, err := app.DB.Download.Active()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, d := range results {
		m, err := app.DB.Medium.FindByID(d.MediumId)
		if err != nil {
			app.Log.Errorf("could not find medium: %s", d.MediumId)
			continue
		}
		app.Log.Infof("found %s: %s", m.ID, m.Title)
		results[i].Medium = *m
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
