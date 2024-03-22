package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/minion"
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
	Cache *Cache

	//golem:template:app/app_partial_definitions

	Fanart   *Fanart
	Flame    *Flame
	Scry     *Scry
	Plex     *plex.Client
	Tmdb     *tmdb.Client
	Tvdb     *tvdb.Client
	Importer *importer.Importer
}

func Setup() error {
	if app != nil {
		return errors.New("application already setup")
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
	app.Workers.Subscribe(app.MinionNotification)
	app.Workers.SubscribeStats(app.MinionStats)

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
