package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Application) MoviesIndex(c *gin.Context, page, limit int) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = pagesize
	}

	kind := QueryString(c, "kind")
	source := QueryString(c, "source")
	completed := QueryBool(c, "completed")
	downloaded := QueryBool(c, "downloaded")
	broken := QueryBool(c, "broken")

	q := app.DB.Movie.Query()
	if kind != "" {
		q = q.Where("kind", kind)
	}
	if source != "" {
		q = q.Where("source", source)
	}
	if completed {
		q = q.Where("completed", true)
	}
	if downloaded {
		q = q.Where("downloaded", true)
	}
	if broken {
		q = q.Where("broken", true)
	}

	count, err := q.Count()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
