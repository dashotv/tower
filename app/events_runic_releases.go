package app

import (
	runic "github.com/dashotv/runic/app"
)

func onRunicReleases(a *Application, msg *runic.Release) error {
	// handle *runic.Release
	series, err := getSeriesBySearch(msg.Title)
	if series == nil {
		return err
	}

	episode, err := getEpisodeBySeries(series, msg.Season, msg.Episode)
	if episode == nil {
		return err
	}

	d := &Download{}
	d.MediumId = episode.ID
	d.Status = "reviewing"
	d.Url = "metube://" + msg.Download

	if err := a.DB.Download.Save(d); err != nil {
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
	if !list[0].Active {
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
