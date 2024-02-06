package app

import (
	"github.com/dashotv/tower/internal/plex"
)

func init() {
	initializers = append(initializers, setupPlex)
}

func setupPlex(app *Application) error {
	p := plex.New(&plex.ClientOptions{
		URL:               app.Config.PlexServerURL,
		Token:             app.Config.PlexToken,
		Debug:             false,
		MachineIdentifier: app.Config.PlexMachineIdentifier,
		ClientIdentifier:  app.Config.PlexClientIdentifier,
		Product:           app.Config.PlexAppName,
		Device:            app.Config.PlexDevice,
		AppName:           app.Config.PlexAppName,
	})

	app.Plex = p
	return nil
}
