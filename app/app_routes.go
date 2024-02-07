// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers = append(initializers, setupRoutes)
	healthchecks["routes"] = checkRoutes
}

func checkRoutes(app *Application) error {
	// TODO: check routes
	return nil
}

func setupRoutes(app *Application) error {
	if app.Config.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := app.Log.Named("routes").Desugar()

	app.Engine = gin.New()
	app.Engine.Use(
		ginzap.Ginzap(logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger, true),
	)
	// unauthenticated routes
	app.Default = app.Engine.Group("/")
	// authenticated routes (if enabled, otherwise same as default)
	app.Router = app.Engine.Group("/")

	// if app.Config.Auth {
	// 	clerkSecret := app.Config.ClerkSecretKey
	// 	if clerkSecret == "" {
	// 		app.Log.Fatal("CLERK_SECRET_KEY is not set")
	// 	}
	//
	// 	clerkClient, err := clerk.NewClient(clerkSecret)
	// 	if err != nil {
	// 		app.Log.Fatalf("clerk: %s", err)
	// 	}
	//
	// 	app.Router.Use(requireSession(clerkClient))
	// }

	return nil
}

// Enable Auth and uncomment to use Clerk to manage auth
// also add this import: "github.com/clerkinc/clerk-sdk-go/clerk"
//
// requireSession wraps the clerk.RequireSession middleware
// func requireSession(client clerk.Client) gin.HandlerFunc {
// 	requireActiveSession := clerk.RequireSessionV2(client)
// 	return func(gctx *gin.Context) {
// 		var skip = true
// 		var handler http.HandlerFunc = func(http.ResponseWriter, *http.Request) {
// 			skip = false
// 		}
// 		requireActiveSession(handler).ServeHTTP(gctx.Writer, gctx.Request)
// 		switch {
// 		case skip:
// 			gctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "session required"})
// 		default:
// 			gctx.Next()
// 		}
// 	}
// }

