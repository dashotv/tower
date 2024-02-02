package app

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

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

type Stuff struct {
	RatingKey    string `json:"ratingKey"`
	Key          string `json:"key"`
	GUID         string `json:"guid"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	LibraryID    int64  `json:"librarySectionID"`
	LibraryTitle string `json:"librarySectionTitle"`
	LibraryKey   string `json:"librarySectionKey"`
	Summary      string `json:"summary"`
	Thumb        string `json:"thumb"`
	Total        int    `json:"total"`
	Viewed       int    `json:"viewed"`
	LastViewedAt int64  `json:"lastViewedAt"`
	AddedAt      int64  `json:"addedAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}

type stuffSorter struct {
	list []*Stuff
	by   func(p1, p2 *Stuff) bool
}

// Len is part of sort.Interface.
func (s *stuffSorter) Len() int {
	return len(s.list)
}

// Swap is part of sort.Interface.
func (s *stuffSorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *stuffSorter) Less(i, j int) bool {
	return s.by(s.list[i], s.list[j])
}

func (a *Application) PlexStuff(c *gin.Context) {
	list := []*Stuff{}

	for _, i := range []string{"234979", "228425", "228426"} {
		children, err := a.Plex.GetCollectionChildren(i)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, child := range children {
			metadata, err := a.Plex.GetViewedByKey(child.RatingKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			stuff := &Stuff{
				RatingKey:    child.RatingKey,
				Key:          child.Key,
				GUID:         child.GUID,
				Type:         child.Type,
				Title:        child.Title,
				LibraryID:    child.LibraryID,
				LibraryTitle: child.LibraryTitle,
				LibraryKey:   child.LibraryKey,
				Summary:      child.Summary,
				Thumb:        child.Thumb,
				Total:        metadata.Leaves,
				Viewed:       metadata.Viewed,
				LastViewedAt: metadata.LastViewedAt,
				AddedAt:      child.AddedAt,
				UpdatedAt:    child.UpdatedAt,
			}
			list = append(list, stuff)
		}
	}
	sorter := &stuffSorter{
		list: list,
		by: func(p1, p2 *Stuff) bool {
			return p1.LastViewedAt > p2.LastViewedAt
		},
	}
	sort.Sort(sorter)
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
