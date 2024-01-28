package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) IndexersIndex(c *gin.Context, page int, limit int) {
	indexers, count, err := a.DB.IndexerList(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"count": count, "results": indexers})
}

func (a *Application) IndexersCreate(c *gin.Context) {
	indexer := &Indexer{}
	if err := c.BindJSON(indexer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Indexer.Save(indexer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "result": indexer})
}

func (a *Application) IndexersShow(c *gin.Context, id string) {
	subject, err := a.DB.Indexer.Get(id, &Indexer{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "result": subject})
}

func (a *Application) IndexersUpdate(c *gin.Context, id string) {
	subject := &Indexer{}
	if err := c.BindJSON(subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Indexer.Save(subject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "result": subject})
}

func (a *Application) IndexersSettings(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	// subject, err := a.DB.Indexers.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func (a *Application) IndexersDelete(c *gin.Context, id string) {
	// asssuming this is a CRUD route, get the subject from the database
	subject, err := a.DB.Indexer.Get(id, &Indexer{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Indexer.Delete(subject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}