func (a *Application) Routes() {
	a.Default.GET("/", a.indexHandler)
	a.Default.GET("/health", a.healthHandler)

	collections := a.Router.Group("/collections")
	collections.GET("/", a.CollectionsIndexHandler)
	collections.POST("/", a.CollectionsCreateHandler)
	collections.GET("/:id", a.CollectionsShowHandler)
	collections.PUT("/:id", a.CollectionsUpdateHandler)
	collections.PATCH("/:id", a.CollectionsSettingsHandler)
	collections.DELETE("/:id", a.CollectionsDeleteHandler)

	downloads := a.Router.Group("/downloads")
	downloads.GET("/", a.DownloadsIndexHandler)
	downloads.POST("/", a.DownloadsCreateHandler)
	downloads.GET("/:id", a.DownloadsShowHandler)
	downloads.PUT("/:id", a.DownloadsUpdateHandler)
	downloads.PATCH("/:id", a.DownloadsSettingsHandler)
	downloads.DELETE("/:id", a.DownloadsDeleteHandler)
	downloads.GET("/last", a.DownloadsLastHandler)
	downloads.GET("/:id/medium", a.DownloadsMediumHandler)
	downloads.GET("/recent", a.DownloadsRecentHandler)
	downloads.PUT("/:id/select", a.DownloadsSelectHandler)

	episodes := a.Router.Group("/episodes")
	episodes.PATCH("/:id", a.EpisodesSettingHandler)
	episodes.PUT("/:id", a.EpisodesUpdateHandler)
	episodes.POST("/settings", a.EpisodesSettingsHandler)

	feeds := a.Router.Group("/feeds")
	feeds.GET("/", a.FeedsIndexHandler)
	feeds.POST("/", a.FeedsCreateHandler)
	feeds.GET("/:id", a.FeedsShowHandler)
	feeds.PUT("/:id", a.FeedsUpdateHandler)
	feeds.PATCH("/:id", a.FeedsSettingsHandler)
	feeds.DELETE("/:id", a.FeedsDeleteHandler)

	hooks := a.Router.Group("/hooks")
	hooks.GET("/plex", a.HooksPlexHandler)

	jobs := a.Router.Group("/jobs")
	jobs.GET("/", a.JobsIndexHandler)
	jobs.POST("/", a.JobsCreateHandler)
	jobs.DELETE("/:id", a.JobsDeleteHandler)

	messages := a.Router.Group("/messages")
	messages.GET("/", a.MessagesIndexHandler)
	messages.POST("/", a.MessagesCreateHandler)

	movies := a.Router.Group("/movies")
	movies.GET("/", a.MoviesIndexHandler)
	movies.POST("/", a.MoviesCreateHandler)
	movies.GET("/:id", a.MoviesShowHandler)
	movies.PUT("/:id", a.MoviesUpdateHandler)
	movies.PATCH("/:id", a.MoviesSettingsHandler)
	movies.DELETE("/:id", a.MoviesDeleteHandler)
	movies.PUT("/:id/refresh", a.MoviesRefreshHandler)
	movies.GET("/:id/paths", a.MoviesPathsHandler)

	plex := a.Router.Group("/plex")
	plex.GET("/auth", a.PlexAuthHandler)
	plex.GET("/", a.PlexIndexHandler)
	plex.GET("/update", a.PlexUpdateHandler)
	plex.GET("/search", a.PlexSearchHandler)
	plex.GET("/libraries", a.PlexLibrariesHandler)
	plex.GET("/libraries/:section/collections", a.PlexCollectionsIndexHandler)
	plex.GET("/libraries/:section/collections/:ratingKey", a.PlexCollectionsShowHandler)
	plex.GET("/stuff", a.PlexStuffHandler)
	plex.GET("/metadata/:key", a.PlexMetadataHandler)
	plex.GET("/clients", a.PlexClientsHandler)
	plex.GET("/devices", a.PlexDevicesHandler)
	plex.GET("/resources", a.PlexResourcesHandler)
	plex.GET("/play", a.PlexPlayHandler)
	plex.GET("/sessions", a.PlexSessionsHandler)

	releases := a.Router.Group("/releases")
	releases.GET("/", a.ReleasesIndexHandler)
	releases.POST("/", a.ReleasesCreateHandler)
	releases.GET("/:id", a.ReleasesShowHandler)
	releases.PUT("/:id", a.ReleasesUpdateHandler)
	releases.PATCH("/:id", a.ReleasesSettingsHandler)
	releases.DELETE("/:id", a.ReleasesDeleteHandler)
	releases.GET("/popular/:interval", a.ReleasesPopularHandler)

	requests := a.Router.Group("/requests")
	requests.GET("/", a.RequestsIndexHandler)
	requests.POST("/", a.RequestsCreateHandler)
	requests.GET("/:id", a.RequestsShowHandler)
	requests.PUT("/:id", a.RequestsUpdateHandler)
	requests.PATCH("/:id", a.RequestsSettingsHandler)
	requests.DELETE("/:id", a.RequestsDeleteHandler)

	series := a.Router.Group("/series")
	series.GET("/", a.SeriesIndexHandler)
	series.POST("/", a.SeriesCreateHandler)
	series.GET("/:id", a.SeriesShowHandler)
	series.PUT("/:id", a.SeriesUpdateHandler)
	series.PATCH("/:id", a.SeriesSettingsHandler)
	series.DELETE("/:id", a.SeriesDeleteHandler)
	series.GET("/:id/currentseason", a.SeriesCurrentSeasonHandler)
	series.GET("/:id/paths", a.SeriesPathsHandler)
	series.PUT("/:id/refresh", a.SeriesRefreshHandler)
	series.GET("/:id/seasons/all", a.SeriesSeasonEpisodesAllHandler)
	series.GET("/:id/seasons/:season", a.SeriesSeasonEpisodesHandler)
	series.GET("/:id/watches", a.SeriesWatchesHandler)
	series.GET("/:id/covers", a.SeriesCoversHandler)
	series.GET("/:id/backgrounds", a.SeriesBackgroundsHandler)

	upcoming := a.Router.Group("/upcoming")
	upcoming.GET("/", a.UpcomingIndexHandler)

	users := a.Router.Group("/users")
	users.GET("/", a.UsersIndexHandler)

	watches := a.Router.Group("/watches")
	watches.GET("/", a.WatchesIndexHandler)

}

func (a *Application) indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": "tower",
		"routes": gin.H{
			"collections": "/collections",
			"downloads":   "/downloads",
			"episodes":    "/episodes",
			"feeds":       "/feeds",
			"hooks":       "/hooks",
			"jobs":        "/jobs",
			"messages":    "/messages",
			"movies":      "/movies",
			"plex":        "/plex",
			"releases":    "/releases",
			"requests":    "/requests",
			"series":      "/series",
			"upcoming":    "/upcoming",
			"users":       "/users",
			"watches":     "/watches",
		},
	})
}

func (a *Application) healthHandler(c *gin.Context) {
	health, err := a.Health()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": "tower", "health": health})
}

// Collections (/collections)
func (a *Application) CollectionsIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.CollectionsIndex(c, page, limit)
}
func (a *Application) CollectionsCreateHandler(c *gin.Context) {
	a.CollectionsCreate(c)
}
func (a *Application) CollectionsShowHandler(c *gin.Context) {
	id := c.Param("id")
	a.CollectionsShow(c, id)
}
func (a *Application) CollectionsUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.CollectionsUpdate(c, id)
}
func (a *Application) CollectionsSettingsHandler(c *gin.Context) {
	id := c.Param("id")
	a.CollectionsSettings(c, id)
}
func (a *Application) CollectionsDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	a.CollectionsDelete(c, id)
}

