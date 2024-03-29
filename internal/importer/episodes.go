package importer

import (
	"github.com/dashotv/fae"
	"github.com/dashotv/tvdb"
)

func (i *Importer) loadEpisodes(tvdbid int64) ([]*Episode, error) {
	def, err := i.loadEpisodesMap(tvdbid, EpisodeOrderDefault)
	if err != nil {
		return nil, fae.Wrap(err, "default")
	}

	abs, err := i.loadEpisodesMap(tvdbid, EpisodeOrderAbsolute)
	if err != nil {
		return nil, fae.Wrap(err, "absolute")
	}

	episodes := make([]*Episode, 0)
	for id, ep := range def {
		if aep, ok := abs[id]; ok {
			ep.Absolute = aep.Episode
		}
		episodes = append(episodes, ep)
	}

	return episodes, nil
}

func (i *Importer) loadEpisodesMap(tvdbid int64, episodeOrder int) (map[int64]*Episode, error) {
	req := tvdb.GetSeriesEpisodesRequest{
		ID:         tvdbid,
		Page:       0,
		SeasonType: episodeOrderString(episodeOrder),
	}
	resp, err := i.Tvdb.GetSeriesEpisodes(req)
	if err != nil {
		return nil, fae.Wrap(err, "episodes")
	}
	if resp.Data == nil {
		return nil, fae.New("episodes: no data")
	}

	episodeMap := make(map[int64]*Episode)
	for _, e := range resp.Data.Episodes {
		ep := &Episode{
			ID:          tvdb.Int64Value(e.ID),
			Title:       tvdb.StringValue(e.Name),
			Description: tvdb.StringValue(e.Overview),
			Airdate:     tvdb.StringValue(e.Aired),
			Season:      int(tvdb.Int64Value(e.SeasonNumber)),
			Episode:     int(tvdb.Int64Value(e.Number)),
		}
		episodeMap[ep.ID] = ep
	}

	trans, err := i.Tvdb.GetSeriesSeasonEpisodesTranslated(tvdbid, i.Opts.Language, 0, episodeOrderString(episodeOrder))
	if err != nil {
		return nil, fae.Wrap(err, "translated")
	}
	if trans.Data == nil {
		return nil, fae.New("translated: no data")
	}
	for _, e := range trans.Data.Episodes {
		if ep, ok := episodeMap[tvdb.Int64Value(e.ID)]; ok {
			if e.Name != nil {
				ep.Title = tvdb.StringValue(e.Name)
			}
			if e.Overview != nil {
				ep.Description = tvdb.StringValue(e.Overview)
			}
		}
	}

	return episodeMap, nil
}
