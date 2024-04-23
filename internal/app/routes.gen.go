// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/plugins/router"
	"github.com/labstack/echo/v4"
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
	e, err := router.New(logger)
	if err != nil {
		return fae.Wrap(err, "router plugin")
	}
	app.Engine = e
	// unauthenticated routes
	app.Default = app.Engine.Group("")
	// authenticated routes (if enabled, otherwise same as default)
	app.Router = app.Engine.Group("")

	// TODO: fix auth
	if app.Config.Auth {
		clerkSecret := app.Config.ClerkSecretKey
		if clerkSecret == "" {
			app.Log.Fatal("CLERK_SECRET_KEY is not set")
		}
		clerkToken := app.Config.ClerkToken
		if clerkToken == "" {
			app.Log.Fatal("CLERK_TOKEN is not set")
		}

		app.Router.Use(router.ClerkAuth(clerkSecret, clerkToken))
	}

	return nil
}

type Setting struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

type SettingsBatch struct {
	IDs   []string `json:"ids"`
	Name  string   `json:"name"`
	Value bool     `json:"value"`
}

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Total   int64       `json:"total,omitempty"`
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
	combinations.GET("/:name", a.CombinationsShowHandler)
	combinations.POST("/", a.CombinationsCreateHandler)
	combinations.PUT("/:id", a.CombinationsUpdateHandler)

	config := a.Router.Group("/config")
	config.PATCH("/:id", a.ConfigSettingsHandler)

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
	episodes.PATCH("/:id", a.EpisodesSettingsHandler)
	episodes.PUT("/:id", a.EpisodesUpdateHandler)
	episodes.PATCH("/settings", a.EpisodesSettingsBatchHandler)

	feeds := a.Router.Group("/feeds")
	feeds.GET("/", a.FeedsIndexHandler)
	feeds.POST("/", a.FeedsCreateHandler)
	feeds.GET("/:id", a.FeedsShowHandler)
	feeds.PUT("/:id", a.FeedsUpdateHandler)
	feeds.PATCH("/:id", a.FeedsSettingsHandler)
	feeds.DELETE("/:id", a.FeedsDeleteHandler)

	hooks := a.Router.Group("/hooks")
	hooks.GET("/plex", a.HooksPlexHandler)
	hooks.POST("/nzbget", a.HooksNzbgetHandler)

	library := a.Router.Group("/library")
	library.GET("/", a.LibraryIndexHandler)
	library.POST("/", a.LibraryCreateHandler)
	library.GET("/:id", a.LibraryShowHandler)
	library.PUT("/:id", a.LibraryUpdateHandler)
	library.PATCH("/:id", a.LibrarySettingsHandler)
	library.DELETE("/:id", a.LibraryDeleteHandler)

	library_template := a.Router.Group("/library_template")
	library_template.GET("/", a.LibraryTemplateIndexHandler)
	library_template.POST("/", a.LibraryTemplateCreateHandler)
	library_template.GET("/:id", a.LibraryTemplateShowHandler)
	library_template.PUT("/:id", a.LibraryTemplateUpdateHandler)
	library_template.PATCH("/:id", a.LibraryTemplateSettingsHandler)
	library_template.DELETE("/:id", a.LibraryTemplateDeleteHandler)

	library_type := a.Router.Group("/library_type")
	library_type.GET("/", a.LibraryTypeIndexHandler)
	library_type.POST("/", a.LibraryTypeCreateHandler)
	library_type.GET("/:id", a.LibraryTypeShowHandler)
	library_type.PUT("/:id", a.LibraryTypeUpdateHandler)
	library_type.PATCH("/:id", a.LibraryTypeSettingsHandler)
	library_type.DELETE("/:id", a.LibraryTypeDeleteHandler)

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
	movies.POST("/:id/jobs", a.MoviesJobsHandler)

	paths := a.Router.Group("/paths")
	paths.POST("/:id", a.PathsUpdateHandler)
	paths.DELETE("/:id", a.PathsDeleteHandler)

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

	want := a.Router.Group("/want")
	want.GET("/series", a.WantSeriesHandler)
	want.GET("/movie", a.WantMovieHandler)

	watches := a.Router.Group("/watches")
	watches.GET("/", a.WatchesIndexHandler)
	watches.POST("/", a.WatchesCreateHandler)
	watches.DELETE("/:id", a.WatchesDeleteHandler)
	watches.DELETE("/medium", a.WatchesDeleteMediumHandler)

}