// Downloads (/downloads)
func (a *Application) DownloadsIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.DownloadsIndex(c, page, limit)
}
func (a *Application) DownloadsCreateHandler(c *gin.Context) {
	a.DownloadsCreate(c)
}
func (a *Application) DownloadsShowHandler(c *gin.Context) {
	id := c.Param("id")
	a.DownloadsShow(c, id)
}
func (a *Application) DownloadsUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.DownloadsUpdate(c, id)
}
func (a *Application) DownloadsSettingsHandler(c *gin.Context) {
	id := c.Param("id")
	a.DownloadsSettings(c, id)
}
func (a *Application) DownloadsDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	a.DownloadsDelete(c, id)
}
func (a *Application) DownloadsLastHandler(c *gin.Context) {
	a.DownloadsLast(c)
}
func (a *Application) DownloadsMediumHandler(c *gin.Context) {
	id := c.Param("id")
	a.DownloadsMedium(c, id)
}
func (a *Application) DownloadsRecentHandler(c *gin.Context) {
	a.DownloadsRecent(c)
}
func (a *Application) DownloadsSelectHandler(c *gin.Context) {
	id := c.Param("id")
	a.DownloadsSelect(c, id)
}

// Episodes (/episodes)
func (a *Application) EpisodesSettingHandler(c *gin.Context) {
	id := c.Param("id")
	a.EpisodesSetting(c, id)
}
func (a *Application) EpisodesUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.EpisodesUpdate(c, id)
}
func (a *Application) EpisodesSettingsHandler(c *gin.Context) {
	a.EpisodesSettings(c)
}

// Feeds (/feeds)
func (a *Application) FeedsIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.FeedsIndex(c, page, limit)
}
func (a *Application) FeedsCreateHandler(c *gin.Context) {
	a.FeedsCreate(c)
}
func (a *Application) FeedsShowHandler(c *gin.Context) {
	id := c.Param("id")
	a.FeedsShow(c, id)
}
func (a *Application) FeedsUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.FeedsUpdate(c, id)
}
func (a *Application) FeedsSettingsHandler(c *gin.Context) {
	id := c.Param("id")
	a.FeedsSettings(c, id)
}
func (a *Application) FeedsDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	a.FeedsDelete(c, id)
}

// Hooks (/hooks)
func (a *Application) HooksPlexHandler(c *gin.Context) {
	a.HooksPlex(c)
}

// Jobs (/jobs)
func (a *Application) JobsIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.JobsIndex(c, page, limit)
}
func (a *Application) JobsCreateHandler(c *gin.Context) {
	job := QueryString(c, "job")
	a.JobsCreate(c, job)
}
func (a *Application) JobsDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	hard := QueryBool(c, "hard")
	a.JobsDelete(c, id, hard)
}

// Messages (/messages)
func (a *Application) MessagesIndexHandler(c *gin.Context) {
	a.MessagesIndex(c)
}
func (a *Application) MessagesCreateHandler(c *gin.Context) {
	a.MessagesCreate(c)
}

// Movies (/movies)
func (a *Application) MoviesIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.MoviesIndex(c, page, limit)
}
func (a *Application) MoviesCreateHandler(c *gin.Context) {
	a.MoviesCreate(c)
}
func (a *Application) MoviesShowHandler(c *gin.Context) {
	id := c.Param("id")
	a.MoviesShow(c, id)
}
func (a *Application) MoviesUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.MoviesUpdate(c, id)
}
func (a *Application) MoviesSettingsHandler(c *gin.Context) {
	id := c.Param("id")
	a.MoviesSettings(c, id)
}
func (a *Application) MoviesDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	a.MoviesDelete(c, id)
}
func (a *Application) MoviesRefreshHandler(c *gin.Context) {
	id := c.Param("id")
	a.MoviesRefresh(c, id)
}
func (a *Application) MoviesPathsHandler(c *gin.Context) {
	id := c.Param("id")
	a.MoviesPaths(c, id)
}

