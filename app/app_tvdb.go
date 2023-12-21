package app

import (
	"github.com/pkg/errors"

	"github.com/dashotv/tvdb"
)

func init() {
	initializers = append(initializers, setupTvdb)
}

func setupTvdb(app *Application) error {
	c, err := tvdb.Login(app.Config.TvdbKey)
	if err != nil {
		return errors.Wrap(err, "tvdb login")
	}
	app.Tvdb = c
	return nil
}
