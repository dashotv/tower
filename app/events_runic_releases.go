package app

import (
	"fmt"
	"time"

	runic "github.com/dashotv/runic/app"
)

func onRunicReleases(a *Application, msg *runic.Release) error {
	// handle *runic.Release
	series, err := getSeriesBySearch(msg.Title)
	if series == nil {
		return err
	}

	// disable for now, because I want to see if the matching is working
	// if !series.Active {
	// 	return nil, nil
	// }

	episode, err := getEpisodeBySeries(series, msg.Season, msg.Episode)
	if episode == nil {
		return err
	}

	a.Log.Named("runic.releases").Warnf("found: %s S%02dE%02d", msg.Title, msg.Season, msg.Episode)

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

func getSeriesBySearch(title string) (*Series, error) {
	list, err := app.DB.Series.Query().Where("search", title).Run()
	if err != nil {
		return nil, err
	}
	if len(list) != 1 {
		return nil, nil
	}

	return list[0], nil
}

func getEpisodeBySeries(s *Series, season, episode int) (*Episode, error) {
	list, err := app.DB.Episode.Query().Where("series_id", s.ID).Where("season", season).Where("episode", episode).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 1 {
		return list[0], nil
	}

	list, err = app.DB.Episode.Query().Where("series_id", s.ID).Where("absolute_number", episode).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 1 && !list[0].Completed && !list[0].Downloaded {
		return list[0], nil
	}

	return nil, nil
}
