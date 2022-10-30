package app

import "fmt"

var activeStates = []string{"searching", "loading", "managing", "downloading", "reviewing"}

func (c *Connector) ActiveDownloads() ([]*Download, error) {
	q := c.Download.Query()
	list, err := q.In("status", activeStates).Run()
	if err != nil {
		return nil, err
	}

	processDownloads(list)
	return list, nil
}

func (c *Connector) RecentDownloads() ([]*Download, error) {
	q := c.Download.Query()
	list, err := q.Where("status", "done").
		Desc("updated_at").Desc("created_at").
		Limit(pagesize).
		Run()
	if err != nil {
		return nil, err
	}

	processDownloads(list)
	return list, nil
}

func processDownloads(list []*Download) {
	for i, d := range list {
		m := &Medium{}
		err := App().DB.Medium.FindByID(d.MediumId, m)
		if err != nil {
			App().Log.Errorf("could not find medium: %s", d.MediumId)
			continue
		}

		paths := m.Paths
		m.Display = m.Type
		if m.Type == "Episode" && m.SeriesId.Hex() != "" {
			s := &Series{}
			err := App().DB.Series.FindByID(m.SeriesId, s)
			if err != nil {
				App().Log.Errorf("could not find series: %s", d.MediumId)
				continue
			}

			m.Display = fmt.Sprintf("%dx%d %s", m.SeasonNumber, m.EpisodeNumber, m.Title)
			m.Title = s.Title
			paths = s.Paths
		}

		for _, p := range paths {
			if p.Type == "cover" {
				m.Cover = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
			if p.Type == "background" {
				m.Background = fmt.Sprintf("%s/%s.%s", imagesBaseURL, p.Local, p.Extension)
				continue
			}
		}

		list[i].Medium = *m
	}
}
