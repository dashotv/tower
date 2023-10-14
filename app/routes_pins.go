package app

import (
	"github.com/gin-gonic/gin"
)

func PinsCreate(c *gin.Context) {
	p := &Pin{}
	c.BindJSON(p)
	db.Pin.Save(p)
	id, err := schedulePinTask(p)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"job": id, "pin": p})
}

func PinsShow(c *gin.Context, id string) {
	// start polling job
	// return json
}

func PlexIndex(c *gin.Context) {
	// get pin
	// get auth url
	c.Redirect(302, "/blarg")
}
func PlexAuth(c *gin.Context) {
	// get plex token
}