// Plex (/plex)
func (a *Application) PlexAuthHandler(c *gin.Context) {
	a.PlexAuth(c)
}
func (a *Application) PlexIndexHandler(c *gin.Context) {
	a.PlexIndex(c)
}
func (a *Application) PlexUpdateHandler(c *gin.Context) {
	a.PlexUpdate(c)
}
func (a *Application) PlexSearchHandler(c *gin.Context) {
	query := QueryString(c, "query")
	section := QueryString(c, "section")
	a.PlexSearch(c, query, section)
}
func (a *Application) PlexLibrariesHandler(c *gin.Context) {
	a.PlexLibraries(c)
}
func (a *Application) PlexCollectionsIndexHandler(c *gin.Context) {
	section := c.Param("section")
	a.PlexCollectionsIndex(c, section)
}
func (a *Application) PlexCollectionsShowHandler(c *gin.Context) {
	section := c.Param("section")
	ratingKey := c.Param("ratingKey")
	a.PlexCollectionsShow(c, section, ratingKey)
}
func (a *Application) PlexStuffHandler(c *gin.Context) {
	a.PlexStuff(c)
}
func (a *Application) PlexMetadataHandler(c *gin.Context) {
	key := c.Param("key")
	a.PlexMetadata(c, key)
}
func (a *Application) PlexClientsHandler(c *gin.Context) {
	a.PlexClients(c)
}
func (a *Application) PlexDevicesHandler(c *gin.Context) {
	a.PlexDevices(c)
}
func (a *Application) PlexResourcesHandler(c *gin.Context) {
	a.PlexResources(c)
}
func (a *Application) PlexPlayHandler(c *gin.Context) {
	ratingKey := QueryString(c, "ratingKey")
	player := QueryString(c, "player")
	a.PlexPlay(c, ratingKey, player)
}
func (a *Application) PlexSessionsHandler(c *gin.Context) {
	a.PlexSessions(c)
}

// Releases (/releases)
func (a *Application) ReleasesIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.ReleasesIndex(c, page, limit)
}
func (a *Application) ReleasesCreateHandler(c *gin.Context) {
	a.ReleasesCreate(c)
}
func (a *Application) ReleasesShowHandler(c *gin.Context) {
	id := c.Param("id")
	a.ReleasesShow(c, id)
}
func (a *Application) ReleasesUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.ReleasesUpdate(c, id)
}
func (a *Application) ReleasesSettingsHandler(c *gin.Context) {
	id := c.Param("id")
	a.ReleasesSettings(c, id)
}
func (a *Application) ReleasesDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	a.ReleasesDelete(c, id)
}
func (a *Application) ReleasesPopularHandler(c *gin.Context) {
	interval := c.Param("interval")
	a.ReleasesPopular(c, interval)
}

// Requests (/requests)
func (a *Application) RequestsIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.RequestsIndex(c, page, limit)
}
func (a *Application) RequestsCreateHandler(c *gin.Context) {
	a.RequestsCreate(c)
}
func (a *Application) RequestsShowHandler(c *gin.Context) {
	id := c.Param("id")
	a.RequestsShow(c, id)
}
func (a *Application) RequestsUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.RequestsUpdate(c, id)
}
func (a *Application) RequestsSettingsHandler(c *gin.Context) {
	id := c.Param("id")
	a.RequestsSettings(c, id)
}
func (a *Application) RequestsDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	a.RequestsDelete(c, id)
}

// Series (/series)
func (a *Application) SeriesIndexHandler(c *gin.Context) {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	a.SeriesIndex(c, page, limit)
}
func (a *Application) SeriesCreateHandler(c *gin.Context) {
	a.SeriesCreate(c)
}
func (a *Application) SeriesShowHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesShow(c, id)
}
func (a *Application) SeriesUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesUpdate(c, id)
}
func (a *Application) SeriesSettingsHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesSettings(c, id)
}
func (a *Application) SeriesDeleteHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesDelete(c, id)
}
func (a *Application) SeriesCurrentSeasonHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesCurrentSeason(c, id)
}
func (a *Application) SeriesPathsHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesPaths(c, id)
}
func (a *Application) SeriesRefreshHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesRefresh(c, id)
}
func (a *Application) SeriesSeasonEpisodesAllHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesSeasonEpisodesAll(c, id)
}
func (a *Application) SeriesSeasonEpisodesHandler(c *gin.Context) {
	id := c.Param("id")
	season := c.Param("season")
	a.SeriesSeasonEpisodes(c, id, season)
}
func (a *Application) SeriesWatchesHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesWatches(c, id)
}
func (a *Application) SeriesCoversHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesCovers(c, id)
}
func (a *Application) SeriesBackgroundsHandler(c *gin.Context) {
	id := c.Param("id")
	a.SeriesBackgrounds(c, id)
}

// Upcoming (/upcoming)
func (a *Application) UpcomingIndexHandler(c *gin.Context) {
	a.UpcomingIndex(c)
}

// Users (/users)
func (a *Application) UsersIndexHandler(c *gin.Context) {
	a.UsersIndex(c)
}

// Watches (/watches)
func (a *Application) WatchesIndexHandler(c *gin.Context) {
	medium_id := QueryString(c, "medium_id")
	username := QueryString(c, "username")
	a.WatchesIndex(c, medium_id, username)
}
