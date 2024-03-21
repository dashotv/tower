package app

import (
	"fmt"
	"regexp"
)

var titleRegex = regexp.MustCompile(`(?i)^(?:episode|chapter)`)

func Destination(m *Medium) (string, error) {
	switch m.Type {
	case "Series", "Movie":
		return fmt.Sprintf("%s/%s/%s", m.Kind, m.Directory, m.Directory), nil
	case "Episode":
		return destinationEpisode(m)
	default:
		return "", fmt.Errorf("unknown type: %s", m.Type)
	}
}

func destinationEpisode(m *Medium) (string, error) {
	s := &Series{}
	err := app.DB.Series.FindByID(m.SeriesId, s)
	if err != nil {
		return "", err
	}

	e := &Episode{}
	err = app.DB.Episode.FindByID(m.ID, e)
	if err != nil {
		return "", err
	}

	out := ""

	if isAnimeKind(string(s.Kind)) && m.AbsoluteNumber > 0 {
		out = fmt.Sprintf("%s/%s/%s - %02dx%02d #%03d", s.Kind, s.Directory, s.Directory, m.SeasonNumber, m.EpisodeNumber, m.AbsoluteNumber)
	} else {
		out = fmt.Sprintf("%s/%s/%s - %02dx%02d", s.Kind, s.Directory, s.Directory, m.SeasonNumber, m.EpisodeNumber)
	}

	if e.Title != "" && !titleRegex.MatchString(e.Title) {
		out = fmt.Sprintf("%s - %s", out, path(e.Title))
	}

	return out, nil
}
