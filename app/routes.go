// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package app

import (
	"net/http"

	"github.com/dashotv/golem/web"
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() {
	s.Default.GET("/", homeHandler)

	downloads := s.Router.Group("/downloads")
	downloads.POST("/", downloadsCreateHandler)
	downloads.DELETE("/:id", downloadsDeleteHandler)
	downloads.GET("/", downloadsIndexHandler)
	downloads.GET("/last", downloadsLastHandler)
	downloads.GET("/:id/medium", downloadsMediumHandler)
	downloads.GET("/recent", downloadsRecentHandler)
	downloads.PUT("/:id/select", downloadsSelectHandler)
	downloads.PATCH("/:id", downloadsSettingHandler)
	downloads.GET("/:id", downloadsShowHandler)
	downloads.PUT("/:id", downloadsUpdateHandler)

	episodes := s.Router.Group("/episodes")
	episodes.PATCH("/:id", episodesSettingHandler)
	episodes.PUT("/:id", episodesUpdateHandler)

	feeds := s.Router.Group("/feeds")
	feeds.POST("/", feedsCreateHandler)
	feeds.DELETE("/:id", feedsDeleteHandler)
	feeds.GET("/", feedsIndexHandler)
	feeds.PATCH("/:id", feedsSettingHandler)
	feeds.GET("/:id", feedsShowHandler)
	feeds.PUT("/:id", feedsUpdateHandler)

	jobs := s.Router.Group("/jobs")
	jobs.GET("/", jobsIndexHandler)

	messages := s.Router.Group("/messages")
	messages.GET("/", messagesIndexHandler)

	movies := s.Router.Group("/movies")
	movies.POST("/", moviesCreateHandler)
	movies.DELETE("/:id", moviesDeleteHandler)
	movies.GET("/", moviesIndexHandler)
	movies.GET("/:id/paths", moviesPathsHandler)
	movies.PUT("/:id/refresh", moviesRefreshHandler)
	movies.PATCH("/:id", moviesSettingHandler)
	movies.GET("/:id", moviesShowHandler)
	movies.PUT("/:id", moviesUpdateHandler)

	plex := s.Router.Group("/plex")
	plex.GET("/auth", plexAuthHandler)
	plex.GET("/", plexIndexHandler)
	plex.GET("/update", plexUpdateHandler)

	releases := s.Router.Group("/releases")
	releases.POST("/", releasesCreateHandler)
	releases.DELETE("/:id", releasesDeleteHandler)
	releases.GET("/", releasesIndexHandler)
	releases.GET("/popular/:interval", releasesPopularHandler)
	releases.PATCH("/:id", releasesSettingHandler)
	releases.GET("/:id", releasesShowHandler)
	releases.PUT("/:id", releasesUpdateHandler)

	requests := s.Router.Group("/requests")
	requests.GET("/", requestsIndexHandler)
	requests.GET("/:id", requestsShowHandler)
	requests.PUT("/:id", requestsUpdateHandler)

	series := s.Router.Group("/series")
	series.POST("/", seriesCreateHandler)
	series.GET("/:id/currentseason", seriesCurrentSeasonHandler)
	series.DELETE("/:id", seriesDeleteHandler)
	series.GET("/", seriesIndexHandler)
	series.GET("/:id/paths", seriesPathsHandler)
	series.PUT("/:id/refresh", seriesRefreshHandler)
	series.GET("/:id/seasons/:season", seriesSeasonEpisodesHandler)
	series.GET("/:id/seasons/all", seriesSeasonEpisodesAllHandler)
	series.GET("/:id/seasons", seriesSeasonsHandler)
	series.PATCH("/:id", seriesSettingHandler)
	series.GET("/:id", seriesShowHandler)
	series.PUT("/:id", seriesUpdateHandler)
	series.GET("/:id/watches", seriesWatchesHandler)

	upcoming := s.Router.Group("/upcoming")
	upcoming.GET("/", upcomingIndexHandler)

	users := s.Router.Group("/users")
	users.GET("/", usersIndexHandler)

	watches := s.Router.Group("/watches")
	watches.GET("/", watchesIndexHandler)

}

func homeHandler(c *gin.Context) {
	Index(c)
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": "tower",
		"routes": gin.H{
			"downloads": "/downloads",
			"episodes":  "/episodes",
			"feeds":     "/feeds",
			"jobs":      "/jobs",
			"messages":  "/messages",
			"movies":    "/movies",
			"plex":      "/plex",
			"releases":  "/releases",
			"requests":  "/requests",
			"series":    "/series",
			"upcoming":  "/upcoming",
			"users":     "/users",
			"watches":   "/watches",
		},
	})
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

func downloadsLastHandler(c *gin.Context) {

	DownloadsLast(c)
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

func downloadsSettingHandler(c *gin.Context) {
	id := c.Param("id")

	DownloadsSetting(c, id)
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
func episodesSettingHandler(c *gin.Context) {
	id := c.Param("id")

	EpisodesSetting(c, id)
}

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

func feedsSettingHandler(c *gin.Context) {
	id := c.Param("id")

	FeedsSetting(c, id)
}

func feedsShowHandler(c *gin.Context) {
	id := c.Param("id")

	FeedsShow(c, id)
}

func feedsUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	FeedsUpdate(c, id)
}

// /jobs
func jobsIndexHandler(c *gin.Context) {

	JobsIndex(c)
}

// /messages
func messagesIndexHandler(c *gin.Context) {

	MessagesIndex(c)
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

func moviesRefreshHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesRefresh(c, id)
}

func moviesSettingHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesSetting(c, id)
}

func moviesShowHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesShow(c, id)
}

func moviesUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	MoviesUpdate(c, id)
}

// /plex
func plexAuthHandler(c *gin.Context) {

	PlexAuth(c)
}

func plexIndexHandler(c *gin.Context) {

	PlexIndex(c)
}

func plexUpdateHandler(c *gin.Context) {

	PlexUpdate(c)
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

func releasesPopularHandler(c *gin.Context) {
	interval := c.Param("interval")

	ReleasesPopular(c, interval)
}

func releasesSettingHandler(c *gin.Context) {
	id := c.Param("id")

	ReleasesSetting(c, id)
}

func releasesShowHandler(c *gin.Context) {
	id := c.Param("id")

	ReleasesShow(c, id)
}

func releasesUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	ReleasesUpdate(c, id)
}

// /requests
func requestsIndexHandler(c *gin.Context) {

	RequestsIndex(c)
}

func requestsShowHandler(c *gin.Context) {
	id := c.Param("id")

	RequestsShow(c, id)
}

func requestsUpdateHandler(c *gin.Context) {
	id := c.Param("id")

	RequestsUpdate(c, id)
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

func seriesRefreshHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesRefresh(c, id)
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

func seriesSettingHandler(c *gin.Context) {
	id := c.Param("id")

	SeriesSetting(c, id)
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

// /users
func usersIndexHandler(c *gin.Context) {

	UsersIndex(c)
}

// /watches
func watchesIndexHandler(c *gin.Context) {
	medium_id := web.QueryString(c, "medium_id")
	username := web.QueryString(c, "username")

	WatchesIndex(c, medium_id, username)
}
