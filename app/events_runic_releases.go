package app

import (
	"fmt"
	"time"

	runic "github.com/dashotv/runic/app"
)

func onRunicReleases(a *Application, msg *runic.Release) error {
	log := a.Log.Named("runic.releases")

	// handle *runic.Release
	series, err := a.DB.SeriesBySearch(msg.Title)
	if series == nil {
		return err
	}

	// disable for now, because I want to see if the matching is working
	// if !series.Active {
	// 	return nil, nil
	// }

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

	notice := &EventNotices{
		Event:   "Download Created",
		Time:    time.Now().String(),
		Class:   "runic",
		Level:   "warn",
		Message: fmt.Sprintf("found: %s S%02dE%02d", msg.Title, msg.Season, msg.Episode),
	}
	if err := a.Events.Send("tower.notices", notice); err != nil {
		return err
	}

	return nil
}
