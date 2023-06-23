// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() {
	s.Router.GET("/", homeHandler)

	downloads := s.Router.Group("/downloads")
	downloads.POST("/", downloadsCreateHandler)
	downloads.DELETE("/:id", downloadsDeleteHandler)
	downloads.GET("/", downloadsIndexHandler)
	downloads.GET("/:id/medium", downloadsMediumHandler)
	downloads.GET("/recent", downloadsRecentHandler)
	downloads.PUT("/:id/select", downloadsSelectHandler)
	downloads.GET("/:id", downloadsShowHandler)
	downloads.PUT("/:id", downloadsUpdateHandler)

	episodes := s.Router.Group("/episodes")
	episodes.PUT("/:id", episodesUpdateHandler)

	feeds := s.Router.Group("/feeds")
	feeds.POST("/", feedsCreateHandler)
	feeds.DELETE("/:id", feedsDeleteHandler)
	feeds.GET("/", feedsIndexHandler)
	feeds.GET("/:id", feedsShowHandler)
	feeds.PUT("/:id", feedsUpdateHandler)

	movies := s.Router.Group("/movies")
	movies.POST("/", moviesCreateHandler)
	movies.DELETE("/:id", moviesDeleteHandler)
	movies.GET("/", moviesIndexHandler)
	movies.GET("/:id/paths", moviesPathsHandler)
	movies.GET("/:id", moviesShowHandler)
	movies.PUT("/:id", moviesUpdateHandler)

	releases := s.Router.Group("/releases")
	releases.POST("/", releasesCreateHandler)
	releases.DELETE("/:id", releasesDeleteHandler)
	releases.GET("/", releasesIndexHandler)
	releases.GET("/:id", releasesShowHandler)
	releases.PUT("/:id", releasesUpdateHandler)

	series := s.Router.Group("/series")
	series.POST("/", seriesCreateHandler)
	series.GET("/:id/currentseason", seriesCurrentSeasonHandler)
	series.DELETE("/:id", seriesDeleteHandler)
	series.GET("/", seriesIndexHandler)
	series.GET("/:id/paths", seriesPathsHandler)
	series.GET("/:id/seasons/:season", seriesSeasonEpisodesHandler)
	series.GET("/:id/seasons/all", seriesSeasonEpisodesAllHandler)
	series.GET("/:id/seasons", seriesSeasonsHandler)
	series.GET("/:id", seriesShowHandler)
	series.PUT("/:id", seriesUpdateHandler)
	series.GET("/:id/watches", seriesWatchesHandler)

	upcoming := s.Router.Group("/upcoming")
	upcoming.GET("/", upcomingIndexHandler)

}

func homeHandler(c *gin.Context) {
	Index(c)
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "home")
}

// /downloads
func downloadsCreateHandler(c *gin.Context) {

	DownloadsCreate(c)
}

func downloadsDeleteHandler(c *gin.Context) {
	id := c.Param("id")

	DownloadsDelete(c, id)
}

func downloadsIndexHandler(c *gin.Context) {

	DownloadsIndex(c)
}

func downloadsMediumHandler(c *gin.Context) {
	id := c.Param("id")

	DownloadsMedium(c, id)
}

func downloadsRecentHandler(c *gin.Context) {

	DownloadsRecent(c)
}

func downloadsSelectHandler(c *gin.Context) {
	id := c.Param("id")

	DownloadsSelect(c, id)
}

func downloadsShowHandler(c *gin.Context) {
	id := c.Param("id")

	DownloadsShow(c, id)
}

func downloadsUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	DownloadsUpdate(c, id)
}

// /episodes
func episodesUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	EpisodesUpdate(c, id)
}

// /feeds
func feedsCreateHandler(c *gin.Context) {

	FeedsCreate(c)
}

func feedsDeleteHandler(c *gin.Context) {
	id := c.Param("id")

	FeedsDelete(c, id)
}

func feedsIndexHandler(c *gin.Context) {

	FeedsIndex(c)
}

func feedsShowHandler(c *gin.Context) {
	id := c.Param("id")

	FeedsShow(c, id)
}

func feedsUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	FeedsUpdate(c, id)
}

// /movies
func moviesCreateHandler(c *gin.Context) {

	MoviesCreate(c)
}

func moviesDeleteHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesDelete(c, id)
}

func moviesIndexHandler(c *gin.Context) {

	MoviesIndex(c)
}

func moviesPathsHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesPaths(c, id)
}

func moviesShowHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesShow(c, id)
}

func moviesUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesUpdate(c, id)
}

// /releases
func releasesCreateHandler(c *gin.Context) {

	ReleasesCreate(c)
}

func releasesDeleteHandler(c *gin.Context) {
	id := c.Param("id")

	ReleasesDelete(c, id)
}

func releasesIndexHandler(c *gin.Context) {

	ReleasesIndex(c)
}

func releasesShowHandler(c *gin.Context) {
	id := c.Param("id")

	ReleasesShow(c, id)
}

func releasesUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	ReleasesUpdate(c, id)
}

// /series
func seriesCreateHandler(c *gin.Context) {

	SeriesCreate(c)
}

func seriesCurrentSeasonHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesCurrentSeason(c, id)
}

func seriesDeleteHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesDelete(c, id)
}

func seriesIndexHandler(c *gin.Context) {

	SeriesIndex(c)
}

func seriesPathsHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesPaths(c, id)
}

func seriesSeasonEpisodesHandler(c *gin.Context) {
	id := c.Param("id")
	season := c.Param("season")

	SeriesSeasonEpisodes(c, id, season)
}

func seriesSeasonEpisodesAllHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesSeasonEpisodesAll(c, id)
}

func seriesSeasonsHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesSeasons(c, id)
}

func seriesShowHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesShow(c, id)
}

func seriesUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesUpdate(c, id)
}

func seriesWatchesHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesWatches(c, id)
}

// /upcoming
func upcomingIndexHandler(c *gin.Context) {

	UpcomingIndex(c)
}
