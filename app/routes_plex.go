package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *Application) PlexIndex(c *gin.Context) {
	// get pin
	pin, err := app.Plex.CreatePin()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app.Log.Debugf("PlexIndex: saving pin %+v", pin)
	err = app.DB.Pin.Save(pin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authUrl := app.Plex.getAuthUrl(pin)
	// c.JSON(200, gin.H{"pin": pin, "authUrl": authUrl})
	c.Redirect(302, authUrl)
}

func (a *Application) PlexAuth(c *gin.Context) {
	id := c.Query("pin")
	pinId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	list, err := app.DB.Pin.Query().Where("pin", int64(pinId)).Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(list) != 1 {
		c.AbortWithStatusJSON(404, gin.H{"error": "pin not found"})
		return
	}

	pin := list[0]
	ok, err := app.Plex.CheckPin(pin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "something went wrong..."})
		return
	}

	// TODO: get user from token (myplex api), then scheduled job to handle watchlist?
	//app.Workers.Enqueue(runJob(&MinionJob{Name: "PlexAuth"}, func() error {
	// 	return nil
	// }))

	if err := app.Workers.Enqueue(&PlexPinToUsers{}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Authorization complete!")
}

func (a *Application) PlexUpdate(c *gin.Context) {
	if err := app.Workers.Enqueue(&PlexPinToUsers{}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Updating users...")
}

func (a *Application) PlexLibraries(c *gin.Context) {
	list, err := a.Plex.GetLibraries()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *Application) PlexSearch(c *gin.Context, query, section string) {
	list, err := a.Plex.Search(query, section)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (a *Application) PlexCollectionsIndex(c *gin.Context, section string) {
	list, err := a.Plex.ListCollections(section)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
func (a *Application) PlexCollectionsShow(c *gin.Context, section, ratingKey string) {
	list, err := a.Plex.GetCollection(ratingKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
