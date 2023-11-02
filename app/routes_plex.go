package app

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func PlexIndex(c *gin.Context) {
	// get pin
	pin, err := plex.CreatePin()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	server.Log.Debugf("PlexIndex: saving pin %+v", pin)
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

	// TODO: get user from token (myplex api), then scheduled job to handle watchlist?
	// workers.Enqueue(runJob(&MinionJob{Name: "PlexAuth"}, func() error {
	// 	return nil
	// }))

	workers.Enqueue("PlexUserUpdates")
	c.String(200, "Authorization complete!")
}

func PlexUpdate(c *gin.Context) {
	workers.Enqueue("PlexUserUpdates")
	c.String(200, "Updating users...")
}