func (a *Application) indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, router.H{
		"name": "tower",
		"routes": router.H{
			"collections":      "/collections",
			"combinations":     "/combinations",
			"config":           "/config",
			"downloads":        "/downloads",
			"episodes":         "/episodes",
			"feeds":            "/feeds",
			"hooks":            "/hooks",
			"library":          "/library",
			"library_template": "/library_template",
			"library_type":     "/library_type",
			"messages":         "/messages",
			"movies":           "/movies",
			"paths":            "/paths",
			"plex":             "/plex",
			"releases":         "/releases",
			"requests":         "/requests",
			"series":           "/series",
			"upcoming":         "/upcoming",
			"users":            "/users",
			"want":             "/want",
			"watches":          "/watches",
		},
	})
}

func (a *Application) healthHandler(c echo.Context) error {
	health, err := a.Health()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, router.H{"name": "tower", "health": health})
}

// Collections (/collections)
func (a *Application) CollectionsIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.CollectionsIndex(c, page, limit)
}
func (a *Application) CollectionsCreateHandler(c echo.Context) error {
	subject := &Collection{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.CollectionsCreate(c, subject)
}
func (a *Application) CollectionsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CollectionsShow(c, id)
}
func (a *Application) CollectionsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Collection{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.CollectionsUpdate(c, id, subject)
}
func (a *Application) CollectionsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.CollectionsSettings(c, id, setting)
}
func (a *Application) CollectionsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.CollectionsDelete(c, id)
}

// Combinations (/combinations)
func (a *Application) CombinationsIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.CombinationsIndex(c, page, limit)
}
func (a *Application) CombinationsShowHandler(c echo.Context) error {
	name := c.Param("name")
	return a.CombinationsShow(c, name)
}
func (a *Application) CombinationsCreateHandler(c echo.Context) error {
	subject := &Combination{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.CombinationsCreate(c, subject)
}
func (a *Application) CombinationsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Combination{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.CombinationsUpdate(c, id, subject)
}

// Config (/config)
func (a *Application) ConfigSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	settings := &Setting{}
	if err := c.Bind(settings); err != nil {
		return err
	}
	return a.ConfigSettings(c, id, settings)
}

// Downloads (/downloads)
func (a *Application) DownloadsIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.DownloadsIndex(c, page, limit)
}
func (a *Application) DownloadsCreateHandler(c echo.Context) error {
	subject := &Download{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.DownloadsCreate(c, subject)
}
func (a *Application) DownloadsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsShow(c, id)
}
func (a *Application) DownloadsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Download{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.DownloadsUpdate(c, id, subject)
}
func (a *Application) DownloadsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.DownloadsSettings(c, id, setting)
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
	page := router.QueryParamInt(c, "page")
	medium_id := router.QueryParamString(c, "medium_id")
	return a.DownloadsRecent(c, page, medium_id)
}
func (a *Application) DownloadsSelectHandler(c echo.Context) error {
	id := c.Param("id")
	medium_id := router.QueryParamString(c, "medium_id")
	num := router.QueryParamInt(c, "num")
	return a.DownloadsSelect(c, id, medium_id, num)
}
func (a *Application) DownloadsTorrentHandler(c echo.Context) error {
	id := c.Param("id")
	return a.DownloadsTorrent(c, id)
}

// Episodes (/episodes)
func (a *Application) EpisodesSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.EpisodesSettings(c, id, setting)
}
func (a *Application) EpisodesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	episode := &Episode{}
	if err := c.Bind(episode); err != nil {
		return err
	}
	return a.EpisodesUpdate(c, id, episode)
}
func (a *Application) EpisodesSettingsBatchHandler(c echo.Context) error {
	settings := &SettingsBatch{}
	if err := c.Bind(settings); err != nil {
		return err
	}
	return a.EpisodesSettingsBatch(c, settings)
}

