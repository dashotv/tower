package app

import "fmt"

var activeStates = []string{"searching", "loading", "managing", "downloading", "reviewing"}

func (c *Connector) ActiveDownloads() ([]*Download, error) {
	q := c.Download.Query()
	list, err := q.In("status", activeStates).Run()
	if err != nil {
		return nil, err
	}

	for i, d := range list {
		m := &Medium{}
		err := App().DB.Medium.FindByID(d.MediumId, m)
		if err != nil {
			App().Log.Errorf("could not find medium: %s", d.MediumId)
			continue
		}

		if m.Type == "Episode" && m.SeriesId.Hex() != "" {
			s := &Series{}
			err := App().DB.Series.FindByID(m.SeriesId, s)
			if err != nil {
				App().Log.Errorf("could not find series: %s", d.MediumId)
				continue
			}

			m.Display = fmt.Sprintf("%dx%d %s", m.SeasonNumber, m.EpisodeNumber, m.Title)
			m.Title = s.Title
			for _, p := range s.Paths {
				if p.Type == "cover" {
					m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
					continue
				}
				if p.Type == "background" {
					m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
					continue
				}
			}
		}

		list[i].Medium = *m
	}

	return list, nil
}
