package app

import (
	"fmt"

	runic "github.com/dashotv/runic/app"
)

func onRunicReleases(a *Application, msg *runic.Release) error {
	log := a.Log.Named("runic.releases")

	if msg.Size > 0 && msg.Size < 100000000 {
		// log.Warnf("skipping: %s %02dx%02d: size %d < 100mb", msg.Title, msg.Season, msg.Episode, msg.Size)
		return nil
	}

	series, err := a.DB.SeriesBySearch(msg.Title)
	if series == nil {
		return err
	}

	// disable for now, because I want to see if the matching is working
	// if !series.Active {
	// 	return nil
	// }

	log.Infof("received: %s %02dx%02d", msg.Title, msg.Season, msg.Episode)

	episode, err := app.DB.SeriesEpisodeBy(series, msg.Season, msg.Episode)
	if episode == nil {
		return err
	}

	log.Warnf("found: %s s%02de%02d", msg.Title, msg.Season, msg.Episode)

	d := &Download{}
	d.MediumId = episode.ID
	d.Status = "reviewing"
	d.Url = "metube://" + msg.Download

	if err := a.DB.Download.Save(d); err != nil {
		return err
	}

	notifier.Info("Download Created", fmt.Sprintf("found: %s S%02dE%02d", msg.Title, msg.Season, msg.Episode))
	return nil
}
