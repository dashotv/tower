// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.infratographer.com/x/echox/echozap"
)

func init() {
	initializers = append(initializers, setupRoutes)
	healthchecks["routes"] = checkRoutes
	starters = append(starters, startRoutes)
}

func checkRoutes(app *Application) error {
	// TODO: check routes
	return nil
}

func startRoutes(ctx context.Context, app *Application) error {
	go func() {
		app.Routes()
		app.Log.Info("starting routes...")
		if err := app.Engine.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
			app.Log.Errorf("routes: %s", err)
		}
	}()
	return nil
}

func setupRoutes(app *Application) error {
	logger := app.Log.Named("routes").Desugar()
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(echozap.Middleware(logger))

	app.Engine = e
	// unauthenticated routes
	app.Default = app.Engine.Group("")
	// authenticated routes (if enabled, otherwise same as default)
	app.Router = app.Engine.Group("")

	// TODO: fix auth
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

	combinations := a.Router.Group("/combinations")
	combinations.GET("/", a.CombinationsIndexHandler)
	combinations.POST("/", a.CombinationsCreateHandler)
	combinations.GET("/:id", a.CombinationsShowHandler)
	combinations.PUT("/:id", a.CombinationsUpdateHandler)
	combinations.PATCH("/:id", a.CombinationsSettingsHandler)
	combinations.DELETE("/:id", a.CombinationsDeleteHandler)

	config := a.Router.Group("/config")
	config.GET("/", a.ConfigIndexHandler)
	config.POST("/", a.ConfigCreateHandler)
	config.GET("/:id", a.ConfigShowHandler)
	config.PUT("/:id", a.ConfigUpdateHandler)
	config.PATCH("/:id", a.ConfigSettingsHandler)
	config.DELETE("/:id", a.ConfigDeleteHandler)

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
	downloads.GET("/:id/torrent", a.DownloadsTorrentHandler)

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
	plex.GET("/metadata/:key", a.PlexMetadataHandler)
	plex.GET("/clients", a.PlexClientsHandler)
	plex.GET("/devices", a.PlexDevicesHandler)
	plex.GET("/resources", a.PlexResourcesHandler)
	plex.GET("/play", a.PlexPlayHandler)
	plex.GET("/sessions", a.PlexSessionsHandler)
	plex.GET("/stop", a.PlexStopHandler)

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
	series.POST("/:id/jobs", a.SeriesJobsHandler)

	upcoming := a.Router.Group("/upcoming")
	upcoming.GET("/", a.UpcomingIndexHandler)

	users := a.Router.Group("/users")
	users.GET("/", a.UsersIndexHandler)

	watches := a.Router.Group("/watches")
	watches.GET("/", a.WatchesIndexHandler)

}

func (a *Application) indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, H{
		"name": "tower",
		"routes": H{
			"collections":  "/collections",
			"combinations": "/combinations",
			"config":       "/config",
			"downloads":    "/downloads",
			"episodes":     "/episodes",
			"feeds":        "/feeds",
			"hooks":        "/hooks",
			"jobs":         "/jobs",
			"messages":     "/messages",
			"movies":       "/movies",
			"plex":         "/plex",
			"releases":     "/releases",
			"requests":     "/requests",
			"series":       "/series",
			"upcoming":     "/upcoming",
			"users":        "/users",
			"watches":      "/watches",
		},
	})
}

func (a *Application) healthHandler(c echo.Context) error {
	health, err := a.Health()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, H{"name": "tower", "health": health})
}

// Collections (/collections)
func (a *Application) CollectionsIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.CollectionsIndex(c, page, limit)
}
func (a *Application) CollectionsCreateHandler(c echo.Context) error {
	return a.CollectionsCreate(c)
}
func (a *Application) CollectionsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CollectionsShow(c, id)
}
func (a *Application) CollectionsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CollectionsUpdate(c, id)
}
func (a *Application) CollectionsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CollectionsSettings(c, id)
}
func (a *Application) CollectionsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CollectionsDelete(c, id)
}

// Combinations (/combinations)
func (a *Application) CombinationsIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.CombinationsIndex(c, page, limit)
}
func (a *Application) CombinationsCreateHandler(c echo.Context) error {
	return a.CombinationsCreate(c)
}
func (a *Application) CombinationsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CombinationsShow(c, id)
}
func (a *Application) CombinationsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CombinationsUpdate(c, id)
}
func (a *Application) CombinationsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CombinationsSettings(c, id)
}
func (a *Application) CombinationsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CombinationsDelete(c, id)
}

