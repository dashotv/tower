package importer

import (
	"errors"
	"fmt"

	"github.com/dashotv/tvdb"
)

func (i *Importer) loadEpisodes(tvdbid int64, episodeOrder int) ([]*Episode, error) {
	resp, err := i.Tvdb.GetSeriesSeasonEpisodesTranslated(tvdbid, i.Opts.Language, 0, episodeOrderString(episodeOrder))
	if err != nil {
		return nil, fmt.Errorf("translated: %w", err)
	}
	if resp.Data == nil {
		return nil, errors.New("translated: no data")
	}

	episodes := make([]*Episode, 0)
	for _, e := range resp.Data.Episodes {
		ep := &Episode{
			Title:       tvdb.StringValue(e.Name),
			Description: tvdb.StringValue(e.Overview),
			Airdate:     tvdb.StringValue(e.Aired),
			Season:      int(tvdb.Int64Value(e.SeasonNumber)),
			Episode:     int(tvdb.Int64Value(e.Number)),
		}
		episodes = append(episodes, ep)
	}

	return episodes, nil
}
