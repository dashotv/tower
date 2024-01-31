package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/minion"
	"github.com/dashotv/tmdb"
	"github.com/dashotv/tower/internal/plex"
	"github.com/dashotv/tvdb"
)

var app *Application

type setupFunc func(app *Application) error
type healthFunc func(app *Application) error

var initializers = []setupFunc{setupConfig, setupLogger}
var healthchecks = map[string]healthFunc{}

type Application struct {
	Config *Config
	Log    *zap.SugaredLogger

	//golem:template:app/app_partial_definitions
	// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
	// Routes
	Engine  *gin.Engine
	Default *gin.RouterGroup
	Router  *gin.RouterGroup

	// Models
	DB *Connector

	// Events
	Events *Events

	// Workers
	Workers *minion.Minion

	//Cache
	Cache *Cache

	//golem:template:app/app_partial_definitions

	Fanart *Fanart
	Flame  *Flame
	Scry   *Scry
	Plex   *plex.Client
	Tmdb   *tmdb.Client
	Tvdb   *tvdb.Client
}

func Start() error {
	if app != nil {
		return errors.New("application already started")
	}

	app = &Application{}

	for _, f := range initializers {
		// fmt.Printf("running initializer %s\n", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		if err := f(app); err != nil {
			return err
		}
	}

	app.DB.Episode.SetQueryDefaults([]bson.M{{"_type": "Episode"}})
	app.DB.Movie.SetQueryDefaults([]bson.M{{"_type": "Movie"}})
	app.DB.Series.SetQueryDefaults([]bson.M{{"_type": "Series"}})

	app.Workers.Subscribe(app.MinionNotification)
	app.Workers.SubscribeStats(app.MinionStats)

	app.Log.Info("starting tower...")

	//golem:template:app/app_partial_start
	// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
	go app.Events.Start()

	go func() {
		app.Log.Infof("starting workers (%d)...", app.Config.MinionConcurrency)
		app.Workers.Start()
	}()

	app.Routes()
	app.Log.Info("starting routes...")
	if err := app.Engine.Run(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	//golem:template:app/app_partial_start

	return nil
}

func (a *Application) Health() (map[string]bool, error) {
	resp := make(map[string]bool)
	for n, f := range healthchecks {
		err := f(a)
		resp[n] = err == nil
	}

	return resp, nil
}
