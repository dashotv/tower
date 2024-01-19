package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RunicSourceSimple struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (a *Application) RunicIndex(c *gin.Context, page int, limit int) {
	out := make([]*RunicSourceSimple, 0)
	list := a.Runic.Sources()
	for _, n := range list {
		s, ok := a.Runic.Source(n)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("indexer does not exist")})
			return
		}
		out = append(out, &RunicSourceSimple{s.Name, s.Type, s.URL})
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"results": out,
	})
}

func (a *Application) RunicCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) RunicShow(c *gin.Context, id string) {
	s, ok := a.Runic.Source(id)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("indexer does not exist")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"source": s,
	})
}

func (a *Application) RunicRead(c *gin.Context, id string) {
	results, err := a.Runic.Read(id, []int{5000})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"source":  id,
		"results": results,
	})
}

func (a *Application) RunicSearch(c *gin.Context, id string, query string, searchType string) {
	results, err := a.Runic.Search(id, []int{5000}, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"source":  id,
		"results": results,
	})
}

func (a *Application) RunicUpdate(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Runic.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) RunicSettings(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Runic.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) RunicDelete(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Runic.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}
