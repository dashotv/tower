package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dashotv/golem/web"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func MoviesIndex(c *gin.Context) {
	page, err := web.QueryDefaultInteger(c, "page", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := db.Movie.Count(bson.M{"_type": "Movie"})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := db.Movie.Query()
	results, err := q.
		Where("_type", "Movie").
		Limit(pagesize).
		Skip((page - 1) * pagesize).
		Desc("created_at").Run()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: do this with custom unmarshaling?
	for _, m := range results {
		m.Display = fmt.Sprintf("%s (%s)", m.Source, m.SourceId)
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

func MoviesCreate(c *gin.Context) {
	r := &CreateRequest{}
	c.BindJSON(r)
	if r.ID == "" || r.Source == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id and source are required"})
		return
	}
	m := &Movie{
		Type:     "Movie",
		SourceId: r.ID,
		Source:   r.Source,
		Title:    r.Title,
		Kind:     "movies",
	}
	d, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ReleaseDate = d

	err = db.Movie.Save(m)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: queue update job
	c.JSON(http.StatusOK, gin.H{"error": false, "movie": m})
}

func MoviesShow(c *gin.Context, id string) {
	m := &Movie{}
	err := db.Movie.Find(id, m)
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

func MoviesUpdate(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.MovieSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func MoviesSetting(c *gin.Context, id string) {
	data := &Setting{}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.MovieSetting(id, data.Setting, data.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errors": false, "data": data})
}

func MoviesDelete(c *gin.Context, id string) {
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func MoviesPaths(c *gin.Context, id string) {
	results, err := db.MoviePaths(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