// Config (/config)
func (a *Application) ConfigIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.ConfigIndex(c, page, limit)
}
func (a *Application) ConfigCreateHandler(c echo.Context) error {
	return a.ConfigCreate(c)
}
func (a *Application) ConfigShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ConfigShow(c, id)
}
func (a *Application) ConfigUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ConfigUpdate(c, id)
}
func (a *Application) ConfigSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ConfigSettings(c, id)
}
func (a *Application) ConfigDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ConfigDelete(c, id)
}

// Downloads (/downloads)
func (a *Application) DownloadsIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.DownloadsIndex(c, page, limit)
}
func (a *Application) DownloadsCreateHandler(c echo.Context) error {
	return a.DownloadsCreate(c)
}
func (a *Application) DownloadsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsShow(c, id)
}
func (a *Application) DownloadsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsUpdate(c, id)
}
func (a *Application) DownloadsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsSettings(c, id)
}
func (a *Application) DownloadsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsDelete(c, id)
}
func (a *Application) DownloadsLastHandler(c echo.Context) error {
	return a.DownloadsLast(c)
}
func (a *Application) DownloadsMediumHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsMedium(c, id)
}
func (a *Application) DownloadsRecentHandler(c echo.Context) error {
	return a.DownloadsRecent(c)
}
func (a *Application) DownloadsSelectHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsSelect(c, id)
}
func (a *Application) DownloadsTorrentHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsTorrent(c, id)
}

// Episodes (/episodes)
func (a *Application) EpisodesSettingHandler(c echo.Context) error {
	id := c.Param("id")
	return a.EpisodesSetting(c, id)
}
func (a *Application) EpisodesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.EpisodesUpdate(c, id)
}
func (a *Application) EpisodesSettingsHandler(c echo.Context) error {
	return a.EpisodesSettings(c)
}

// Feeds (/feeds)
func (a *Application) FeedsIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.FeedsIndex(c, page, limit)
}
func (a *Application) FeedsCreateHandler(c echo.Context) error {
	return a.FeedsCreate(c)
}
func (a *Application) FeedsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.FeedsShow(c, id)
}
func (a *Application) FeedsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.FeedsUpdate(c, id)
}
func (a *Application) FeedsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.FeedsSettings(c, id)
}
func (a *Application) FeedsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.FeedsDelete(c, id)
}

// Hooks (/hooks)
func (a *Application) HooksPlexHandler(c echo.Context) error {
	return a.HooksPlex(c)
}

// Jobs (/jobs)
func (a *Application) JobsIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.JobsIndex(c, page, limit)
}
func (a *Application) JobsCreateHandler(c echo.Context) error {
	job := QueryString(c, "job")
	return a.JobsCreate(c, job)
}
func (a *Application) JobsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	hard := QueryBool(c, "hard")
	return a.JobsDelete(c, id, hard)
}

// Messages (/messages)
func (a *Application) MessagesIndexHandler(c echo.Context) error {
	return a.MessagesIndex(c)
}
func (a *Application) MessagesCreateHandler(c echo.Context) error {
	return a.MessagesCreate(c)
}

// Movies (/movies)
func (a *Application) MoviesIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.MoviesIndex(c, page, limit)
}
func (a *Application) MoviesCreateHandler(c echo.Context) error {
	return a.MoviesCreate(c)
}
func (a *Application) MoviesShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.MoviesShow(c, id)
}
func (a *Application) MoviesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.MoviesUpdate(c, id)
}
func (a *Application) MoviesSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.MoviesSettings(c, id)
}
func (a *Application) MoviesDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.MoviesDelete(c, id)
}
func (a *Application) MoviesRefreshHandler(c echo.Context) error {
	id := c.Param("id")
	return a.MoviesRefresh(c, id)
}
func (a *Application) MoviesPathsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.MoviesPaths(c, id)
}

