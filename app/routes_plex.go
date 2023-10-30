package app

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PlexIndex(c *gin.Context) {
	// get pin
	pin, err := plex.CreatePin()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	err = db.Pin.Save(pin)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	authUrl := plex.getAuthUrl(pin)
	// c.JSON(200, gin.H{"pin": pin, "authUrl": authUrl})
	c.Redirect(302, authUrl)
}

func PlexAuth(c *gin.Context) {
	id := c.Query("pin")
	pinId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	list, err := db.Pin.Query().Where("pin", int64(pinId)).Run()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	server.Log.Infof("list: %d %+v", pinId, list)
	if len(list) != 1 {
		c.AbortWithStatusJSON(404, gin.H{"error": "pin not found"})
		return
	}

	pin := list[0]
	ok, err := plex.CheckPin(pin)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(401, gin.H{"error": "something went wrong..."})
		return
	}

	minion.Add("plex get user", func(id int, log *zap.SugaredLogger) error {
		// TODO: need working plex client
		return nil
	})

	// TODO: get user from token (call myplex api), maybe background this?
	c.String(200, "Authorization complete!")
}