// Feeds (/feeds)
func (a *Application) FeedsIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.FeedsIndex(c, page, limit)
}
func (a *Application) FeedsCreateHandler(c echo.Context) error {
	subject := &Feed{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.FeedsCreate(c, subject)
}
func (a *Application) FeedsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.FeedsShow(c, id)
}
func (a *Application) FeedsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Feed{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.FeedsUpdate(c, id, subject)
}
func (a *Application) FeedsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.FeedsSettings(c, id, setting)
}
func (a *Application) FeedsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.FeedsDelete(c, id)
}

// Hooks (/hooks)
func (a *Application) HooksPlexHandler(c echo.Context) error {
	return a.HooksPlex(c)
}
func (a *Application) HooksNzbgetHandler(c echo.Context) error {
	payload := &NzbgetPayload{}
	if err := c.Bind(payload); err != nil {
		return err
	}
	return a.HooksNzbget(c, payload)
}

// Library (/library)
func (a *Application) LibraryIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.LibraryIndex(c, page, limit)
}
func (a *Application) LibraryCreateHandler(c echo.Context) error {
	subject := &Library{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.LibraryCreate(c, subject)
}
func (a *Application) LibraryShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.LibraryShow(c, id)
}
func (a *Application) LibraryUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Library{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.LibraryUpdate(c, id, subject)
}
func (a *Application) LibrarySettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.LibrarySettings(c, id, setting)
}
func (a *Application) LibraryDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.LibraryDelete(c, id)
}

// LibraryTemplate (/library_template)
func (a *Application) LibraryTemplateIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.LibraryTemplateIndex(c, page, limit)
}
func (a *Application) LibraryTemplateCreateHandler(c echo.Context) error {
	subject := &LibraryTemplate{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.LibraryTemplateCreate(c, subject)
}
func (a *Application) LibraryTemplateShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.LibraryTemplateShow(c, id)
}
func (a *Application) LibraryTemplateUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &LibraryTemplate{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.LibraryTemplateUpdate(c, id, subject)
}
func (a *Application) LibraryTemplateSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.LibraryTemplateSettings(c, id, setting)
}
func (a *Application) LibraryTemplateDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.LibraryTemplateDelete(c, id)
}

// LibraryType (/library_type)
func (a *Application) LibraryTypeIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.LibraryTypeIndex(c, page, limit)
}
func (a *Application) LibraryTypeCreateHandler(c echo.Context) error {
	subject := &LibraryType{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.LibraryTypeCreate(c, subject)
}
func (a *Application) LibraryTypeShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.LibraryTypeShow(c, id)
}
func (a *Application) LibraryTypeUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &LibraryType{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.LibraryTypeUpdate(c, id, subject)
}
func (a *Application) LibraryTypeSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.LibraryTypeSettings(c, id, setting)
}
func (a *Application) LibraryTypeDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.LibraryTypeDelete(c, id)
}

// Messages (/messages)
func (a *Application) MessagesIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.MessagesIndex(c, page, limit)
}
func (a *Application) MessagesCreateHandler(c echo.Context) error {
	message := &Message{}
	if err := c.Bind(message); err != nil {
		return err
	}
	return a.MessagesCreate(c, message)
}

// Movies (/movies)
func (a *Application) MoviesIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	kind := router.QueryParamString(c, "kind")
	source := router.QueryParamString(c, "source")
	downloaded := router.QueryParamBool(c, "downloaded")
	completed := router.QueryParamBool(c, "completed")
	broken := router.QueryParamBool(c, "broken")
	return a.MoviesIndex(c, page, limit, kind, source, downloaded, completed, broken)
}
func (a *Application) MoviesCreateHandler(c echo.Context) error {
	subject := &Movie{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.MoviesCreate(c, subject)
}
func (a *Application) MoviesShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.MoviesShow(c, id)
}
func (a *Application) MoviesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Movie{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.MoviesUpdate(c, id, subject)
}
func (a *Application) MoviesSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.MoviesSettings(c, id, setting)
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
func (a *Application) MoviesJobsHandler(c echo.Context) error {
	id := c.Param("id")
	name := router.QueryParamString(c, "name")
	return a.MoviesJobs(c, id, name)
}

