package app

import "github.com/dashotv/tower/internal/importer"

func init() {
	initializers = append(initializers, setupImporter)
}

func setupImporter(a *Application) error {
	opts := &importer.Options{
		TvdbKey:   a.Config.TvdbKey,
		TmdbToken: a.Config.TmdbToken,
		FanartKey: a.Config.FanartApiKey,
		FanartURL: a.Config.FanartApiURL,
	}

	i, err := importer.New(opts)
	if err != nil {
		return err
	}

	a.Importer = i
	return nil
}