// Plex (/plex)
func (a *Application) PlexAuthHandler(c echo.Context) error {
	return a.PlexAuth(c)
}
func (a *Application) PlexIndexHandler(c echo.Context) error {
	return a.PlexIndex(c)
}
func (a *Application) PlexUpdateHandler(c echo.Context) error {
	return a.PlexUpdate(c)
}
func (a *Application) PlexSearchHandler(c echo.Context) error {
	query := QueryString(c, "query")
	section := QueryString(c, "section")
	return a.PlexSearch(c, query, section)
}
func (a *Application) PlexLibrariesHandler(c echo.Context) error {
	return a.PlexLibraries(c)
}
func (a *Application) PlexCollectionsIndexHandler(c echo.Context) error {
	section := c.Param("section")
	return a.PlexCollectionsIndex(c, section)
}
func (a *Application) PlexCollectionsShowHandler(c echo.Context) error {
	section := c.Param("section")
	ratingKey := c.Param("ratingKey")
	return a.PlexCollectionsShow(c, section, ratingKey)
}
func (a *Application) PlexMetadataHandler(c echo.Context) error {
	key := c.Param("key")
	return a.PlexMetadata(c, key)
}
func (a *Application) PlexClientsHandler(c echo.Context) error {
	return a.PlexClients(c)
}
func (a *Application) PlexDevicesHandler(c echo.Context) error {
	return a.PlexDevices(c)
}
func (a *Application) PlexResourcesHandler(c echo.Context) error {
	return a.PlexResources(c)
}
func (a *Application) PlexPlayHandler(c echo.Context) error {
	ratingKey := QueryString(c, "ratingKey")
	player := QueryString(c, "player")
	return a.PlexPlay(c, ratingKey, player)
}
func (a *Application) PlexSessionsHandler(c echo.Context) error {
	return a.PlexSessions(c)
}
func (a *Application) PlexStopHandler(c echo.Context) error {
	session := QueryString(c, "session")
	return a.PlexStop(c, session)
}

// Releases (/releases)
func (a *Application) ReleasesIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.ReleasesIndex(c, page, limit)
}
func (a *Application) ReleasesCreateHandler(c echo.Context) error {
	return a.ReleasesCreate(c)
}
func (a *Application) ReleasesShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ReleasesShow(c, id)
}
func (a *Application) ReleasesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ReleasesUpdate(c, id)
}
func (a *Application) ReleasesSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ReleasesSettings(c, id)
}
func (a *Application) ReleasesDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ReleasesDelete(c, id)
}
func (a *Application) ReleasesPopularHandler(c echo.Context) error {
	interval := c.Param("interval")
	return a.ReleasesPopular(c, interval)
}

// Requests (/requests)
func (a *Application) RequestsIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.RequestsIndex(c, page, limit)
}
func (a *Application) RequestsCreateHandler(c echo.Context) error {
	return a.RequestsCreate(c)
}
func (a *Application) RequestsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.RequestsShow(c, id)
}
func (a *Application) RequestsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.RequestsUpdate(c, id)
}
func (a *Application) RequestsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.RequestsSettings(c, id)
}
func (a *Application) RequestsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.RequestsDelete(c, id)
}

// Series (/series)
func (a *Application) SeriesIndexHandler(c echo.Context) error {
	page := QueryInt(c, "page")
	limit := QueryInt(c, "limit")
	return a.SeriesIndex(c, page, limit)
}
func (a *Application) SeriesCreateHandler(c echo.Context) error {
	return a.SeriesCreate(c)
}
func (a *Application) SeriesShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesShow(c, id)
}
func (a *Application) SeriesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesUpdate(c, id)
}
func (a *Application) SeriesSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesSettings(c, id)
}
func (a *Application) SeriesDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesDelete(c, id)
}
func (a *Application) SeriesCurrentSeasonHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesCurrentSeason(c, id)
}
func (a *Application) SeriesPathsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesPaths(c, id)
}
func (a *Application) SeriesRefreshHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesRefresh(c, id)
}
func (a *Application) SeriesSeasonEpisodesAllHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesSeasonEpisodesAll(c, id)
}
func (a *Application) SeriesSeasonEpisodesHandler(c echo.Context) error {
	id := c.Param("id")
	season := c.Param("season")
	return a.SeriesSeasonEpisodes(c, id, season)
}
func (a *Application) SeriesWatchesHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesWatches(c, id)
}
func (a *Application) SeriesCoversHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesCovers(c, id)
}
func (a *Application) SeriesBackgroundsHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesBackgrounds(c, id)
}
func (a *Application) SeriesJobsHandler(c echo.Context) error {
	id := c.Param("id")
	name := QueryString(c, "name")
	return a.SeriesJobs(c, id, name)
}

// Upcoming (/upcoming)
func (a *Application) UpcomingIndexHandler(c echo.Context) error {
	return a.UpcomingIndex(c)
}

// Users (/users)
func (a *Application) UsersIndexHandler(c echo.Context) error {
	return a.UsersIndex(c)
}

// Watches (/watches)
func (a *Application) WatchesIndexHandler(c echo.Context) error {
	medium_id := QueryString(c, "medium_id")
	username := QueryString(c, "username")
	return a.WatchesIndex(c, medium_id, username)
}
