package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/golem/web"
)

func (a *Application) MoviesIndex(c *gin.Context, page, limit int) {
	page, err := web.QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := app.DB.Movie.Count(bson.M{"_type": "Movie"})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := app.DB.Movie.Query()
	results, err := q.
		Limit(pagesize).
		Skip((page - 1) * pagesize).
		Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: do this with custom unmarshaling?
	for _, m := range results {
		for _, p := range m.Paths {
			if p.Type == "cover" {
				m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
			if p.Type == "background" {
				m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "results": results})
}

func (a *Application) MoviesCreate(c *gin.Context) {
	r := &CreateRequest{}
	c.BindJSON(r)
	if r.ID == "" || r.Source == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id and source are required"})
		return
	}

	m := &Movie{
		Type:         "Movie",
		SourceId:     r.ID,
		Source:       r.Source,
		Title:        r.Title,
		Description:  r.Description,
		Kind:         primitive.Symbol(r.Kind),
		SearchParams: &SearchParams{Type: "movies", Resolution: 1080, Verified: true},
	}

	d, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ReleaseDate = d

	err = app.DB.Movie.Save(m)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: m.ID.Hex()}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "movie": m})
}

func (a *Application) MoviesShow(c *gin.Context, id string) {
	m := &Movie{}
	err := app.DB.Movie.Find(id, m)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, p := range m.Paths {
		if p.Type == "cover" {
			m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
		if p.Type == "background" {
			m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
			continue
		}
	}

	c.JSON(http.StatusOK, m)
}

func (a *Application) MoviesUpdate(c *gin.Context, id string) {
	data := &Movie{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.MovieUpdate(id, data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) MoviesRefresh(c *gin.Context, id string) {
	if err := app.Workers.Enqueue(&TmdbUpdateMovie{ID: id}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) MoviesSettings(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.DB.MovieSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func (a *Application) MoviesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) MoviesPaths(c *gin.Context, id string) {
	results, err := app.DB.MoviePaths(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
