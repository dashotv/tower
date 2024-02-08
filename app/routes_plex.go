package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/dashotv/tower/internal/plex"
)

func plexPinToPin(pin *plex.Pin) *Pin {
	return &Pin{
		Pin:        pin.ID,
		Code:       pin.Code,
		Product:    pin.Product,
		Identifier: pin.Identifier,
		Token:      pin.Token,
	}
}

func pinToPlexPin(pin *Pin) *plex.Pin {
	return &plex.Pin{
		ID:         pin.Pin,
		Code:       pin.Code,
		Product:    pin.Product,
		Identifier: pin.Identifier,
		Token:      pin.Token,
	}
}

func (a *Application) PlexIndex(c *gin.Context) {
	// get pin
	plexPin, err := app.Plex.CreatePin()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pin := plexPinToPin(plexPin)

	app.Log.Debugf("PlexIndex: saving pin %+v", pin)
	err = app.DB.Pin.Save(pin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authUrl := app.Plex.GetAuthUrl(app.Config.Plex, plexPin)
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

	plexPin := pinToPlexPin(list[0])
	ok, err := app.Plex.CheckPin(plexPin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "something went wrong..."})
		return
	}

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

func (a *Application) PlexMetadata(c *gin.Context, key string) {
	a.Log.Debugf("PlexMetadata: key=%s", key)
	resp, err := a.Plex.GetMetadataByKey(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *Application) PlexClients(c *gin.Context) {
	list, err := a.Plex.GetClients()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *Application) PlexDevices(c *gin.Context) {
	list, err := a.Plex.GetDevices()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *Application) PlexResources(c *gin.Context) {
	provides := QueryString(c, "provides")
	if provides == "" {
		provides = "player"
	}

	list, err := a.Plex.GetResources()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// filter by provides
	filtered := lo.Filter(list, func(r *plex.Resource, i int) bool {
		provided := strings.Split(r.Provides, ",")
		return lo.Contains(provided, provides)
	})

	c.JSON(http.StatusOK, filtered)
}
func (a *Application) PlexPlay(c *gin.Context, ratingKey, player string) {
	err := a.Plex.Play(player, ratingKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Playing...")
}
func (a *Application) PlexSessions(c *gin.Context) {
	list, err := a.Plex.GetSessions()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
