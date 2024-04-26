package app

import (
	"context"

	"github.com/dashotv/tower/internal/plex"
)

func init() {
	initializers = append(initializers, setupPlex)
	initializers = append(initializers, setupPlexFiles)
	starters = append(starters, startPlexFiles)
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

func setupPlexFiles(a *Application) error {
	a.PlexFileCache = &plexFileCache{files: make(map[string]string)}
	return nil
}

func startPlexFiles(ctx context.Context, a *Application) error {
	return a.Workers.Enqueue(&PlexMatch{})
}