// Paths (/paths)
func (a *Application) PathsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	medium_id := router.QueryParamString(c, "medium_id")
	path := &Path{}
	if err := c.Bind(path); err != nil {
		return err
	}
	return a.PathsUpdate(c, id, medium_id, path)
}
func (a *Application) PathsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	medium_id := router.QueryParamString(c, "medium_id")
	return a.PathsDelete(c, id, medium_id)
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
	query := router.QueryParamString(c, "query")
	section := router.QueryParamString(c, "section")
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
	ratingKey := router.QueryParamString(c, "ratingKey")
	player := router.QueryParamString(c, "player")
	return a.PlexPlay(c, ratingKey, player)
}
func (a *Application) PlexSessionsHandler(c echo.Context) error {
	return a.PlexSessions(c)
}
func (a *Application) PlexStopHandler(c echo.Context) error {
	session := router.QueryParamString(c, "session")
	return a.PlexStop(c, session)
}

// Releases (/releases)
func (a *Application) ReleasesIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.ReleasesIndex(c, page, limit)
}
func (a *Application) ReleasesCreateHandler(c echo.Context) error {
	subject := &Release{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.ReleasesCreate(c, subject)
}
func (a *Application) ReleasesShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.ReleasesShow(c, id)
}
func (a *Application) ReleasesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Release{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.ReleasesUpdate(c, id, subject)
}
func (a *Application) ReleasesSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.ReleasesSettings(c, id, setting)
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
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	return a.RequestsIndex(c, page, limit)
}
func (a *Application) RequestsCreateHandler(c echo.Context) error {
	subject := &Request{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.RequestsCreate(c, subject)
}
func (a *Application) RequestsShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.RequestsShow(c, id)
}
func (a *Application) RequestsUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Request{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.RequestsUpdate(c, id, subject)
}
func (a *Application) RequestsSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.RequestsSettings(c, id, setting)
}
func (a *Application) RequestsDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.RequestsDelete(c, id)
}

// Series (/series)
func (a *Application) SeriesIndexHandler(c echo.Context) error {
	page := router.QueryParamIntDefault(c, "page", "1")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	kind := router.QueryParamString(c, "kind")
	source := router.QueryParamString(c, "source")
	active := router.QueryParamBool(c, "active")
	favorite := router.QueryParamBool(c, "favorite")
	broken := router.QueryParamBool(c, "broken")
	return a.SeriesIndex(c, page, limit, kind, source, active, favorite, broken)
}
func (a *Application) SeriesCreateHandler(c echo.Context) error {
	subject := &Series{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.SeriesCreate(c, subject)
}
func (a *Application) SeriesShowHandler(c echo.Context) error {
	id := c.Param("id")
	return a.SeriesShow(c, id)
}
func (a *Application) SeriesUpdateHandler(c echo.Context) error {
	id := c.Param("id")
	subject := &Series{}
	if err := c.Bind(subject); err != nil {
		return err
	}
	return a.SeriesUpdate(c, id, subject)
}
func (a *Application) SeriesSettingsHandler(c echo.Context) error {
	id := c.Param("id")
	setting := &Setting{}
	if err := c.Bind(setting); err != nil {
		return err
	}
	return a.SeriesSettings(c, id, setting)
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
	name := router.QueryParamString(c, "name")
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

// Want (/want)
func (a *Application) WantSeriesHandler(c echo.Context) error {
	id := router.QueryParamString(c, "id")
	return a.WantSeries(c, id)
}
func (a *Application) WantMovieHandler(c echo.Context) error {
	id := router.QueryParamString(c, "id")
	return a.WantMovie(c, id)
}

// Watches (/watches)
func (a *Application) WatchesIndexHandler(c echo.Context) error {
	medium_id := router.QueryParamString(c, "medium_id")
	username := router.QueryParamString(c, "username")
	return a.WatchesIndex(c, medium_id, username)
}
func (a *Application) WatchesCreateHandler(c echo.Context) error {
	medium_id := router.QueryParamString(c, "medium_id")
	username := router.QueryParamString(c, "username")
	return a.WatchesCreate(c, medium_id, username)
}
func (a *Application) WatchesDeleteHandler(c echo.Context) error {
	id := c.Param("id")
	return a.WatchesDelete(c, id)
}
func (a *Application) WatchesDeleteMediumHandler(c echo.Context) error {
	medium_id := router.QueryParamString(c, "medium_id")
	return a.WatchesDeleteMedium(c, medium_id)
}
