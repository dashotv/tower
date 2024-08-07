package app

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/fae"
	flame "github.com/dashotv/flame/client"
	"github.com/dashotv/golem/plugins/cache"
	"github.com/dashotv/minion"
	runic "github.com/dashotv/runic/client"
	scry "github.com/dashotv/scry/client"
	"github.com/dashotv/tmdb"
	"github.com/dashotv/tower/internal/importer"
	"github.com/dashotv/tower/internal/plex"
	"github.com/dashotv/tvdb"
)

var app *Application

type setupFunc func(app *Application) error
type healthFunc func(app *Application) error
type startFunc func(ctx context.Context, app *Application) error

var initializers = []setupFunc{setupConfig, setupLogger}
var healthchecks = map[string]healthFunc{}
var starters = []startFunc{}

type Application struct {
	Config *Config
	Log    *zap.SugaredLogger

	//golem:template:app/app_partial_definitions
	// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
	// Routes
	Engine  *echo.Echo
	Default *echo.Group
	Router  *echo.Group

	// Models
	DB *Connector

	// Events
	Events *Events

	// Workers
	Workers *minion.Minion

	//Cache
	Cache *cache.Cache

	//golem:template:app/app_partial_definitions

	Fanart        *Fanart
	Flame         *flame.Client
	Scry          *scry.Client
	Runic         *runic.Client
	Plex          *plex.Client
	PlexFileCache *plexFileCache
	Tmdb          *tmdb.Client
	Tvdb          *tvdb.Client
	Importer      *importer.Importer
	Want          *Want
	Destinator    *Destinator
	Libs          map[string]*Library
}

type Plugin interface {
	Name() string
	Setup() error
	Start(context.Context) error
}

func Setup() error {
	if app != nil {
		return fae.New("application already setup")
	}

	app = &Application{}

	for _, f := range initializers {
		if err := f(app); err != nil {
			return err
		}
	}

	app.DB.Episode.SetQueryDefaults([]bson.M{{"_type": "Episode"}})
	app.DB.Movie.SetQueryDefaults([]bson.M{{"_type": "Movie"}})
	app.DB.Series.SetQueryDefaults([]bson.M{{"_type": "Series"}})

	return nil
}

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if app == nil {
		if err := Setup(); err != nil {
			return err
		}
	}

	ctx = ContextSet(ctx, "app", app)

	// app.Engine.Debug = true
	// app.Engine.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	fmt.Printf("req: %s\n", reqBody)
	// 	fmt.Printf("res: %s\n", resBody)
	// }))
	app.Engine.HTTPErrorHandler = customHTTPErrorHandler

	for _, f := range starters {
		if err := f(ctx, app); err != nil {
			return err
		}
	}

	app.Log.Info("started")

	for {
		select {
		case <-ctx.Done():
			app.Log.Info("stopping")
			return nil
		}
	}
}

func (a *Application) Health() (map[string]bool, error) {
	resp := make(map[string]bool)
	for n, f := range healthchecks {
		err := f(a)
		resp[n] = err == nil
	}

	return resp, nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	if !c.Response().Committed {
		c.JSON(code, &Response{Error: true, Message: err.Error()})
	}
}
