package app

import (
	"github.com/dashotv/tmdb"
)

func init() {
	initializers = append(initializers, setupTmdb)
}

var posterRatio float32 = 0.6666666666666666
var backgroundRatio float32 = 1.7777777777777777

func setupTmdb(app *Application) error {
	app.Tmdb = tmdb.New(app.Config.TmdbToken)
	return nil
}
